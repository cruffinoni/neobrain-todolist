package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/cruffinoni/neobrain-todolist/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoutes struct {
	mock.Mock
	db database.Database
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) DeleteTask(taskID int64) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func (m *MockDB) MarkTaskAsDone(taskID int64) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func (m *MockDB) GetTasks(filter database.TaskFilter) ([]*database.Task, error) {
	args := m.Called(filter)
	return args.Get(0).([]*database.Task), args.Error(1)
}

func (m *MockDB) ImportTasks(tasks []*database.Task) error {
	args := m.Called(tasks)
	return args.Error(0)
}

func (m *MockDB) AddTask(task string) (int64, error) {
	args := m.Called(task)
	return args.Get(0).(int64), args.Error(1)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestAddTask(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		mockDBResponse func(task string) (int64, error)
	}{
		{
			name:           "valid_input",
			payload:        `{"task": "Sample task"}`,
			expectedStatus: http.StatusCreated,
			mockDBResponse: func(task string) (int64, error) {
				return 1, nil
			},
		},
		{
			name:           "empty_task",
			payload:        `{"task": ""}`,
			expectedStatus: http.StatusBadRequest,
			mockDBResponse: func(task string) (int64, error) {
				return 0, nil // Not used in this case
			},
		},
		{
			name:           "invalid_json_input",
			payload:        `{"invalid":}`,
			expectedStatus: http.StatusBadRequest,
			mockDBResponse: func(task string) (int64, error) {
				return 0, nil // Not used in this case
			},
		},
		{
			name:           "db_error",
			payload:        `{"task": "Sample task"}`,
			expectedStatus: http.StatusInternalServerError,
			mockDBResponse: func(task string) (int64, error) {
				return 0, errors.New("database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			var requestBody AddTaskBody
			_ = json.Unmarshal([]byte(tt.payload), &requestBody)
			mockDB.On("AddTask", requestBody.Task).Return(tt.mockDBResponse(requestBody.Task))
			routes := Routes{db: mockDB}
			router := gin.Default()
			router.POST("/tasks", routes.AddTask)

			req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response AddTaskResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), response.TaskID)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		expectedStatus int
		mockDBResponse func(taskID int64) error
	}{
		{
			name:           "valid_task_id",
			taskID:         "1",
			expectedStatus: http.StatusOK,
			mockDBResponse: func(taskID int64) error {
				return nil
			},
		},
		{
			name:           "invalid_task_id",
			taskID:         "abc",
			expectedStatus: http.StatusBadRequest,
			mockDBResponse: func(taskID int64) error {
				return nil // Not used in this case
			},
		},
		{
			name:           "db_error",
			taskID:         "1",
			expectedStatus: http.StatusInternalServerError,
			mockDBResponse: func(taskID int64) error {
				return errors.New("database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			taskID, err := strconv.ParseInt(tt.taskID, 10, 64)
			if err == nil {
				mockDB.On("DeleteTask", taskID).Return(tt.mockDBResponse(taskID))
			}

			routes := Routes{db: mockDB}
			router := gin.Default()
			router.DELETE("/tasks/:task_id", routes.DeleteTask)

			req, _ := http.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

func TestListTasks(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		filter         database.TaskFilter
		expectedStatus int
		mockDBResponse func(filter database.TaskFilter) ([]*database.Task, error)
	}{
		{
			name:           "no_filter",
			queryParam:     "",
			expectedStatus: http.StatusOK,
			filter:         database.TaskFilterAll,
			mockDBResponse: func(filter database.TaskFilter) ([]*database.Task, error) {
				return []*database.Task{
					{ID: 1, Task: "Task 1", Done: false},
					{ID: 2, Task: "Task 2", Done: true},
				}, nil
			},
		},
		{
			name:           "filter_completed",
			queryParam:     "completed=true",
			expectedStatus: http.StatusOK,
			filter:         database.TaskFilterDone,
			mockDBResponse: func(filter database.TaskFilter) ([]*database.Task, error) {
				return []*database.Task{
					{ID: 2, Task: "Task 2", Done: true},
				}, nil
			},
		},
		{
			name:           "filter_not_completed",
			queryParam:     "completed=false",
			filter:         database.TaskFilterNotDone,
			expectedStatus: http.StatusOK,
			mockDBResponse: func(filter database.TaskFilter) ([]*database.Task, error) {
				return []*database.Task{
					{ID: 1, Task: "Task 1", Done: false},
				}, nil
			},
		},
		{
			name:           "db_error",
			queryParam:     "",
			filter:         database.TaskFilterAll,
			expectedStatus: http.StatusInternalServerError,
			mockDBResponse: func(filter database.TaskFilter) ([]*database.Task, error) {
				return nil, errors.New("database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			routes := Routes{db: mockDB}
			mockDB.On("GetTasks", tt.filter).Return(tt.mockDBResponse(tt.filter))
			router := gin.Default()
			router.GET("/tasks", routes.ListTasks)

			req, _ := http.NewRequest("GET", "/tasks?"+tt.queryParam, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if tt.expectedStatus == http.StatusOK {
				var tasks []*database.Task
				err := json.Unmarshal(resp.Body.Bytes(), &tasks)
				assert.NoError(t, err)
				assert.NotEmpty(t, tasks)
			}
		})
	}
}
