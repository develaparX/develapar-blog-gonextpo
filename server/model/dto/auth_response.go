package dto

type LoginDto struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password"`
}

type LoginResponseDto struct {
	Token string `json:"token"`
}
