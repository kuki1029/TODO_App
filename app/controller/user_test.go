package controller

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
	"todo/app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

// First we need to setup mock functions for all our db functions
// This will allow us to control the returns for these functions and make it
// so that we don't need to mock the db directly
type mockUserRepo struct {
	mock.Mock
}

type mockRedisClient struct {
	mock.Mock
}

// ==================================== Mock Functions Setup ====================================
func (m *mockUserRepo) AddUser(userInfo models.UserDTO) error {
	args := m.Called(userInfo)
	return args.Error(0)
}

func (m *mockUserRepo) UserExists(userInfo *models.UserDTO, email string) error {
	args := m.Called(userInfo, email)
	return args.Error(0)
}

func (m *mockUserRepo) AuthenticateUser(userInfo *models.UserDTO, userInfo2 *models.UserDTO) bool {
	args := m.Called(userInfo, userInfo2)
	return args.Get(0).(bool)
}

func (m *mockUserRepo) ReturnUserID(userInfo models.UserDTO) uint {
	args := m.Called(userInfo)
	return args.Get(0).(uint)
}

func (m *mockRedisClient) GetFromRedis(key string) (uint, error) {
	args := m.Called(key)
	return args.Get(0).(uint), args.Error(1)
}

func (m *mockRedisClient) SetInRedis(value string, redisVal string, timeAmt time.Duration) {
	_ = m.Called(value, redisVal, timeAmt)
	return
}

func (m *mockRedisClient) DelInRedis(value string) {
	_ = m.Called(value)
	return
}

func (m *mockRedisClient) Ping() error {
	args := m.Called()
	return args.Error(0)
}

// ==================================== Mock Functions Setup End ====================================

type UserControllerTestSuite struct {
	suite.Suite
	controller *UserController
	mockRepo   *mockUserRepo
	mockRedis  *mockRedisClient
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.mockRepo = &mockUserRepo{}
	suite.mockRedis = &mockRedisClient{}
	suite.controller = &UserController{
		userRepo:    suite.mockRepo,
		redisClient: suite.mockRedis,
	}
}

// This tests the login function with basic functionality
// Tests the case when everything is okay and the function returns no errors
func (suite *UserControllerTestSuite) TestLogin() {
	user1 := models.UserDTO{
		Name:     "",
		Email:    "test@gmail.com",
		Password: "password",
	}
	// Set up expected behavior for the mock repository
	userTemp := new(models.UserDTO)
	// The .Run is used as we modify a pointer in the UserExists function.
	// To mock this, we need to use .RUn
	suite.mockRepo.On("UserExists", userTemp, "test@gmail.com").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.UserDTO)
		arg.Name = "Test"
		arg.Email = "test@gmail.com"
		arg.Password = "hashedPassword"
	})
	// Instead of actually hashing the passwords, we just type it in ourselves as the hashedPassword
	// to reduce any external dependencies
	userTempReturned := models.UserDTO{
		Name:     "Test",
		Email:    "test@gmail.com",
		Password: "hashedPassword",
	}
	suite.mockRepo.On("AuthenticateUser", &userTempReturned, &user1).Return(true)
	suite.mockRepo.On("ReturnUserID", user1).Return(uint(1))
	suite.mockRedis.On("SetInRedis", mock.Anything, "email: test@gmail.com id: 1", 24*time.Hour).Return()

	// Create a new Fiber test context
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	// Set request and set the body
	c.Request().SetRequestURI("/login")
	c.Request().Header.SetMethod("POST")
	c.Request().Header.Set("Content-Type", "application/json")
	c.Context().Request.SetBodyString(`{"email":"test@gmail.com","password":"password"}`)
	// Call the Login function
	err := suite.controller.Login(c)
	// Make sure everything went okay
	suite.NoError(err)
	suite.Equal(http.StatusOK, c.Response().StatusCode())
	// Assert that the CreateTeam method of the mock repository was called with the correct arguments
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockRedis.AssertExpectations(suite.T())
	// Parse the JSON response
	var response map[string]interface{}
	err = json.Unmarshal(
		c.Response().Body(), &response)
	suite.NoError(err)
	suite.Equal(true, response["success"])

}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}