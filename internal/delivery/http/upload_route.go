package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/victorlui/sma-api/internal/frameworks/backblaze"
)

func UploadRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		file, err := ctx.FormFile("file")

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Falha ao obter o arquivo"})
			return
		}

		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao abrir o arquivo"})
			return
		}
		defer src.Close()

		name := filepath.Base(file.Filename)

		fileName, err := backblaze.UploadFileService(name, src)

		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
			return
		}

		fmt.Println("File enviado", fileName)

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Arquivo enviado com sucesso",
			"file":    fileName,
		})
	}
}
