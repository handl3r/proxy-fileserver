package controllers

import (
	"context"
	"proxy-fileserver/services"
)

type ControllerProvider interface {
	GetStreamFileController() *StreamFileController
	GetAuthController() *AuthController
}

type controllerProviderImpl struct {
	streamFileController *StreamFileController
	authController       *AuthController
}

func NewControllerProvider(ctx context.Context, fileSystemService *services.FileSystemService, authService *services.AuthService) ControllerProvider {
	return &controllerProviderImpl{
		streamFileController: NewStreamFileController(fileSystemService),
		authController:       NewAuthController(authService),
	}
}

func (c controllerProviderImpl) GetStreamFileController() *StreamFileController {
	return c.streamFileController
}

func (c controllerProviderImpl) GetAuthController() *AuthController {
	return c.authController
}
