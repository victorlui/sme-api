package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/victorlui/sma-api/internal/db"
	"github.com/victorlui/sma-api/internal/handler"
	"github.com/victorlui/sma-api/internal/logger"
	"github.com/victorlui/sma-api/internal/repository"
)

func SetupRouter(db *db.PostgresDB) *gin.Engine {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "*"}, // Ou especifique os domínios permitidos
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(logger.CustomLogger())

	userRepo := repository.NewUserRepository(db.Pool)
	userHandler := handler.NewUserHandler(userRepo)

	schoolRepo := repository.NewSchoolRepository(db.Pool)
	schoolHandler := handler.NewSchoolHandler(schoolRepo)

	studentRepo := repository.NewStudentRepository(db.Pool, schoolRepo)
	studentHandler := handler.NewStudentHandler(studentRepo)

	studentServiceRepo := repository.NewServiceStudentRepository(db.Pool, studentRepo, userRepo)
	studentServiceHandler := handler.NewServiceStudentHandler(studentServiceRepo)

	SetupRoutes(r, userHandler, schoolHandler, studentHandler, studentServiceHandler)

	return r
}
