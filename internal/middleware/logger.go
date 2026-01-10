package middleware

import (
	"time"

	"github.com/VoLKoDaV13579/GoToDoList/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Logger - middleware для логирования HTTP-запросов
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Засекаем время начала обработки запроса
		startTime := time.Now()

		// Получаем request ID из контекста (если установлен)
		requestID := c.GetString("request_id")

		// Логируем входящий запрос
		log.Info("Incoming request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		// Обрабатываем запрос
		c.Next()

		// Вычисляем время обработки
		duration := time.Since(startTime)

		// Логируем результат обработки запроса
		log.Info("Request completed",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"ip", c.ClientIP(),
		)

		// Если были ошибки, логируем их отдельно
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Error("Request error",
					"request_id", requestID,
					"error", err.Error(),
				)
			}
		}
	}
}
