package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/VoLKoDaV13579/GoToDoList/internal/model"
	"github.com/VoLKoDaV13579/GoToDoList/internal/service"
	"github.com/VoLKoDaV13579/GoToDoList/pkg/validator"
	"github.com/gin-gonic/gin"
)

// TodoHandler обрабатывает HTTP-запросы для работы с задачами
type TodoHandler struct {
	service   service.TodoService
	validator *validator.CustomValidator
}

// NewTodoHandler создает новый экземпляр обработчика задач
func NewTodoHandler(service service.TodoService, validator *validator.CustomValidator) *TodoHandler {
	return &TodoHandler{
		service:   service,
		validator: validator,
	}
}

// CreateTodo godoc
// @Summary Создать новую задачу
// @Description Создает новую задачу в системе
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body model.CreateTodoRequest true "Данные задачи"
// @Success 201 {object} model.Todo
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req model.CreateTodoRequest

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Неверный формат запроса",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Валидируем данные
	if err := h.validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Создаем задачу через сервис
	todo, err := h.service.CreateTodo(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Не удалось создать задачу",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodoByID godoc
// @Summary Получить задачу по ID
// @Description Возвращает задачу по её идентификатору
// @Tags todos
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} model.Todo
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	// Получаем ID из параметров URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Неверный формат ID",
		})
		return
	}

	// Получаем задачу через сервис
	todo, err := h.service.GetTodoByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "not_found",
				Message: "Задача не найдена",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Не удалось получить задачу",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// GetAllTodos godoc
// @Summary Получить список задач
// @Description Возвращает список задач с фильтрацией и пагинацией
// @Tags todos
// @Produce json
// @Param status query string false "Фильтр по статусу" Enums(pending, in_progress, completed)
// @Param priority query int false "Фильтр по приоритету"
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(10)
// @Success 200 {object} model.TodoListResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/todos [get]
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	var filter model.TodoFilter

	// Парсим query параметры
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_query",
			Message: "Неверные параметры запроса",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Получаем список задач через сервис
	response, err := h.service.GetAllTodos(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Не удалось получить список задач",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTodo godoc
// @Summary Обновить задачу
// @Description Частично обновляет данные задачи
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param todo body model.UpdateTodoRequest true "Данные для обновления"
// @Success 200 {object} model.Todo
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/todos/{id} [patch]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	// Получаем ID из параметров URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Неверный формат ID",
		})
		return
	}

	var req model.UpdateTodoRequest

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Неверный формат запроса",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Валидируем данные
	if err := h.validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Обновляем задачу через сервис
	todo, err := h.service.UpdateTodo(uint(id), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "not_found",
				Message: "Задача не найдена",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Не удалось обновить задачу",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по её идентификатору
// @Tags todos
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	// Получаем ID из параметров URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Неверный формат ID",
		})
		return
	}

	// Удаляем задачу через сервис
	if err := h.service.DeleteTodo(uint(id)); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "not_found",
				Message: "Задача не найдена",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Не удалось удалить задачу",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Message: "Задача успешно удалена",
	})
}

// UpdateTodoStatus godoc
// @Summary Обновить статус задачи
// @Description Обновляет только статус задачи
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param status body map[string]string true "Новый статус" example({"status": "completed"})
// @Success 200 {object} model.Todo
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/todos/{id}/status [patch]
func (h *TodoHandler) UpdateTodoStatus(c *gin.Context) {
	// Получаем ID из параметров URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Неверный формат ID",
		})
		return
	}

	var req struct {
		Status model.TodoStatus `json:"status" binding:"required"`
	}

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Неверный формат запроса",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Проверяем валидность статуса
	if req.Status != model.StatusPending && req.Status != model.StatusInProgress && req.Status != model.StatusCompleted {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_status",
			Message: "Неверный статус. Допустимые значения: pending, in_progress, completed",
		})
		return
	}

	// Обновляем статус через сервис
	todo, err := h.service.UpdateTodoStatus(uint(id), req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "not_found",
				Message: "Задача не найдена",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_error",
			Message: "Не удалось обновить статус задачи",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// HealthCheck godoc
// @Summary Проверка здоровья API
// @Description Возвращает статус работоспособности API
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *TodoHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "API работает нормально",
	})
}
