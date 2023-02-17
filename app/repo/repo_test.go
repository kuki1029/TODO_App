package repo

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"time"

	"testing"

	"todo/app/models"
	"todo/app/utils/password"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Contains all the items required for testing. This just makes it easier when passing variables through functions
// and also ties into the testify suite package
type RepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

// This is needed to redefine the match command for time in the gorm model
type anyTime struct{}

// This makes it so that sqlmock only checks if the time in the db is of time type as
// it is hard to check if the actual time matches due to tiny differences
func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// Setup the mock connection
func (s *RepoTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	s.NoError(err)
	s.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	s.NoError(err)

}

// ==================================================== User Repo Tests ====================================================

// Test AddUser function. This will add a user to the db using the repo functions and
// make sure no errors are returned and the user is succesfully added to the db
func (s *RepoTestSuite) TestAddUser() {
	// Pass the mock db as the new db
	repo := NewUserRepo(s.db)
	// User details
	name := "test"
	email := "user@gmail.com"
	pass := "password"
	// Here we have the expected raw sql that is done by gorm
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"users\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"name\",\"email\",\"password\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, name, email, pass).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()
	// Create DTO model to transfer the user details into the AddUser function
	userCreate := models.UserDTO{
		Name:     name,
		Email:    email,
		Password: pass,
	}
	// Call the actual function that is being tested.
	err := repo.AddUser(userCreate)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the UserExists function when the user does not exist. This attempts to find a user
// that does not exist within the db. This should return an error as the user isn't in the db
func (s *RepoTestSuite) TestUserExistNone() {
	// Pass the mock db as the new db
	repo := NewUserRepo(s.db)
	// User details
	email := "user1@gmail.com"
	user_found := &models.UserDTO{}
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	// Call the actual function that is being tested.
	err := repo.UserExists(user_found, email)
	// As record does not exist within db, we expect this error
	s.Error(err, string(gorm.ErrRecordNotFound.Error()))
	s.Equal(user_found.Email, "")
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the UserExists function when the user does exist. This attempts to find a user
// that does  exist within the db. This should not return an error as the user is in the db
func (s *RepoTestSuite) TestUserExistsInDb() {
	// Pass the mock db as the new db
	repo := NewUserRepo(s.db)
	// User details
	name := "test"
	email := "user2@gmail.com"
	pass := "password"
	// As we also want to return the email and pass, we need to create custom Newrows for that
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).FromCSVString("1, user2@gmail.com, password")
	// Here we have the expected raw sql that is done by gorm
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"users\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"name\",\"email\",\"password\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, name, email, pass).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(rows)
	// Create DTO model to transfer the user details into the AddUser function
	userCreate := models.UserDTO{
		Name:     name,
		Email:    email,
		Password: pass,
	}
	// Call the actual function that is being tested. We need to call the AddUser func first to properly add
	// the user to the db before we can search for them
	err_create := repo.AddUser(userCreate)
	err_find := repo.UserExists(&userCreate, email)
	// As record does not exist within db, we expect this error
	s.NoError(err_create)
	s.NoError(err_find)
	s.Equal(email, userCreate.Email)
	s.Equal(pass, userCreate.Password)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the AuthenticateUser function. This will test the case when the users credentails are valid
// and returns true
func (s *RepoTestSuite) TestAuthenticateUserValid() {
	// Pass the mock db as the new db
	repo := NewUserRepo(s.db)
	// User details
	user1 := &models.UserDTO{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "password",
	}
	// Second user struct should contain a hashed password
	user2 := &models.UserDTO{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: password.Generate("password"),
	}
	// Call the actual function that is being tested.
	isSame := repo.AuthenticateUser(user2, user1)
	// As record does not exist within db, we expect this error
	s.True(isSame)
}

// Tests the AuthenticateUser function. This will test the case when the users credentails are NOT valid
// due to wrong password and returns false
func (s *RepoTestSuite) TestAuthenticateUserWrongPassword() {
	// Pass the mock db as the new db
	repo := NewUserRepo(s.db)
	// User details
	user1 := &models.UserDTO{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "wrongPassword",
	}
	// Second user struct should contain a hashed password
	user2 := &models.UserDTO{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: password.Generate("password"),
	}
	// Call the actual function that is being tested.
	isSame := repo.AuthenticateUser(user2, user1)
	// As record does not exist within db, we expect this error
	s.False(isSame)
}

// Tests the AuthenticateUser function. This will test the case when the users credentails are NOT valid
// due to wrong email and returns false
func (s *RepoTestSuite) TestAuthenticateUserWrongEmail() {
	// Pass the mock db as the new db
	repo := NewUserRepo(s.db)
	// User details
	user1 := &models.UserDTO{
		Name:     "test",
		Email:    "wrong@gmail.com",
		Password: "password",
	}
	// Second user struct should contain a hashed password
	user2 := &models.UserDTO{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: password.Generate("password"),
	}
	// Call the actual function that is being tested.
	isSame := repo.AuthenticateUser(user2, user1)
	// As record does not exist within db, we expect this error
	s.False(isSame)
}

// Tests the ReturnUserId function. This will test the case when the user doesn't exist within db and doesn't have a valid ID
// Should return zero
func (s *RepoTestSuite) TestReturnUserIDDoesNotExist() {
	// Pass the mock db as the new db
	repo := NewUserRepo(s.db)
	// User details
	name := "test"
	email := "user4@gmail.com"
	pass := "password"
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	// Create DTO model to transfer the user details into the AddUser function
	userCreate := models.UserDTO{
		Name:     name,
		Email:    email,
		Password: pass,
	}
	// Call the actual function that is being tested.
	ID := repo.ReturnUserID(userCreate)
	// As record does not exist within db, we expect this error
	s.Equal(ID, uint(0))
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the ReturnUserId function. This will test the case when the user exists within db and has a valid ID
// Should return a positive non-zero number
func (s *RepoTestSuite) TestReturnUserIDValid() {
	// Pass the mock db as the new db
	repo := NewUserRepo(s.db)
	// User details
	name := "test"
	email := "user3@gmail.com"
	pass := "password"
	// Here we have the expected raw sql that is done by gorm
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"users\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"name\",\"email\",\"password\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, name, email, pass).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	// Create DTO model to transfer the user details into the AddUser function
	userCreate := models.UserDTO{
		Name:     name,
		Email:    email,
		Password: pass,
	}
	// Call the actual function that is being tested. We need to call the AddUser func first to properly add
	// the user to the db before we can search for them
	err_create := repo.AddUser(userCreate)
	ID := repo.ReturnUserID(userCreate)
	// As record does not exist within db, we expect this error
	s.NoError(err_create)
	s.Greater(ID, uint(0))
	s.NoError(s.mock.ExpectationsWereMet())
}

// ==================================================== Task Repo Tests ====================================================

// Tests the ReturnTasksWithID function. This will test the case when the tasks for an user is empty. So the returned array
// should be empty
func (s *RepoTestSuite) TestReturnTasksWithIDEmpty() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	// User details
	ID := uint(1)
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE user_id = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
		WillReturnError(gorm.ErrRecordNotFound)
	// Call the actual function that is being tested.
	tasks, err := repo.ReturnTasksWithID(ID)
	// As record does not exist within db, we expect this error
	s.Error(err, string(gorm.ErrRecordNotFound.Error()))
	s.Empty(tasks)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the ReturnTasksWithID function. This will test the case when the tasks for an user isn't empty. So the returned array
// should contain some array of tasks
func (s *RepoTestSuite) TestReturnTasksWithIDNotEmpty() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	// User details
	ID := uint(1)
	// As we also want to return the email and pass, we need to create custom Newrows for that
	rows := sqlmock.NewRows([]string{"id", "TaskName", "isDone"}).FromCSVString("1, Task1, false")
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE user_id = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(ID).
		WillReturnRows(rows)
	// Call the actual function that is being tested.
	tasks, err := repo.ReturnTasksWithID(ID)
	// As record does not exist within db, we expect this error
	s.NoError(err)
	s.Equal("Task1", tasks[0].TaskName)
	s.Equal(false, tasks[0].IsDone)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the AddTask function. This will simply add a task to the db and expect no error
// This function shouldn't return an error unless there is something wrong with the db itself
func (s *RepoTestSuite) TestAddTaskSimple() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	// User details
	name := "Task1"
	isDone := false
	userID := 4
	taskID_expected := uint(6)
	// Here we have the expected raw sql that is done by gorm
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"tasks\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"task_name\",\"is_done\",\"user_id\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, name, isDone, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(taskID_expected))
	s.mock.ExpectCommit()
	// Create model for it
	task := models.Task{
		TaskName: name,
		IsDone:   false,
		UserID:   4,
	}
	// Call the actual function that is being tested.
	taskID, err := repo.AddTask(task)
	// As record does not exist within db, we expect this error
	s.NoError(err)
	s.Equal(taskID_expected, taskID)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the DelTask function. This will simply delete a task from the db.
// This function won't return as error as it is only called on tasks that already exist.
// Also the gorm db won't throw errors for deleting something that doesn't exist
func (s *RepoTestSuite) TestDelTaskSimple() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	taskID_expected := uint(6)
	// Here we have the expected raw sql that is done by gorm
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "deleted_at"=$1 WHERE "tasks"."id" = $2 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(anyTime{}, taskID_expected).
		WillReturnResult(sqlmock.NewResult(int64(taskID_expected), 1))
	s.mock.ExpectCommit()
	// Call the actual function that is being tested.
	err := repo.DelTask(taskID_expected)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the MarkTaskDone function. This will simply change the isDone field
// on the task. If the task doesn't exist, it will throw an error but that shouldn't happen
// as we only call this on existing tasks. This will check the case when the task does exist
func (s *RepoTestSuite) TestMarkTaskDone() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	taskID_expected := uint(6)
	isDone := true
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE ID = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(taskID_expected).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(taskID_expected))
	// Here we have the expected raw sql that is done by gorm
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "is_done"=$1,"updated_at"=$2 WHERE "tasks"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(isDone, anyTime{}, taskID_expected).
		WillReturnResult(sqlmock.NewResult(int64(taskID_expected), 1))
	s.mock.ExpectCommit()
	// Call the actual function that is being tested.
	err := repo.MarkTaskDone(taskID_expected)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the MarkTaskDone function. This will simply change the isDone field
// on the task. If the task doesn't exist, it will throw an error but that shouldn't happen
// as we only call this on existing tasks. This will check the case when the task doesn't exist
func (s *RepoTestSuite) TestMarkTaskDoneError() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	taskID_expected := uint(6)
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE ID = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(taskID_expected).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(taskID_expected)).
		WillReturnError(gorm.ErrRecordNotFound)
	// Call the actual function that is being tested.
	err := repo.MarkTaskDone(taskID_expected)
	s.Error(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the EditTask function. This will simply change the task_name field
// on the task. If the task doesn't exist, it will throw an error but that shouldn't happen
// as we only call this on existing tasks. This will check the case when the task does exist
func (s *RepoTestSuite) TestEditTask() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	taskID_expected := uint(6)
	name := "New Name"
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE ID = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(taskID_expected).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(taskID_expected))
	// Here we have the expected raw sql that is done by gorm
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "task_name"=$1,"updated_at"=$2 WHERE "tasks"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(name, anyTime{}, taskID_expected).
		WillReturnResult(sqlmock.NewResult(int64(taskID_expected), 1))
	s.mock.ExpectCommit()
	// Call the actual function that is being tested.
	err := repo.EditTask(taskID_expected, name)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the EditTask function. This will simply change the task_name field
// on the task. If the task doesn't exist, it will throw an error but that shouldn't happen
// as we only call this on existing tasks. This will check the case when the task doesn't exist
func (s *RepoTestSuite) TestEditTaskError() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	taskID_expected := uint(6)
	name := "New Name"
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE ID = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(taskID_expected).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(taskID_expected)).
		WillReturnError(gorm.ErrRecordNotFound)
	// Call the actual function that is being tested.
	err := repo.EditTask(taskID_expected, name)
	s.Error(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the ReturnName function. This will return the name of the user.
// This case is for when the user exists and the function doesn't throw an error
func (s *RepoTestSuite) TestReturnName() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	userID := uint(1)
	expectedName := "User"
	// As we also want to return the email and pass, we need to create custom Newrows for that
	rows := sqlmock.NewRows([]string{"id", "name"}).FromCSVString("1, User")
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE ID = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(userID).
		WillReturnRows(rows)
	// Call the actual function that is being tested.
	name, err := repo.ReturnName(userID)
	s.NoError(err)
	s.Equal(expectedName, name)
	s.NoError(s.mock.ExpectationsWereMet())
}

// Tests the ReturnName function. This will return the name of the user.
// This case is for when the user doesn't exist so the function returns an error and blank name
func (s *RepoTestSuite) TestReturnNameError() {
	// Pass the mock db as the new db
	repo := NewTaskRepo(s.db)
	userID := uint(1)
	expectedName := ""
	// As we also want to return the email and pass, we need to create custom Newrows for that
	rows := sqlmock.NewRows([]string{"id"}).FromCSVString("1")
	// Set the expected query from the func
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE ID = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(userID).
		WillReturnRows(rows).
		WillReturnError(gorm.ErrRecordNotFound)
	// Call the actual function that is being tested.
	name, err := repo.ReturnName(userID)
	s.Error(err)
	s.Equal(expectedName, name)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
