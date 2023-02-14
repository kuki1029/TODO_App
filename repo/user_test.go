package repo

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"testing"

	"todo/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

type anyTime struct{}

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
	}), &gorm.Config{})
	s.NoError(err)

}

// Test AddUser function
func (s *RepoTestSuite) TestCreateUser() {

	repo := NewRepo(s.db)

	name_test := "test"
	email_test := "user@gmail.com"
	password_test := "password"

	// We use anyArg for the password as it is hashed
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"users\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"name\",\"email\",\"password\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, "test", "user@gmail.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	userCreate := models.UserDTO{
		Name:     name_test,
		Email:    email_test,
		Password: password_test,
	}

	err := AddUser(userCreate, repo.db)

	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}



func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
