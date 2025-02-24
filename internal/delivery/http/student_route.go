package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorlui/sma-api/internal/handler"
	"github.com/victorlui/sma-api/internal/model"
	"github.com/victorlui/sma-api/internal/utils"
)

func CreateStudentRoute(studentHandler *handler.StudentHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var studentModel model.Student

		if err := ctx.ShouldBindJSON(&studentModel); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": utils.CustomErrorHandling(err),
			})

			return
		}

		student, err := studentHandler.CreateNewStudent(studentModel)

		if err != nil {

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, student)
	}
}

func GetStudentsRoute(studentHandler *handler.StudentHandler) gin.HandlerFunc {
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

		name := ctx.DefaultQuery("name", "")
		students, totalRecords, err := studentHandler.GetAllStudents(page, limit, name)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		totalPages := (totalRecords + limit - 1) / limit

		ctx.JSON(http.StatusCreated, gin.H{
			"itens":      students,
			"totalItems": totalRecords,
			"totalPage":  totalPages,
		})
	}
}

func GetStudentsByIdRoute(studentHandler *handler.StudentHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		studentID := ctx.Param("student_id")

		//converter de string para int
		id, _ := strconv.Atoi(studentID)

		student, err := studentHandler.GetStudentByIDHandler(id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, student)
	}
}

func DeleteStudentRoute(studentHandler *handler.StudentHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao pegar id"})
			return
		}

		err = studentHandler.DeleteHandler(id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"error": "Aluno deletado com sucesso"})
	}
}
