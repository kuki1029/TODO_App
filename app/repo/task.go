package repo

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"todo/app/models"
)

// TaskRepo is a repository for interacting with tasks in the database
type TaskRepo struct {
	db *gorm.DB
}

// NewTaskRepo creates a new instance of TaskRepo
func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

// This function returns the tasks for a particular user by looking them up through their ID
// This function is under loginDatabase as it deals with the user side of things
func (tr TaskRepo) ReturnTasksWithID(ID uint) ([]models.TaskResponse, error) {
	tempTasks := []models.Task{}
	// As the user model stores a task struct, and not TaskResponse, we need to create
	// another variable so we can return the TaskResponse
	resTasks := []models.TaskResponse{}
	err := tr.db.Where("user_id = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return resTasks, err
	}
	// If no error, we can copy the tasks into the resTasks. The copier function handles this for us
	copier.Copy(&resTasks, &tempTasks)
	return resTasks, nil

}

// This function will add the task to the database by updating the task array
func (tr TaskRepo) AddTask(task models.Task) (uint, error) {
	err := tr.db.Create(&task).Error
	if err != nil {
		// Return 0 for id as if there was an error for creating, there will be no id to return
		return 0, err
	} else {
		return task.ID, nil
	}
}

// This function will delete the task from the database
func (tr TaskRepo) DelTask(ID uint) error {
	tempTasks := models.Task{}
	err := tr.db.Delete(&tempTasks, ID).Error
	return err
}

// Function will change the isDone field on the task corresponding to the ID
func (tr TaskRepo) MarkTaskDone(ID uint) error {
	tempTasks := models.Task{}
	err := tr.db.Where("ID = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return err
	}
	// Change the boolean to inverse as user might need to mark a task as not done
	tempTasks.IsDone = !tempTasks.IsDone
	err = tr.db.Save(&tempTasks).Error
	return err

}

// This function will edit the task name
func (tr TaskRepo) EditTask(ID uint, NewName string) error {
	tempTasks := models.Task{}
	err := tr.db.Where("ID = ?", ID).Find(&tempTasks).Error
	if err != nil {
		return err
	}
	tempTasks.TaskName = NewName
	err = tr.db.Save(&tempTasks).Error
	return err
}

// This function returns the name for a particular user using their ID
func (tr TaskRepo) ReturnName(ID uint) (string, error) {
	tempUser := models.User{}
	// As the user model stores a task struct, and not TaskResponse, we need to create
	// another variable so we can return the TaskResponse
	err := tr.db.Where("ID = ?", ID).Find(&tempUser).Error

	if err != nil {
		// Return blank name incase of error
		return "", err
	}
	return tempUser.Name, nil
}
