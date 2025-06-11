package dto

type LoginDto struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password"`
}

type LoginResponseDto struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
