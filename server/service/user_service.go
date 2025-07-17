package service

import (
	"context"
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
	FindAllUserWithPagination(ctx context.Context, page, limit int) (PaginationResult, error)
	Login(payload dto.LoginDto) (dto.LoginResponseDto, error)
	RefreshToken(refreshToken string) (dto.LoginResponseDto, error)
}

type userService struct {
	repo              repository.UserRepository
	jwtService        JwtService
	passwordHasher    utils.PasswordHasher
	paginationService PaginationService
}

// Login implements UserService.
func (u *userService) Login(payload dto.LoginDto) (dto.LoginResponseDto, error) {
	user, err := u.repo.GetByEmail(context.Background(), payload.Identifier)
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
	err = u.repo.SaveRefreshToken(context.Background(), user.Id, token.RefreshToken, expiresAt)
	if err != nil {
		log.Printf("[Login] Failed saving refresh token: %v", err)
		return dto.LoginResponseDto{}, fmt.Errorf("internal server error")
	}

	return token, nil

}

// FindAllUser implements UserService.
func (u *userService) FindAllUser() ([]model.User, error) {
	return u.repo.GetAllUser(context.Background())
}

// FindAllUserWithPagination implements UserService with pagination support
func (u *userService) FindAllUserWithPagination(ctx context.Context, page, limit int) (PaginationResult, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := u.paginationService.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, fmt.Errorf("pagination validation failed: %v", err)
	}

	// Get paginated users from repository
	users, total, repoErr := u.repo.GetAllUserWithPagination(ctx, query.Offset, query.Limit)
	if repoErr != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, fmt.Errorf("failed to fetch users: %v", repoErr)
	}

	// Remove passwords from user data for security
	for i := range users {
		users[i].Password = "-"
	}

	// Create pagination result
	result, paginationErr := u.paginationService.Paginate(ctx, users, total, query)
	if paginationErr != nil {
		return PaginationResult{}, fmt.Errorf("failed to create pagination result: %v", paginationErr)
	}

	return result, nil
}

// FindUserById implements UserService.
func (u *userService) FindUserById(id string) (model.User, error) {
	newId, err := strconv.Atoi(id)

	if err != nil {
		return model.User{}, err
	}

	return u.repo.GetUserById(context.Background(), newId)
}

func (u *userService) CreateNewUser(payload model.User) (model.User, error) {
	// Hash password dulu sebelum disimpan
	hashedPassword, err := u.passwordHasher.EncryptPassword(payload.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to encrypt password: %v", err)
	}
	payload.Password = hashedPassword

	// Simpan user ke database
	createdUser, err := u.repo.CreateNewUser(context.Background(), payload)
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
	rt, err := u.repo.FindRefreshToken(context.Background(), decodedToken)
	if err != nil || rt.ExpiresAt.Before(time.Now()) {
		return dto.LoginResponseDto{}, err
	}

	// Ambil user
	user, err := u.repo.GetUserById(context.Background(), rt.UserID)
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
	err = u.repo.UpdateRefreshToken(context.Background(), refreshToken, tokenResp.RefreshToken, expiresAt)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed to update refresh token")
	}

	return tokenResp, nil
}

func NewUserservice(repository repository.UserRepository, jS JwtService, ph utils.PasswordHasher, paginationService PaginationService) UserService {
	return &userService{
		repo:              repository,
		jwtService:        jS,
		passwordHasher:    ph,
		paginationService: paginationService,
	}
}
