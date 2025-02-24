package http

import (
	"github.com/gin-gonic/gin"
	"github.com/victorlui/sma-api/internal/handler"
	"github.com/victorlui/sma-api/internal/middleware"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, schoolHandler *handler.SchoolHandler, studentHandler *handler.StudentHandler, studentServiceHandler *handler.ServiceStudentHandler) {
	// Rotas de autenticação
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", RegisterUser(userHandler))
		authRoutes.POST("/login", Login(userHandler))
	}

	// Rotas de usuários
	userRoutes := r.Group("/users", middleware.AuthMiddleware())
	{
		userRoutes.GET("", UsersRoute(userHandler))                // GET /users
		userRoutes.GET("/:user_id", GetUserByIDRoute(userHandler)) // GET /users/:user_id
	}

	// Rotas de escolas
	schoolRoutes := r.Group("/schools", middleware.AuthMiddleware())
	{
		schoolRoutes.POST("", CreateNewSchool(schoolHandler)) // POST /schools
		schoolRoutes.GET("", GETSchools(schoolHandler))
		schoolRoutes.DELETE("/:id", DeleteSchoolRoute((schoolHandler)))
		schoolRoutes.PATCH("/:id", UpdateSchoolRoute((schoolHandler))) // GET /schools
	}

	// Rotas de estudantes
	studentRoutes := r.Group("/students", middleware.AuthMiddleware())
	{
		studentRoutes.POST("", CreateStudentRoute(studentHandler))              // POST /students
		studentRoutes.GET("", GetStudentsRoute(studentHandler))                 // GET /students
		studentRoutes.GET("/:student_id", GetStudentsByIdRoute(studentHandler)) // GET /students/:student_id
		studentRoutes.DELETE("/:id", DeleteStudentRoute(studentHandler))
	}

	// Rota de upload
	r.POST("/upload", middleware.AuthMiddleware(), UploadRoute())

	// Rota de serviço de estudantes
	studentServiceRoute := r.Group("service-students", middleware.AuthMiddleware())
	{
		studentServiceRoute.POST("", middleware.AuthMiddleware(), ServiceStundentRoute(studentServiceHandler))
		studentServiceRoute.GET("", middleware.AuthMiddleware(), GETServiceStudent(studentServiceHandler))
		studentServiceRoute.GET("/:student_service_id", middleware.AuthMiddleware(), ServiceStundentRoute(studentServiceHandler))
		studentServiceRoute.DELETE("/:id", middleware.AuthMiddleware(), DeleteServiceStudentRoute(studentServiceHandler))
		studentServiceRoute.PATCH("/:id", middleware.AuthMiddleware(), UpdateServiceStudentRoute(studentServiceHandler))
	}
}
