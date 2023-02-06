package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todo/models"
)

// DB represents a Database instance
var DB *gorm.DB

// This function will connect to the DB and setup all required variables
func ConnectToDB() {
	err := godotenv.Load("local.env")
	// Check for errors with loading the enviroment variables
	if err != nil {
		// We use log.Fatal as the entire program depends on the successfull
		// loading of this file
		log.Fatal("Error loading env file \n", err)
	}
	// Create the string so we can login and open the PostgreSQL DB
	loginDB := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s",
		os.Getenv("PSQL_USER"), os.Getenv("PSQL_PASS"),
		os.Getenv("PSQL_DBNAME"), os.Getenv("PSQL_PORT"))
	// Here we open the actual connection
	DB, err = gorm.Open(postgres.Open(loginDB), &gorm.Config{})
	// Now we check for errors to be sure everything is okay
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	} else {
		DB.AutoMigrate(models.User{})
	}
}

// This function will add the user credentials to the database
// The password will be hashed for security reasons
func AddUser(userInfo models.User) error {
	// Get the pointer for the model
	//creds := &userInfo
	// First we check for any errors. If there are no errors when retrieving the user
	// from the database, it means that there exists an entry with that email already.
	// To prevent duplicate entries, we check for this and return the error
	var tempUser models.User
	err := DB.Where("email = ?", userInfo.Email).First(&tempUser).Error
	if err == nil {
		return errors.New("there is already an account with this email. please login instead")
	}
	// As the email does not exist in the database, we first hash it before adding it
	// We also salt the password for extra security
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err == nil {
		// If no errors, we can add the user info to the database
		userInfo.Password = string(hashedPass)

		err := DB.Create(&userInfo)
		if err.Error == nil {
			return nil
		}
		return err.Error
	}
	return err
}

// This function will check if the user exists in the database or not
// This will allow us to authenticate logins
func AuthenticateUser(userInfo models.User) error {
	// Need to create temp to store the details if the user does exist
	var tempUser models.User
	err := DB.Where("email = ?", userInfo.Email).First(&tempUser).Error
	if err == nil {
		// Here we can compare the hashed password obtained from the database and the password entered by the user
		err = bcrypt.CompareHashAndPassword([]byte(tempUser.Password), []byte(userInfo.Password))
		// If we get no errors above, it means that the two passwords match and we can login the user
		if err == nil {
			return nil
		} else {
			return errors.New("the password is incorrect. Please try again")
		}
	} else {
		return errors.New("there is no account under this email. Please signup first")
	}
}

// This function returns the user ID
func ReturnUserID(userInfo models.User) uint {
	var tempUser models.User
	DB.Where("email = ?", userInfo.Email).First(&tempUser)
	return tempUser.ID
}

// This function returns the name for a particular user using their ID
func ReturnName(ID uint) (string, error) {
	tempUser := models.User{}
	// As the user model stores a task struct, and not TaskResponse, we need to create
	// another variable so we can return the TaskResponse
	err := DB.Where("ID = ?", ID).Find(&tempUser).Error

	if err != nil {
		// Return blank name incase of error
		return "", err
	}
	return tempUser.Name, nil

}
