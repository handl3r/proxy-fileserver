package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy-fileserver/common/log"
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
	srcStream, err := c.fileSystemService.GetSourceStream(path)
	if err != nil {
		if err == enums.ErrorNoContent {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.AbortWithStatusJSON(err.GetCode(), err)
		return
	}
	ctx.DataFromReader(200, -1, "application/octet-stream", srcStream, nil)
	log.Infof("Finish streaming file %s to client %s", path, ctx.ClientIP())
}

func (c *StreamFileController) GetFileBasicHttp(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Path
	path := rawPath[1:len(rawPath)]
	err := c.fileSystemService.StreamFile(w, path)
	if err != nil {
		w.WriteHeader(err.GetCode())
	}
}
