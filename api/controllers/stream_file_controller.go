package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-fileserver/enums"
	"proxy-fileserver/services"
)

type StreamFileController struct {
	fileSystemService *services.FileSystemService
}

func NewStreamFileController(service *services.FileSystemService) *StreamFileController {
	return &StreamFileController{
		fileSystemService: service,
	}
}

func (c *StreamFileController) GetFile(ctx *gin.Context) {
	rawPath := ctx.Request.URL.Path
	path := rawPath[1:len(rawPath)]
	err := c.fileSystemService.StreamFile(ctx.Writer, path)
	if err != nil {
		if err == enums.ErrorNoContent {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.AbortWithStatusJSON(err.GetCode(), err)
		return
	}
	ctx.Next()
}
