package handler

import (
	"github.com/VoLKoDaV13579/GoToDoList/internal/middleware"
	"github.com/VoLKoDaV13579/GoToDoList/pkg/logger"
	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает все маршруты приложения
func SetupRouter(todoHandler *TodoHandler, log *logger.Logger) *gin.Engine {
	// Создаем роутер
	router := gin.New()

	// Подключаем middleware
	router.Use(middleware.Logger(log))           // Логирование запросов
	router.Use(middleware.Recovery(log))         // Обработка паник
	router.Use(middleware.CORS())                // CORS для фронтенда
	router.Use(middleware.RequestID())           // Уникальный ID для каждого запроса

	// Health check endpoint (без префикса версии API)
	router.GET("/health", todoHandler.HealthCheck)

	// API v1 группа маршрутов
	v1 := router.Group("/api/v1")
	{
		// Группа маршрутов для работы с задачами
		todos := v1.Group("/todos")
		{
			todos.POST("", todoHandler.CreateTodo)                // Создать задачу
			todos.GET("", todoHandler.GetAllTodos)                // Получить все задачи
			todos.GET("/:id", todoHandler.GetTodoByID)            // Получить задачу по ID
			todos.PATCH("/:id", todoHandler.UpdateTodo)           // Обновить задачу
			todos.DELETE("/:id", todoHandler.DeleteTodo)          // Удалить задачу
			todos.PATCH("/:id/status", todoHandler.UpdateTodoStatus) // Обновить статус задачи
		}
	}

	return router
}
