package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorlui/sma-api/internal/handler"
	"github.com/victorlui/sma-api/internal/model"
	"github.com/victorlui/sma-api/internal/utils"
)

func ServiceStundentRoute(serviceStudenthandler *handler.ServiceStudentHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var studentModel model.ServiceStudent

		// validate := validator.New()
		// _ = validate.RegisterValidation("ISO8601date", utils.CustomDate)

		// err := validate.Struct(studentModel)

		// if err != nil {
		// 	fmt.Printf("Err(s):\n%+v\n", err)
		// 	ctx.JSON(http.StatusBadRequest, gin.H{
		// 		"error": utils.CustomErrorHandling(err),
		// 	})

		// 	return

		// }

		if err := ctx.ShouldBindJSON(&studentModel); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": utils.CustomErrorHandling(err),
			})

			return
		}

		student_service, err := serviceStudenthandler.CreateNewServiceStudentHandler(studentModel)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao realizar atendimento",
			})
			return
		}

		ctx.JSON(http.StatusCreated, student_service)
	}
}

func UpdateServiceStudentRoute(serviceStudentHandler *handler.ServiceStudentHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "Erro ao definir id"})
			return
		}

		var studentModel model.ServiceStudent

		if err := ctx.ShouldBindJSON(&studentModel); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": utils.CustomErrorHandling(err),
			})
			return
		}

		student_service, err := serviceStudentHandler.UpdateNewServiceStudentHandler(id, studentModel)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, student_service)
	}
}

func GETServiceStudent(serviceStudenthandler *handler.ServiceStudentHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}
		limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}

		date_initial := ctx.DefaultQuery("date_initial", "")
		date_end := ctx.DefaultQuery("date_end", "")
		idStudent := ctx.DefaultQuery("idStudent", "")

		service_students, err := serviceStudenthandler.GetAllServiceStudents(page, limit, date_initial, date_end, idStudent)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		ctx.JSON(http.StatusOK, service_students)
	}
}

func DeleteServiceStudentRoute(serviceStudentHandler *handler.ServiceStudentHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		err = serviceStudentHandler.DeleteServiceStudentHandler(id)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Atendimento excluido com sucesso"})

	}
}
