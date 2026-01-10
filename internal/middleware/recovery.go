package middleware

import (
	"net/http"

	"github.com/VoLKoDaV13579/GoToDoList/internal/model"
	"github.com/VoLKoDaV13579/GoToDoList/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Recovery - middleware для обработки паник и восстановления приложения
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("request_id")

				// Логируем критическую ошибку
				log.Error("Panic recovered",
					"request_id", requestID,
					"error", err,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"ip", c.ClientIP(),
				)

				// Возвращаем ошибку клиенту
				c.JSON(http.StatusInternalServerError, model.ErrorResponse{
					Error:   "internal_server_error",
					Message: "Произошла внутренняя ошибка сервера",
				})

				// Прерываем обработку запроса
				c.Abort()
			}
		}()

		c.Next()
	}
}
