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
	CreateNewUser(ctx context.Context, payload model.User) (model.User, error)
	FindUserById(ctx context.Context, id string) (model.User, error)
	FindAllUser(ctx context.Context) ([]model.User, error)
	FindAllUserWithPagination(ctx context.Context, page, limit int) (PaginationResult, error)
	Login(ctx context.Context, payload dto.LoginDto) (dto.LoginResponseDto, error)
	RefreshToken(ctx context.Context, refreshToken string) (dto.LoginResponseDto, error)
	UpdateUser(ctx context.Context, requestingUserID int, requestingUserRole string, targetUserID int, req dto.UpdateUserRequest) (model.User, error)
	DeleteUser(ctx context.Context, requestingUserID int, requestingUserRole string, targetUserID int) error
}

type userService struct {
	repo              repository.UserRepository
	jwtService        JwtService
	passwordHasher    utils.PasswordHasher
	paginationService PaginationService
	validationService ValidationService
}

// Login implements UserService.
func (u *userService) Login(ctx context.Context, payload dto.LoginDto) (dto.LoginResponseDto, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return dto.LoginResponseDto{}, ctx.Err()
	default:
	}

	// Basic validation for login payload
	if payload.Identifier == "" {
		return dto.LoginResponseDto{}, fmt.Errorf("email is required")
	}
	if payload.Password == "" {
		return dto.LoginResponseDto{}, fmt.Errorf("password is required")
	}

	// Get user by email with context
	user, err := u.repo.GetByEmail(ctx, payload.Identifier)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return dto.LoginResponseDto{}, ctx.Err()
		}
		return dto.LoginResponseDto{}, fmt.Errorf("invalid credentials")
	}

	// Check context cancellation before password comparison
	select {
	case <-ctx.Done():
		return dto.LoginResponseDto{}, ctx.Err()
	default:
	}

	// Compare password
	if err := u.passwordHasher.ComparePasswordHash(user.Password, payload.Password); err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("invalid credentials")
	}

	// Remove password from user object for security
	user.Password = "-"
	
	// Check context cancellation before token generation
	select {
	case <-ctx.Done():
		return dto.LoginResponseDto{}, ctx.Err()
	default:
	}

	// Generate JWT token
	token, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed to create token")
	}

	// Save refresh token with context
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = u.repo.SaveRefreshToken(ctx, user.Id, token.RefreshToken, expiresAt)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return dto.LoginResponseDto{}, ctx.Err()
		}
		log.Printf("[Login] Failed saving refresh token: %v", err)
		return dto.LoginResponseDto{}, fmt.Errorf("internal server error")
	}

	return token, nil
}

// FindAllUser implements UserService.
func (u *userService) FindAllUser(ctx context.Context) ([]model.User, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Get all users from repository with context
	users, err := u.repo.GetAllUser(ctx)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}

	// Remove passwords from user data for security
	for i := range users {
		users[i].Password = "-"
	}

	return users, nil
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
func (u *userService) FindUserById(ctx context.Context, id string) (model.User, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
	}

	// Validate ID format
	newId, err := strconv.Atoi(id)
	if err != nil {
		return model.User{}, fmt.Errorf("invalid user ID format: %v", err)
	}

	if newId <= 0 {
		return model.User{}, fmt.Errorf("user ID must be greater than 0")
	}

	// Get user from repository with context
	user, err := u.repo.GetUserById(ctx, newId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.User{}, ctx.Err()
		}
		return model.User{}, fmt.Errorf("failed to fetch user: %v", err)
	}

	// Remove password from user data for security
	user.Password = "-"
	return user, nil
}

func (u *userService) CreateNewUser(ctx context.Context, payload model.User) (model.User, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
	}

	// Validate user data using validation service
	if validationErr := u.validationService.ValidateUser(ctx, payload); validationErr != nil {
		return model.User{}, validationErr
	}

	// Check context cancellation after validation
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
	}

	// Hash password before saving
	hashedPassword, err := u.passwordHasher.EncryptPassword(payload.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to encrypt password: %v", err)
	}
	payload.Password = hashedPassword

	// Save user to database with context
	createdUser, err := u.repo.CreateNewUser(ctx, payload)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.User{}, ctx.Err()
		}
		return model.User{}, fmt.Errorf("failed to create user: %v", err)
	}

	// Remove password from response for security
	createdUser.Password = "-"
	return createdUser, nil
}

