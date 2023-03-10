package repo

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/copier"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"todo/models"
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
// This function is under loginDatabase as it deals with the user side of things
func ReturnTasksWithID(ID uint) ([]models.TaskResponse, error) {
	tempTasks := []models.Task{}
	// As the user model stores a task struct, and not TaskResponse, we need to create
	// another variable so we can return the TaskResponse
	resTasks := []models.TaskResponse{}
	err := DB.Where("user_id = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return resTasks, err
	}
	// If no error, we can copy the tasks into the resTasks. The copier function handles this for us
	copier.Copy(&resTasks, &tempTasks)
	return resTasks, nil

}

// This function will add the task to the database by updating the task array
func AddTask(task models.Task) (uint, error) {
	err := DB.Create(&task).Error
	if err != nil {
		// Return 0 for id as if there was an error for creating, there will be no id to return
		return 0, err
	} else {
		return task.ID, nil
	}
}

// This function will delete the task from the database
func DelTask(ID uint) error {
	tempTasks := models.Task{}
	err := DB.Delete(&tempTasks, ID).Error
	return err
}

// Function will change the isDone field on the task corresponding to the ID
func MarkTaskDone(ID uint) error {
	tempTasks := models.Task{}
	err := DB.Where("ID = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return err
	}
	// Change the boolean to inverse as user might need to mark a task as not done
	tempTasks.IsDone = !tempTasks.IsDone
	err = DB.Save(&tempTasks).Error
	return err

}

// This function will edit the task name
func EditTask(ID uint, NewName string) error {
	tempTasks := models.Task{}
	err := DB.Where("ID = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return err
	}
	tempTasks.TaskName = NewName
	err = DB.Save(&tempTasks).Error
	return err
}
