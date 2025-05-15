package service

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"develapar-server/utils"
	"fmt"
	"strconv"
)

type UserService interface {
	CreateNewUser(payload model.User) (model.User, error)
	FindUserById(id string) (model.User, error)
	FindAllUser() ([]model.User, error)
	Login(payload dto.LoginDto) (dto.LoginResponseDto, error)
}

type userService struct {
	repo       repository.UserRepository
	jwtService JwtService
}

// Login implements UserService.
func (u *userService) Login(payload dto.LoginDto) (dto.LoginResponseDto, error) {
	user, err := u.repo.GetByEmail(payload.Identifier)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("invalid email credentials")
	}

	if err := utils.ComparePasswordHash(user.Password, payload.Password); err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("invalid password credentials")
	}

	user.Password = "-"
	token, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed to create token")
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
	hashedPassword, err := utils.EncryptPassword(payload.Password)
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

func NewUserservice(repository repository.UserRepository, jS JwtService) UserService {
	return &userService{repo: repository, jwtService: jS}
}
