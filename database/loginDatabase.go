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

	"TODO/models"
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
	creds := &userInfo
	// First we check for any errors. If there are no errors when retrieving the user
	// from the database, it means that there exists an entry with that email already.
	// To prevent duplicate entries, we check for this and return the error
	err := DB.Take(creds).Error
	if err == nil {
		return errors.New("There is already an account with this email. Please login instead.")
	}
	// As the email does not exist in the database, we first hash it before adding it
	// The 8 represents the cost of hashing. 8 is chosen arbitrarily
	// We also salt the password for extra security
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 8)
	if err == nil {
		// If no errors, we can add the user info to the database
		userInfo.Password = string(hashedPass)
		err := DB.Create(creds)
		if err == nil {
			return nil
		}
		return err.Error
	}
	return err
}
