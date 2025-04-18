package service

import (
	"develapar-server/model"
	"develapar-server/repository"
	"develapar-server/utils"
)

type UserService interface {
	CreateNewUser(payload model.User) (model.User, error)
}

type userService struct {
	repo repository.UserRepository
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
