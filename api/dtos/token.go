package dtos

// ValidateTokenRequest represent request body to validate a token
// input of controllers.AuthController.ValidateToken
type ValidateTokenRequest struct {
	Token string `json:"token" validate:"required"`
	Path  string `json:"path"`
}

// CreateTokenRequest represent request body to create a token
// input of controllers.AuthController.GetToken
type CreateTokenRequest struct {
	Type string `json:"type" validate:"oneof=medium_level_token high_level_token"`
	Path string `json:"path" validate:"required_if=Type high_level_token"`
}

// CreateTokenResponse represent response object of request create token
type CreateTokenResponse struct {
	Token string `json:"token"`
}
