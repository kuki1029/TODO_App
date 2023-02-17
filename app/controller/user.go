package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"fmt"
	"time"
	"todo/app/models"
	"todo/app/utils/password"
)

// Interface for UserRepo
type UserRepo interface {
	AddUser(userInfo models.UserDTO) error
	UserExists(userInfo *models.UserDTO, email string) error
	AuthenticateUser(userInfo *models.UserDTO, userInfo2 *models.UserDTO) bool
	ReturnUserID(userInfo models.UserDTO) uint
}

// Interface for RedisClient
type RedisClient interface {
	Ping() error
	GetFromRedis(key string) (uint, error)
	SetInRedis(value string, redisVal string, timeAmt time.Duration)
	DelInRedis(value string)
}

// UserController handles all routes related to users
type UserController struct {
	userRepo    UserRepo
	redisClient RedisClient
}

// NewUserController creates a new instance of UserController
func NewUserController(ur UserRepo, rc RedisClient) *UserController {
	return &UserController{
		userRepo:    ur,
		redisClient: rc,
	}
}

// This function will create cookies for the user and log them in
// so that they can view their tasks.
func (uc *UserController) Login(ctx *fiber.Ctx) error {
	var creds models.UserDTO
	// First we need to parse the variable ctx to receive the credentials
	err := ctx.BodyParser(&creds)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}
	// Once we have their credentials, we need to check if the user is in the database
	userTemp := new(models.UserDTO)
	err = uc.userRepo.UserExists(userTemp, creds.Email)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Incorrect password or account does not exist. Please create an account or try again.",
		})
	}
	userVerified := uc.userRepo.AuthenticateUser(userTemp, &creds)
	if userVerified {
		// To be able to identify this user on other pages, we need to create a cookie for their browser
		cookie := setCookie(ctx, "sessionKey", uuid.NewString(), 24)
		ID := uc.userRepo.ReturnUserID(creds)
		redisVal := "email: " + creds.Email + " id: " + fmt.Sprint(ID)
		// We set it to expire in 24 hours
		uc.redisClient.SetInRedis(cookie.Value, redisVal, 24*time.Hour)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Login successful.",
		})
	} else {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Incorrect password. Please try again.",
		})
	}

}

// This function will be called through the JS and handle any signup requirements
// the c variable here contains all the required credentials
func (uc *UserController) Signup(ctx *fiber.Ctx) error {
	var creds models.UserDTO
	// First we need to parse the variable ctx to receive the credentials
	err := ctx.BodyParser(&creds)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}
	// Once we have the required data, we need to make sure the user isn't a duplicate
	userTemp := new(models.UserDTO)
	err = uc.userRepo.UserExists(userTemp, creds.Email)
	if err == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "This email has already been used for an account. Please sign in.",
		})
	}

	// If the user isn't in the database, we can add the user to the database and their credentials
	creds.Password = password.Generate(creds.Password)
	err = uc.userRepo.AddUser(creds)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Account created.",
		})
	}
}

// This function logs out the user by replacing the stored cookies
func (uc *UserController) Logout(ctx *fiber.Ctx) error {
	// Delete from redis cache
	uc.redisClient.DelInRedis(ctx.Cookies("sessionKey"))
	// Make new cookie to replace current cookie
	// This makes it expire. We use -100 so it expires for older versions of IE too
	setCookie(ctx, "sessionKey", "", -100)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logged out.",
	})
}

// This function will take in parameters for the cookie and set them to the fiber context
func setCookie(ctx *fiber.Ctx, name string, value string, timeAmt time.Duration) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = name
	// We generate a random key to store in the cookie value. Also stored in redis cache
	cookie.Value = value
	cookie.Expires = time.Now().Add(timeAmt * time.Hour)
	ctx.Cookie(cookie)
	return cookie
}
