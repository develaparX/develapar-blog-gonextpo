package service

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"develapar-server/utils"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
)

type UserService interface {
	CreateNewUser(payload model.User) (model.User, error)
	FindUserById(id string) (model.User, error)
	FindAllUser() ([]model.User, error)
	Login(payload dto.LoginDto) (dto.LoginResponseDto, error)
	RefreshToken(refreshToken string) (dto.LoginResponseDto, error)
}

type userService struct {
	repo           repository.UserRepository
	jwtService     JwtService
	passwordHasher utils.PasswordHasher
}

// Login implements UserService.
func (u *userService) Login(payload dto.LoginDto) (dto.LoginResponseDto, error) {
	user, err := u.repo.GetByEmail(payload.Identifier)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("invalid credentials")
	}

	if err := u.passwordHasher.ComparePasswordHash(user.Password, payload.Password); err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("invalid credentials")
	}

	user.Password = "-"
	token, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed to create token")
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour) // misalnya refresh token 7 hari
	err = u.repo.SaveRefreshToken(user.Id, token.RefreshToken, expiresAt)
	if err != nil {
		log.Printf("[Login] Failed saving refresh token: %v", err)
		return dto.LoginResponseDto{}, fmt.Errorf("internal server error")
	}

	return token, nil

}

// FindAllUser implements UserService.
func (u *userService) FindAllUser() ([]model.User, error) {
	return u.repo.GetAllUser()
}

// FindUserById implements UserService.
func (u *userService) FindUserById(id string) (model.User, error) {
	newId, err := strconv.Atoi(id)

	if err != nil {
		return model.User{}, err
	}

	return u.repo.GetUserById(newId)
}

func (u *userService) CreateNewUser(payload model.User) (model.User, error) {
	// Hash password dulu sebelum disimpan
	hashedPassword, err := u.passwordHasher.EncryptPassword(payload.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to encrypt password: %v", err)
	}
	payload.Password = hashedPassword

	// Simpan user ke database
	createdUser, err := u.repo.CreateNewUser(payload)
	if err != nil {
		return model.User{}, err
	}

	// Jangan balikin password-nya
	createdUser.Password = "-"
	return createdUser, nil
}

func (u *userService) RefreshToken(refreshToken string) (dto.LoginResponseDto, error) {
	// Cek refresh token di database
	decodedToken, err := url.QueryUnescape(refreshToken)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("invalid refresh token format")
	}
	
	// Cek refresh token di database
	rt, err := u.repo.FindRefreshToken(decodedToken)
	if err != nil || rt.ExpiresAt.Before(time.Now()) {
		return dto.LoginResponseDto{}, err
	}

	// Ambil user
	user, err := u.repo.GetUserById(rt.UserID)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("user not found")
	}

	// Generate token baru
	tokenResp, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed to generate new token")
	}

	// Update refresh token lama → bisa regenerasi token atau pakai yang sama (opsional)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = u.repo.UpdateRefreshToken(refreshToken, tokenResp.RefreshToken, expiresAt)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed to update refresh token")
	}

	return tokenResp, nil
}

func NewUserservice(repository repository.UserRepository, jS JwtService, ph utils.PasswordHasher) UserService {
	return &userService{repo: repository, jwtService: jS, passwordHasher: ph}
}
