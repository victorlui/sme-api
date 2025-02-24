package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CustomLogger é um middleware para personalizar o logger no Gin
func CustomLogger() gin.HandlerFunc {
	// Inicializa um logger usando zap
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	return func(c *gin.Context) {
		startTime := time.Now()

		// Processa a requisição
		c.Next()

		// Calcula o tempo de execução
		duration := time.Since(startTime)

		// Loga as informações da requisição
		zapLogger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}