func (u *userService) RefreshToken(ctx context.Context, refreshToken string) (dto.LoginResponseDto, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return dto.LoginResponseDto{}, ctx.Err()
	default:
	}

	// Decode refresh token
	decodedToken, err := url.QueryUnescape(refreshToken)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("invalid refresh token format")
	}
	
	// Check refresh token in database with context
	rt, err := u.repo.FindRefreshToken(ctx, decodedToken)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return dto.LoginResponseDto{}, ctx.Err()
		}
		return dto.LoginResponseDto{}, fmt.Errorf("invalid refresh token")
	}
	
	if rt.ExpiresAt.Before(time.Now()) {
		return dto.LoginResponseDto{}, fmt.Errorf("refresh token expired")
	}

	// Check context cancellation before getting user
	select {
	case <-ctx.Done():
		return dto.LoginResponseDto{}, ctx.Err()
	default:
	}

	// Get user with context
	user, err := u.repo.GetUserById(ctx, rt.UserID)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return dto.LoginResponseDto{}, ctx.Err()
		}
		return dto.LoginResponseDto{}, fmt.Errorf("user not found")
	}

	// Generate new token
	tokenResp, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed to generate new token")
	}

	// Update refresh token with context
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = u.repo.UpdateRefreshToken(ctx, refreshToken, tokenResp.RefreshToken, expiresAt)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return dto.LoginResponseDto{}, ctx.Err()
		}
		return dto.LoginResponseDto{}, fmt.Errorf("failed to update refresh token")
	}

	return tokenResp, nil
}

// UpdateUser implements UserService.
func (u *userService) UpdateUser(ctx context.Context, requestingUserID int, requestingUserRole string, targetUserID int, req dto.UpdateUserRequest) (model.User, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
	}

	// Secondary authorization validation using authorization helper
	if err := utils.ValidateUserPermissions(requestingUserID, requestingUserRole, targetUserID); err != nil {
		// Log security event for authorization failure
		log.Printf("[SECURITY] UpdateUser authorization failed - Requesting User: %d, Role: %s, Target User: %d, Error: %v", 
			requestingUserID, requestingUserRole, targetUserID, err)
		return model.User{}, err
	}

	// Log admin operations for audit purposes
	if utils.ValidateAdminRole(requestingUserRole) && requestingUserID != targetUserID {
		log.Printf("[AUDIT] Admin user %d (role: %s) updating user %d", 
			requestingUserID, requestingUserRole, targetUserID)
	}

	// Validate target user ID
	if targetUserID <= 0 {
		return model.User{}, fmt.Errorf("user ID must be greater than 0")
	}

	// Get existing user with context
	user, err := u.repo.GetUserById(ctx, targetUserID)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.User{}, ctx.Err()
		}
		return model.User{}, fmt.Errorf("failed to fetch user for update: %v", err)
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		// Hash new password
		hashedPassword, err := u.passwordHasher.EncryptPassword(*req.Password)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to encrypt password: %v", err)
		}
		user.Password = hashedPassword
	}

	// Check context cancellation before update
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
	}

	// Update user in repository with context
	updatedUser, err := u.repo.UpdateUser(ctx, user)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.User{}, ctx.Err()
		}
		return model.User{}, fmt.Errorf("failed to update user: %v", err)
	}

	// Log successful update operation
	log.Printf("[INFO] User %d successfully updated by user %d (role: %s)", 
		targetUserID, requestingUserID, requestingUserRole)

	// Remove password from response for security
	updatedUser.Password = "-"
	return updatedUser, nil
}

// DeleteUser implements UserService.
func (u *userService) DeleteUser(ctx context.Context, requestingUserID int, requestingUserRole string, targetUserID int) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Secondary authorization validation using authorization helper
	if err := utils.ValidateUserPermissions(requestingUserID, requestingUserRole, targetUserID); err != nil {
		// Log security event for authorization failure
		log.Printf("[SECURITY] DeleteUser authorization failed - Requesting User: %d, Role: %s, Target User: %d, Error: %v", 
			requestingUserID, requestingUserRole, targetUserID, err)
		return err
	}

	// Log admin operations for audit purposes
	if utils.ValidateAdminRole(requestingUserRole) && requestingUserID != targetUserID {
		log.Printf("[AUDIT] Admin user %d (role: %s) deleting user %d", 
			requestingUserID, requestingUserRole, targetUserID)
	}

	// Validate target user ID
	if targetUserID <= 0 {
		return fmt.Errorf("user ID must be greater than 0")
	}

	// Check context cancellation before deletion
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Delete user from repository with context
	err := u.repo.DeleteUser(ctx, targetUserID)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// Log successful delete operation
	log.Printf("[INFO] User %d successfully deleted by user %d (role: %s)", 
		targetUserID, requestingUserID, requestingUserRole)

	return nil
}

func NewUserservice(repository repository.UserRepository, jS JwtService, ph utils.PasswordHasher, paginationService PaginationService, validationService ValidationService) UserService {
	return &userService{
		repo:              repository,
		jwtService:        jS,
		passwordHasher:    ph,
		paginationService: paginationService,
		validationService: validationService,
	}
}
