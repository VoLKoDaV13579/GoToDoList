package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS - middleware для настройки CORS (Cross-Origin Resource Sharing)
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Разрешаем запросы со всех источников (для production лучше указать конкретные домены)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Разрешаем определенные методы
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Разрешаем определенные заголовки
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		// Разрешаем отправку credentials (cookies, auth headers)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Максимальное время кэширования preflight-запроса (в секундах)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Обрабатываем preflight-запросы
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
