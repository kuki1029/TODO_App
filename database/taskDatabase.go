package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/copier"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"TODO/models"
)

// DB represents a Database instance
var DBTasks *gorm.DB

// This function will connect to the DB and setup all required variables
func ConnectToDBTasks() {
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
		DB.AutoMigrate(models.Task{})
	}
}

// This function returns the tasks for a particular user by looking them up through their ID
func ReturnTasksWithID(ID uint) ([]models.TaskResponse, error) {
	tempTasks := []models.Task{}
	// As the user model stores a task struct, and not TaskResponse, we need to create
	// another variable so we can return the TaskResponse
	resTasks := []models.TaskResponse{}
	err := DB.Where("user_id = ?", ID).First(&tempTasks).Error
	if err != nil {
		return resTasks, err
	}
	// If no error, we can copy the tasks into the resTasks. The copier function handles this for us
	copier.Copy(&resTasks, &tempTasks)
	return resTasks, nil

}
