package repo

import (
	"todo/app/models"
	"todo/app/utils/password"

	"gorm.io/gorm"
)

// UserRepo is a repository for interacting with users in the database
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new instance of UserRepo
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// This function will simply add given credentials to the database.
// Does not check for duplicates or anything else.
func (ur UserRepo) AddUser(userInfo models.UserDTO) error {
	// Add task to database by making a model of user type
	userCreate := models.User{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Password: userInfo.Password,
	}
	err := ur.db.Create(&userCreate)
	return err.Error

}

// This function will check if a user exists in the databse or not.
// It will produce an error if user does not exist within the db, otherwise
// nil error if user does exist
func (ur UserRepo) UserExists(userInfo *models.UserDTO, email string) error {
	var userCreate models.User
	err := ur.db.Where("email = ?", email).First(&userCreate).Error
	// We save the credentials in the userInfo variable so we can use it outside the function
	userInfo.Email = userCreate.Email
	userInfo.Password = userCreate.Password
	return err
}

// This function will check if the user exists in the database or not and
// checks if the passwords and emails match
// This will allow us to authenticate logins
func (ur UserRepo) AuthenticateUser(userInfo *models.UserDTO, userInfo2 *models.UserDTO) bool {
	// Here we can compare the hashed password obtained from the database and the password entered by the user
	// First password is plaintext
	isSame, err := password.Verify(userInfo2.Password, userInfo.Password)
	// Also check if emails match
	if userInfo.Email != userInfo2.Email {
		return false
	}
	// If we get no errors above, it means that the two passwords match and we can login the user
	if err == nil && isSame {
		return true
	}
	return false
}

// This function returns the user ID. Assume user already exists within the db.
// If not, it will return 0 as ID
func (ur UserRepo) ReturnUserID(userInfo models.UserDTO) uint {
	var tempUser models.User
	ur.db.Where("email = ?", userInfo.Email).First(&tempUser)
	return tempUser.ID
}
