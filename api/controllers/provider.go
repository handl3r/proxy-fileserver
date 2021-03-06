package controllers

import (
	"context"
	"proxy-fileserver/services"
)

type ControllerProvider interface {
	GetStreamFileController() *StreamFileController
}

type controllerProviderImpl struct {
	streamFileController *StreamFileController
}

func NewControllerProvider(ctx context.Context, fileSystemService *services.FileSystemService) ControllerProvider {
	return &controllerProviderImpl{
		streamFileController: NewStreamFileController(fileSystemService),
	}
}

func (c controllerProviderImpl) GetStreamFileController() *StreamFileController {
	return c.streamFileController
}
