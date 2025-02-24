package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorlui/sma-api/internal/handler"
	"github.com/victorlui/sma-api/internal/model"
	"github.com/victorlui/sma-api/internal/utils"
)

func CreateNewSchool(schoolHandler *handler.SchoolHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var schoolRequest model.School

		if err := ctx.ShouldBindJSON(&schoolRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": utils.CustomErrorHandling(err),
			})
			return
		}

		school, err := schoolHandler.CreateNewSchool(schoolRequest)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, school)
	}
}

func UpdateSchoolRoute(schoolHandler *handler.SchoolHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Erro ao pegar id"})
			return
		}

		var schoolRequest model.School

		if err := ctx.ShouldBindJSON(&schoolRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": utils.CustomErrorHandling(err)})
			return
		}

		school, err := schoolHandler.UpdateSchoolHandler(id, schoolRequest)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, school)
	}
}

func GETSchools(schoolHandler *handler.SchoolHandler) gin.HandlerFunc {
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
		schools, totalRecords, err := schoolHandler.GetSchoolsAll(page, limit, name)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		totalPages := (totalRecords + limit - 1) / limit

		ctx.JSON(http.StatusCreated, gin.H{
			"itens":      schools,
			"totalItems": totalRecords,
			"totalPage":  totalPages,
		})
	}
}

func DeleteSchoolRoute(schoolHandler *handler.SchoolHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro pegar o id"})
			return
		}

		err = schoolHandler.DeleteSchool(id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"menssage": "Escola deletado com sucesso"})

	}
}
