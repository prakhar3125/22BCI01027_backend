package main
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the auth service
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(email, password string) (int, error) {
	args := m.Called(email, password)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

// TestRegisterEndpoint tests the register endpoint
func TestRegisterEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock auth service
	mockAuthService := new(MockAuthService)
	mockAuthService.On("Register", "test@example.com", "password123").Return(1, nil)

	// Create auth controller with mock service
	authController := NewAuthController(mockAuthService)

	// Create a new router
	router := gin.Default()
	router.POST("/register", authController.Register)

	// Create a request body
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})

	// Create a request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Check response data
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "User registered successfully", response["message"])

	// Verify that the mock was called
	mockAuthService.AssertExpectations(t)
}