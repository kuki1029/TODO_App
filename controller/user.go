package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"fmt"
	"time"
	"todo/middleware"
	"todo/models"
	"todo/repo"
)

var client = middleware.RedisSetUp()
var _ = middleware.Ping(client)

// This function will create cookies for the user and log them in
// so that they can view their tasks.
func Login(ctx *fiber.Ctx) error {
	var creds models.UserDTO
	db := repo.DB.DbConn
	// First we need to parse the variable ctx to receive the credentials
	err := ctx.BodyParser(&creds)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}
	// Once we have their credentials, we need to check if the user is in the database
	err = repo.AuthenticateUser(creds, db)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	// To be able to identify this user on other pages, we need to create a cookie for their browser
	cookie := new(fiber.Cookie)
	cookie.Name = "sessionKey"
	// We generate a random key to store in the cookie value. Also stored in redis cache
	cookie.Value = uuid.NewString()
	cookie.Expires = time.Now().Add(24 * time.Hour)
	ctx.Cookie(cookie)
	ID := repo.ReturnUserID(creds, db)
	redisVal := "email: " + creds.Email + " id: " + fmt.Sprint(ID)
	// We set it to expire in 24 hours
	client.Set(cookie.Value, redisVal, 24*time.Hour)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login successful.",
	})

}

// This function will be called through the JS and handle any signup requirements
// the c variable here contains all the required credentials
func Signup(ctx *fiber.Ctx) error {
	var creds models.UserDTO
	db := repo.DB.DbConn
	// First we need to parse the variable ctx to receive the credentials
	err := ctx.BodyParser(&creds)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}
	// Once we have the required data, we need to make sure the user isn't a duplicate
	err = repo.AddUser(creds, db)
	if err != nil {
		// In this case, we know the user is a duplicate, so we returen an error message
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User already exists. Please login or use a different email address.",
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Account created.",
		})
	}
}

// This function logs out the user by replacing the stored cookies
func Logout(ctx *fiber.Ctx) error {
	// Delete from redis cache
	client.Del(ctx.Cookies("sessionKey"))
	// Make new cookie to replace current cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "sessionKey"
	cookie.Value = ""
	// This makes it expire. We use -100 so it expires for older versions of IE too
	cookie.Expires = time.Now().Add(-100 * time.Hour)
	ctx.Cookie(cookie)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logged out.",
	})
}
