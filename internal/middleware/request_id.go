package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID - middleware для генерации уникального ID для каждого запроса
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, есть ли уже request ID в заголовке
		requestID := c.GetHeader("X-Request-ID")

		// Если нет, генерируем новый UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Сохраняем request ID в контексте для использования в других middleware и handlers
		c.Set("request_id", requestID)

		// Добавляем request ID в заголовок ответа
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
