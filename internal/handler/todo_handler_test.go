package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VoLKoDaV13579/GoToDoList/internal/handler/mocks"
	"github.com/VoLKoDaV13579/GoToDoList/internal/model"
	"github.com/VoLKoDaV13579/GoToDoList/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewTodoHandler(t *testing.T) {
	mockService := new(mocks.MockTodoService)
	v := validator.New()
	handler := NewTodoHandler(mockService, v)
	assert.NotNil(t, handler)
}

func TestTodoHandler_CreateTodo(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.POST("/todos", handler.CreateTodo)

		reqBody := model.CreateTodoRequest{
			Title:       "Test Todo",
			Description: "Test Description",
			Priority:    1,
		}

		expectedTodo := &model.Todo{
			ID:          1,
			Title:       "Test Todo",
			Description: "Test Description",
			Status:      model.StatusPending,
			Priority:    1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockService.On("CreateTodo", mock.AnythingOfType("*model.CreateTodoRequest")).
			Return(expectedTodo, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response model.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Test Todo", response.Title)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.POST("/todos", handler.CreateTodo)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("validation error - empty title", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.POST("/todos", handler.CreateTodo)

		reqBody := model.CreateTodoRequest{
			Title:    "", // Empty title
			Priority: 1,
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation_error", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.POST("/todos", handler.CreateTodo)

		reqBody := model.CreateTodoRequest{
			Title:    "Test Todo",
			Priority: 1,
		}

		mockService.On("CreateTodo", mock.AnythingOfType("*model.CreateTodoRequest")).
			Return(nil, errors.New("database error"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTodoHandler_GetTodoByID(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.GET("/todos/:id", handler.GetTodoByID)

		expectedTodo := &model.Todo{
			ID:          1,
			Title:       "Test Todo",
			Description: "Test Description",
			Status:      model.StatusPending,
			Priority:    1,
		}

		mockService.On("GetTodoByID", uint(1)).Return(expectedTodo, nil)

		req, _ := http.NewRequest("GET", "/todos/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid ID format", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.GET("/todos/:id", handler.GetTodoByID)

		req, _ := http.NewRequest("GET", "/todos/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_id", response.Error)
	})

	t.Run("todo not found", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.GET("/todos/:id", handler.GetTodoByID)

		mockService.On("GetTodoByID", uint(999)).
			Return(nil, errors.New("todo not found"))

		req, _ := http.NewRequest("GET", "/todos/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.GET("/todos/:id", handler.GetTodoByID)

		mockService.On("GetTodoByID", uint(1)).
			Return(nil, errors.New("database error"))

		req, _ := http.NewRequest("GET", "/todos/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTodoHandler_GetAllTodos(t *testing.T) {
	t.Run("successful get all", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.GET("/todos", handler.GetAllTodos)

		expectedResponse := &model.TodoListResponse{
			Data: []model.Todo{
				{ID: 1, Title: "Todo 1", Status: model.StatusPending},
				{ID: 2, Title: "Todo 2", Status: model.StatusCompleted},
			},
			Total:      2,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("GetAllTodos", mock.AnythingOfType("*model.TodoFilter")).
			Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/todos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.TodoListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Data))
		assert.Equal(t, int64(2), response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("get with filters", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.GET("/todos", handler.GetAllTodos)

		expectedResponse := &model.TodoListResponse{
			Data: []model.Todo{
				{ID: 1, Title: "Todo 1", Status: model.StatusCompleted},
			},
			Total:      1,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		mockService.On("GetAllTodos", mock.AnythingOfType("*model.TodoFilter")).
			Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/todos?status=completed&priority=2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.GET("/todos", handler.GetAllTodos)

		mockService.On("GetAllTodos", mock.AnythingOfType("*model.TodoFilter")).
			Return(nil, errors.New("database error"))

		req, _ := http.NewRequest("GET", "/todos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTodoHandler_UpdateTodo(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id", handler.UpdateTodo)

		newTitle := "Updated Title"
		reqBody := model.UpdateTodoRequest{
			Title: &newTitle,
		}

		updatedTodo := &model.Todo{
			ID:     1,
			Title:  "Updated Title",
			Status: model.StatusPending,
		}

		mockService.On("UpdateTodo", uint(1), mock.AnythingOfType("*model.UpdateTodoRequest")).
			Return(updatedTodo, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", response.Title)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid ID format", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id", handler.UpdateTodo)

		newTitle := "Updated Title"
		reqBody := model.UpdateTodoRequest{
			Title: &newTitle,
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/abc", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id", handler.UpdateTodo)

		req, _ := http.NewRequest("PATCH", "/todos/1", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("todo not found", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id", handler.UpdateTodo)

		newTitle := "Updated Title"
		reqBody := model.UpdateTodoRequest{
			Title: &newTitle,
		}

		mockService.On("UpdateTodo", uint(999), mock.AnythingOfType("*model.UpdateTodoRequest")).
			Return(nil, errors.New("todo not found"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id", handler.UpdateTodo)

		newTitle := "Updated Title"
		reqBody := model.UpdateTodoRequest{
			Title: &newTitle,
		}

		mockService.On("UpdateTodo", uint(1), mock.AnythingOfType("*model.UpdateTodoRequest")).
			Return(nil, errors.New("database error"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTodoHandler_DeleteTodo(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.DELETE("/todos/:id", handler.DeleteTodo)

		mockService.On("DeleteTodo", uint(1)).Return(nil)

		req, _ := http.NewRequest("DELETE", "/todos/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid ID format", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.DELETE("/todos/:id", handler.DeleteTodo)

		req, _ := http.NewRequest("DELETE", "/todos/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("todo not found", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.DELETE("/todos/:id", handler.DeleteTodo)

		mockService.On("DeleteTodo", uint(999)).
			Return(errors.New("todo not found"))

		req, _ := http.NewRequest("DELETE", "/todos/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.DELETE("/todos/:id", handler.DeleteTodo)

		mockService.On("DeleteTodo", uint(1)).
			Return(errors.New("database error"))

		req, _ := http.NewRequest("DELETE", "/todos/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTodoHandler_UpdateTodoStatus(t *testing.T) {
	t.Run("successful status update", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id/status", handler.UpdateTodoStatus)

		reqBody := map[string]string{
			"status": "completed",
		}

		updatedTodo := &model.Todo{
			ID:     1,
			Title:  "Test Todo",
			Status: model.StatusCompleted,
		}

		mockService.On("UpdateTodoStatus", uint(1), model.StatusCompleted).
			Return(updatedTodo, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/1/status", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, model.StatusCompleted, response.Status)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid ID format", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id/status", handler.UpdateTodoStatus)

		reqBody := map[string]string{
			"status": "completed",
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/abc/status", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id/status", handler.UpdateTodoStatus)

		req, _ := http.NewRequest("PATCH", "/todos/1/status", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid status", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id/status", handler.UpdateTodoStatus)

		reqBody := map[string]string{
			"status": "invalid_status",
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/1/status", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_status", response.Error)
	})

	t.Run("todo not found", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id/status", handler.UpdateTodoStatus)

		reqBody := map[string]string{
			"status": "completed",
		}

		mockService.On("UpdateTodoStatus", uint(999), model.StatusCompleted).
			Return(nil, errors.New("todo not found"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/999/status", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockTodoService)
		v := validator.New()
		handler := NewTodoHandler(mockService, v)

		router := setupTestRouter()
		router.PATCH("/todos/:id/status", handler.UpdateTodoStatus)

		reqBody := map[string]string{
			"status": "completed",
		}

		mockService.On("UpdateTodoStatus", uint(1), model.StatusCompleted).
			Return(nil, errors.New("database error"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/todos/1/status", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTodoHandler_HealthCheck(t *testing.T) {
	mockService := new(mocks.MockTodoService)
	v := validator.New()
	handler := NewTodoHandler(mockService, v)

	router := setupTestRouter()
	router.GET("/health", handler.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}
