package repo

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"todo/models"
)

// This function returns the tasks for a particular user by looking them up through their ID
// This function is under loginDatabase as it deals with the user side of things
func ReturnTasksWithID(ID uint, db *gorm.DB) ([]models.TaskResponse, error) {
	tempTasks := []models.Task{}
	// As the user model stores a task struct, and not TaskResponse, we need to create
	// another variable so we can return the TaskResponse
	resTasks := []models.TaskResponse{}
	err := db.Where("user_id = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return resTasks, err
	}
	// If no error, we can copy the tasks into the resTasks. The copier function handles this for us
	copier.Copy(&resTasks, &tempTasks)
	return resTasks, nil

}

// This function will add the task to the database by updating the task array
func AddTask(task models.Task, db *gorm.DB) (uint, error) {
	err := db.Create(&task).Error
	if err != nil {
		// Return 0 for id as if there was an error for creating, there will be no id to return
		return 0, err
	} else {
		return task.ID, nil
	}
}

// This function will delete the task from the database
func DelTask(ID uint, db *gorm.DB) error {
	tempTasks := models.Task{}
	err := db.Delete(&tempTasks, ID).Error
	return err
}

// Function will change the isDone field on the task corresponding to the ID
func MarkTaskDone(ID uint, db *gorm.DB) error {
	tempTasks := models.Task{}
	err := db.Where("ID = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return err
	}
	// Change the boolean to inverse as user might need to mark a task as not done
	tempTasks.IsDone = !tempTasks.IsDone
	err = db.Save(&tempTasks).Error
	return err

}

// This function will edit the task name
func EditTask(ID uint, NewName string, db *gorm.DB) error {
	tempTasks := models.Task{}
	err := db.Where("ID = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return err
	}
	tempTasks.TaskName = NewName
	err = db.Save(&tempTasks).Error
	return err
}
