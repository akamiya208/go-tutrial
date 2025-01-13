package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akamiya208/go-tutrial/internal/pkg/dto"
	"github.com/akamiya208/go-tutrial/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockedMySQLClient struct {
	mock.Mock
}

func (m *MockedMySQLClient) GetTask(taskID uint) (models.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockedMySQLClient) GetTasksByName(name string) ([]models.Task, error) {
	args := m.Called(name)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockedMySQLClient) CreateTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockedMySQLClient) UpdateTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockedMySQLClient) DeleteTask(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockedMySQLClient) DB() *gorm.DB {
	return nil
}

func TestHandleGetTask(t *testing.T) {
	description := "description"
	testTask := models.Task{
		ID:          1,
		Name:        "task1",
		Description: &description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   gorm.DeletedAt{},
	}

	t.Run("ok", func(t *testing.T) {
		taskID := uint(1)
		mockedMySQLClient := new(MockedMySQLClient)
		mockedMySQLClient.On("GetTask", taskID).Return(testTask, nil)
		handlers := NewTaskHandler(mockedMySQLClient)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tasks/%d", taskID), nil)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/tasks/{taskId}", handlers.HandleGetTask)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.TaskResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, "task1", response.Name)
	})

	t.Run("not found", func(t *testing.T) {
		taskID := uint(1)
		mockedMySQLClient := new(MockedMySQLClient)
		mockedMySQLClient.On("GetTask", taskID).Return(models.Task{}, fmt.Errorf("record not found"))
		handlers := NewTaskHandler(mockedMySQLClient)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tasks/%d", taskID), nil)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/tasks/{taskId}", handlers.HandleGetTask)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, "Not Found", response.Status)
	})

	t.Run("internal server error", func(t *testing.T) {
		taskID := uint(1)
		mockedMySQLClient := new(MockedMySQLClient)
		mockedMySQLClient.On("GetTask", taskID).Return(models.Task{}, fmt.Errorf("some error"))
		handlers := NewTaskHandler(mockedMySQLClient)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tasks/%d", taskID), nil)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/tasks/{taskId}", handlers.HandleGetTask)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, "Internal Server Error", response.Status)
	})

	t.Run("bad request", func(t *testing.T) {
		taskID := "invalid"
		mockedMySQLClient := new(MockedMySQLClient)
		handlers := NewTaskHandler(mockedMySQLClient)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tasks/%s", taskID), nil)

		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/v1/tasks/{taskId}", handlers.HandleGetTask)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, "Bad Request", response.Status)
	})
}
