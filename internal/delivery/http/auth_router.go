package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorlui/sma-api/internal/handler"
	"github.com/victorlui/sma-api/internal/model"
	"github.com/victorlui/sma-api/internal/utils"
)

func RegisterUser(userHandler *handler.UserHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userRequest model.CreateUserRequest

		if err := ctx.ShouldBindJSON(&userRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": utils.CustomErrorHandling(err),
			})
			return
		}

		hashedPassword, err := utils.HashPassword(userRequest.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao codificar a senha.",
			})
			return
		}

		// Substitui a senha original pela senha codificada
		userRequest.Password = hashedPassword

		user, err := userHandler.CreateNewUser(userRequest)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, user)
	}
}

func Login(userHandler *handler.UserHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.LoginRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.CustomErrorHandling(err)})
			return
		}

		user, err := userHandler.Login(ctx.Request.Context(), req.Email, req.Password)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}
