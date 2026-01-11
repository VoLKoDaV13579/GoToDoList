package service

import (
	"errors"
	"testing"
	"time"

	"github.com/VoLKoDaV13579/GoToDoList/internal/model"
	"github.com/VoLKoDaV13579/GoToDoList/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTodoService(t *testing.T) {
	mockRepo := new(mocks.MockTodoRepository)
	service := NewTodoService(mockRepo)
	assert.NotNil(t, service)
}

func TestTodoService_CreateTodo(t *testing.T) {
	t.Run("successful creation with default status", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		req := &model.CreateTodoRequest{
			Title:       "Test Todo",
			Description: "Test Description",
			Priority:    1,
		}

		mockRepo.On("Create", mock.AnythingOfType("*model.Todo")).
			Return(nil).
			Run(func(args mock.Arguments) {
				todo := args.Get(0).(*model.Todo)
				todo.ID = 1
			})

		todo, err := service.CreateTodo(req)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, "Test Todo", todo.Title)
		assert.Equal(t, model.StatusPending, todo.Status) // Default status
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful creation with custom status", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		req := &model.CreateTodoRequest{
			Title:       "Test Todo",
			Description: "Test Description",
			Status:      model.StatusInProgress,
			Priority:    2,
		}

		mockRepo.On("Create", mock.AnythingOfType("*model.Todo")).Return(nil)

		todo, err := service.CreateTodo(req)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, model.StatusInProgress, todo.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("creation with due date", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		dueDate := time.Now().Add(24 * time.Hour)
		req := &model.CreateTodoRequest{
			Title:    "Test Todo",
			Priority: 0,
			DueDate:  &dueDate,
		}

		mockRepo.On("Create", mock.AnythingOfType("*model.Todo")).Return(nil)

		todo, err := service.CreateTodo(req)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.NotNil(t, todo.DueDate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("creation error", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		req := &model.CreateTodoRequest{
			Title:    "Test Todo",
			Priority: 1,
		}

		mockRepo.On("Create", mock.AnythingOfType("*model.Todo")).
			Return(errors.New("database error"))

		todo, err := service.CreateTodo(req)

		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "failed to create todo")
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_GetTodoByID(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		expectedTodo := &model.Todo{
			ID:          1,
			Title:       "Test Todo",
			Description: "Test Description",
			Status:      model.StatusPending,
			Priority:    1,
		}

		mockRepo.On("GetByID", uint(1)).Return(expectedTodo, nil)

		todo, err := service.GetTodoByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, uint(1), todo.ID)
		assert.Equal(t, "Test Todo", todo.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("todo not found", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		mockRepo.On("GetByID", uint(999)).
			Return(nil, errors.New("todo not found"))

		todo, err := service.GetTodoByID(999)

		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "todo not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_GetAllTodos(t *testing.T) {
	t.Run("successful get with default pagination", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		filter := &model.TodoFilter{
			Page:     0, // Should default to 1
			PageSize: 0, // Should default to 10
		}

		todos := []model.Todo{
			{ID: 1, Title: "Todo 1", Status: model.StatusPending},
			{ID: 2, Title: "Todo 2", Status: model.StatusCompleted},
		}

		mockRepo.On("GetAll", mock.AnythingOfType("*model.TodoFilter")).
			Return(todos, int64(2), nil)

		response, err := service.GetAllTodos(filter)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Data))
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		assert.Equal(t, 1, response.TotalPages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get with custom pagination", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		filter := &model.TodoFilter{
			Page:     2,
			PageSize: 5,
		}

		todos := []model.Todo{
			{ID: 6, Title: "Todo 6", Status: model.StatusPending},
		}

		mockRepo.On("GetAll", filter).Return(todos, int64(11), nil)

		response, err := service.GetAllTodos(filter)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, int64(11), response.Total)
		assert.Equal(t, 2, response.Page)
		assert.Equal(t, 5, response.PageSize)
		assert.Equal(t, 3, response.TotalPages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get with status filter", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		status := model.StatusCompleted
		filter := &model.TodoFilter{
			Status:   &status,
			Page:     1,
			PageSize: 10,
		}

		todos := []model.Todo{
			{ID: 1, Title: "Todo 1", Status: model.StatusCompleted},
		}

		mockRepo.On("GetAll", filter).Return(todos, int64(1), nil)

		response, err := service.GetAllTodos(filter)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Data))
		mockRepo.AssertExpectations(t)
	})

	t.Run("page size exceeds maximum", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		filter := &model.TodoFilter{
			Page:     1,
			PageSize: 200, // Should be capped to 100
		}

		mockRepo.On("GetAll", mock.AnythingOfType("*model.TodoFilter")).
			Return([]model.Todo{}, int64(0), nil)

		response, err := service.GetAllTodos(filter)

		assert.NoError(t, err)
		assert.Equal(t, 100, response.PageSize) // Capped at 100
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		filter := &model.TodoFilter{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("GetAll", mock.AnythingOfType("*model.TodoFilter")).
			Return(nil, int64(0), errors.New("database error"))

		response, err := service.GetAllTodos(filter)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to get todos")
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_UpdateTodo(t *testing.T) {
	t.Run("successful update all fields", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		existingTodo := &model.Todo{
			ID:          1,
			Title:       "Old Title",
			Description: "Old Description",
			Status:      model.StatusPending,
			Priority:    0,
		}

		newTitle := "New Title"
		newDesc := "New Description"
		newStatus := model.StatusCompleted
		newPriority := 2
		newDueDate := time.Now().Add(24 * time.Hour)

		req := &model.UpdateTodoRequest{
			Title:       &newTitle,
			Description: &newDesc,
			Status:      &newStatus,
			Priority:    &newPriority,
			DueDate:     &newDueDate,
		}

		mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.Todo")).Return(nil)

		todo, err := service.UpdateTodo(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, "New Title", todo.Title)
		assert.Equal(t, "New Description", todo.Description)
		assert.Equal(t, model.StatusCompleted, todo.Status)
		assert.Equal(t, 2, todo.Priority)
		assert.NotNil(t, todo.DueDate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("partial update", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		existingTodo := &model.Todo{
			ID:          1,
			Title:       "Old Title",
			Description: "Old Description",
			Status:      model.StatusPending,
			Priority:    1,
		}

		newTitle := "New Title"
		req := &model.UpdateTodoRequest{
			Title: &newTitle,
		}

		mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.Todo")).Return(nil)

		todo, err := service.UpdateTodo(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, "New Title", todo.Title)
		assert.Equal(t, "Old Description", todo.Description) // Unchanged
		mockRepo.AssertExpectations(t)
	})

	t.Run("todo not found", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		newTitle := "New Title"
		req := &model.UpdateTodoRequest{
			Title: &newTitle,
		}

		mockRepo.On("GetByID", uint(999)).
			Return(nil, errors.New("todo not found"))

		todo, err := service.UpdateTodo(999, req)

		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "todo not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		existingTodo := &model.Todo{
			ID:    1,
			Title: "Old Title",
		}

		newTitle := "New Title"
		req := &model.UpdateTodoRequest{
			Title: &newTitle,
		}

		mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.Todo")).
			Return(errors.New("database error"))

		todo, err := service.UpdateTodo(1, req)

		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "failed to update todo")
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_DeleteTodo(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		mockRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeleteTodo(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete error", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		mockRepo.On("Delete", uint(999)).
			Return(errors.New("todo not found"))

		err := service.DeleteTodo(999)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete todo")
		mockRepo.AssertExpectations(t)
	})
}

func TestTodoService_UpdateTodoStatus(t *testing.T) {
	t.Run("successful status update", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		existingTodo := &model.Todo{
			ID:     1,
			Title:  "Test Todo",
			Status: model.StatusPending,
		}

		mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.Todo")).Return(nil)

		todo, err := service.UpdateTodoStatus(1, model.StatusCompleted)

		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, model.StatusCompleted, todo.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("todo not found", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		mockRepo.On("GetByID", uint(999)).
			Return(nil, errors.New("todo not found"))

		todo, err := service.UpdateTodoStatus(999, model.StatusCompleted)

		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "todo not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		mockRepo := new(mocks.MockTodoRepository)
		service := NewTodoService(mockRepo)

		existingTodo := &model.Todo{
			ID:     1,
			Title:  "Test Todo",
			Status: model.StatusPending,
		}

		mockRepo.On("GetByID", uint(1)).Return(existingTodo, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.Todo")).
			Return(errors.New("database error"))

		todo, err := service.UpdateTodoStatus(1, model.StatusCompleted)

		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Contains(t, err.Error(), "failed to update todo status")
		mockRepo.AssertExpectations(t)
	})
}
