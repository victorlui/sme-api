package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorlui/sma-api/internal/handler"
)

func UsersRoute(userHandler *handler.UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userHandler.GetAllUsers()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})

			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUserByIDRoute(userHandler *handler.UserHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("user_id")

		//converter de string para int
		id, _ := strconv.Atoi(userID)

		user, err := userHandler.GetUserByID(id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, user)
	}
}
