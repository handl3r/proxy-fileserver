package dtos

type Token struct {
	Token string `json:"token"`
	Path  string `json:"path"`
}

type CreateTokenRequest struct {
	Path string `json:"path"`
}
