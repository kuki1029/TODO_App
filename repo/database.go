package repo

import (
	"fmt"
	"log"
	"os"
	"todo/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	DbConn *gorm.DB
}

// DB represents a Database instance
var DB Dbinstance

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
	DB.DbConn, err = gorm.Open(postgres.Open(loginDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	// Now we check for errors to be sure everything is okay
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	} else {
		DB.DbConn.AutoMigrate(models.User{})
	}
}
