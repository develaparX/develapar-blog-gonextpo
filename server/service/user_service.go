package service

import (
	"develapar-server/model"
	"develapar-server/repository"
	"develapar-server/utils"
	"strconv"
)

type UserService interface {
	CreateNewUser(payload model.User) (model.User, error)
	FindUserById(id string) (model.User, error)
	FindAllUser() ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
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
	passwordHash, err := utils.EncryptPassword(payload.Password)
	if err != nil {
		return model.User{}, err
	}
	payload.Password = passwordHash

	return u.repo.CreateNewUser(payload)
}

func NewUserservice(repository repository.UserRepository) UserService {
	return &userService{repo: repository}
}
