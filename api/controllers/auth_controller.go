package controllers

import "proxy-fileserver/services"

type AuthController struct {
	AuthService *services.AuthService
}
