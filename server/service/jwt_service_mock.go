package service

import (
	"develapar-server/model"
	"develapar-server/model/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func (j *JwtServiceMock) GenerateToken(payload model.User) (dto.LoginResponseDto, error) {
	args := j.Called(payload)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

func (j *JwtServiceMock) VerifyToken(token string) (jwt.MapClaims, error) {
	args := j.Called(token)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}

func (j *JwtServiceMock) GenerateRefreshToken() (string, error) {
	args := j.Called()
	return args.Get(0).(string), args.Error(1)
}
