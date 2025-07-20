package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"
	"develapar-server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service        service.UserService
	mD             middleware.AuthMiddleware
	rg             *gin.RouterGroup
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

// @Summary User login
// @Description Authenticate user and return access token
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body dto.LoginDto true "Login credentials"
// @Success 200 {object} middleware.SuccessResponse "Success Login"
// @Failure 400 {object} middleware.ErrorResponse "Invalid request payload"
// @Failure 401 {object} middleware.ErrorResponse "Invalid credentials"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Router /auth/login [post]
func (u *UserController) loginHandler(c *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var payload dto.LoginDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		appErr := u.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Call service with context
	response, err := u.service.Login(requestCtx, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := u.errorHandler.TimeoutError(requestCtx, "login")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := u.errorHandler.CancellationError(requestCtx, "login")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Wrap as unauthorized error
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Authentication failed")
		appErr.StatusCode = 401
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Set HttpOnly cookie for refresh token
	c.SetCookie(
		"refreshToken",
		response.RefreshToken,
		60*60*24*7, // 7 days
		"/",
		"localhost", // change to domain in production
		false,       // secure: true if HTTPS
		true,        // httpOnly
	)

	// Create success response with context
	responseData := gin.H{
		"message":      "Login successful",
		"access_token": response.AccessToken,
	}
	u.responseHelper.SendSuccess(c, responseData)
}

// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body model.User true "User registration details"
// @Success 200 {object} middleware.SuccessResponse "User successfully registered"
// @Failure 400 {object} middleware.ErrorResponse "Invalid payload"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 409 {object} middleware.ErrorResponse "User already exists"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (u *UserController) registerUser(c *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var payload model.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		appErr := u.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}


	// Call service with context
	data, err := u.service.CreateNewUser(requestCtx, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := u.errorHandler.TimeoutError(requestCtx, "user registration")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := u.errorHandler.CancellationError(requestCtx, "user registration")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Wrap as internal error
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to create user")
		appErr.StatusCode = 500
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "User successfully registered",
		"user":    data,
	}
	u.responseHelper.SendCreated(c, responseData)
}

// @Summary Get user by ID
// @Description Get user details by their ID
// @Tags Users
// @Produce json
// @Param user_id path string true "ID of the user to retrieve"
// @Success 200 {object} middleware.SuccessResponse "User details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid user ID"
// @Failure 404 {object} middleware.ErrorResponse "User not found"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /users/{user_id} [get]
func (u *UserController) findUserByIdHandler(c *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	userId := c.Param("user_id")
	if userId == "" {
		appErr := u.errorHandler.ValidationError(requestCtx, "user_id", "User ID is required")
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Call service with context
	user, err := u.service.FindUserById(requestCtx, userId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := u.errorHandler.TimeoutError(requestCtx, "get user")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := u.errorHandler.CancellationError(requestCtx, "get user")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Wrap as internal error
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve user")
		appErr.StatusCode = 500
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "User retrieved successfully",
		"user":    user,
	}
	u.responseHelper.SendSuccess(c, responseData)
}

// @Summary Get all users
// @Description Get a list of all registered users
// @Tags Users
// @Produce json
// @Success 200 {object} middleware.SuccessResponse "List of users"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /users [get]
func (u *UserController) findAllUserHandler(c *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	// Call service with context
	users, err := u.service.FindAllUser(requestCtx)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := u.errorHandler.TimeoutError(requestCtx, "get all users")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := u.errorHandler.CancellationError(requestCtx, "get all users")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Wrap as internal error
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve users")
		appErr.StatusCode = 500
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Users retrieved successfully",
		"users":   users,
	}
	u.responseHelper.SendSuccess(c, responseData)
}

// @Summary Get all users with pagination
// @Description Get a paginated list of all registered users
// @Tags Users
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} middleware.SuccessResponse "Paginated list of users"
// @Failure 400 {object} middleware.ErrorResponse "Invalid pagination parameters"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /users/paginated [get]
func (u *UserController) findAllUserWithPaginationHandler(c *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()

	// Get pagination parameters from query string
	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil || p <= 0 {
			appErr := u.errorHandler.ValidationError(requestCtx, "page", "Page must be a positive integer")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		} else {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err != nil || l <= 0 || l > 100 {
			appErr := u.errorHandler.ValidationError(requestCtx, "limit", "Limit must be a positive integer between 1 and 100")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		} else {
			limit = l
		}
	}

	// Call service with pagination and context
	result, err := u.service.FindAllUserWithPagination(requestCtx, page, limit)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := u.errorHandler.TimeoutError(requestCtx, "get paginated users")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := u.errorHandler.CancellationError(requestCtx, "get paginated users")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Wrap as internal error
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve paginated users")
		appErr.StatusCode = 500
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Create success response with context and pagination
	responseData := gin.H{
		"message": "Users retrieved successfully",
		"users":   result.Data,
	}
	u.responseHelper.SendSuccessWithServicePagination(c, responseData, result.Metadata)
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token from cookie
// @Tags Users
// @Produce json
// @Success 200 {object} middleware.SuccessResponse "Access token refreshed successfully"
// @Failure 400 {object} middleware.ErrorResponse "Refresh token not found"
// @Failure 401 {object} middleware.ErrorResponse "Invalid or expired refresh token"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Router /auth/refresh [post]
func (u *UserController) refreshTokenHandler(c *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	// Get refresh token from cookie
	cookie, err := c.Cookie("refreshToken")
	if err != nil || cookie == "" {
		appErr := u.errorHandler.ValidationError(requestCtx, "refresh_token", "Refresh token not found in cookies")
		appErr.StatusCode = 400
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Call service with context
	tokenResp, err := u.service.RefreshToken(requestCtx, cookie)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := u.errorHandler.TimeoutError(requestCtx, "refresh token")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := u.errorHandler.CancellationError(requestCtx, "refresh token")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Wrap as unauthorized error
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Invalid or expired refresh token")
		appErr.StatusCode = 401
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Set new refresh token in cookie
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	c.SetCookie("refreshToken", tokenResp.RefreshToken, int(refreshExpiry.Sub(time.Now()).Seconds()), "/", "", true, true)

	// Create success response with context
	responseData := gin.H{
		"message":      "Token refreshed successfully",
		"access_token": tokenResp.AccessToken,
	}
	u.responseHelper.SendSuccess(c, responseData)
}

// @Summary Update user profile
// @Description Update user profile information
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path string true "ID of the user to update"
// @Param payload body dto.UpdateUserRequest true "User update details"
// @Success 200 {object} middleware.SuccessResponse "User updated successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid user ID or payload"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 403 {object} middleware.ErrorResponse "Forbidden"
// @Failure 404 {object} middleware.ErrorResponse "User not found"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{user_id} [put]
func (u *UserController) updateUserHandler(c *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	userIdStr := c.Param("user_id")
	if userIdStr == "" {
		appErr := u.errorHandler.ValidationError(requestCtx, "user_id", "User ID is required")
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		appErr := u.errorHandler.ValidationError(requestCtx, "user_id", "Invalid user ID format")
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Extract user ID and role from Gin context (set by auth middleware)
	requestingUserID, err := utils.GetUserIDFromGinContext(c)
	if err != nil {
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Failed to extract user ID from context")
		appErr.StatusCode = 401
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	requestingUserRole, err := utils.GetUserRoleFromContext(c)
	if err != nil {
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Failed to extract user role from context")
		appErr.StatusCode = 401
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Validate authorization using helper function
	if err := utils.ValidateUserPermissions(requestingUserID, requestingUserRole, userId); err != nil {
		// Check if it's an authorization error
		if utils.IsAuthorizationError(err) {
			authErr := err.(*utils.AuthorizationError)
			if authErr.Code == utils.ErrForbidden {
				appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrForbidden, authErr.Message)
				appErr.StatusCode = 403
				u.errorHandler.HandleError(requestCtx, c, appErr)
				return
			}
			// Handle other authorization errors as bad request
			appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrBadRequest, authErr.Message)
			appErr.StatusCode = 400
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		// Handle unexpected errors
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Authorization validation failed")
		appErr.StatusCode = 500
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	var payload dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		appErr := u.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Call service with context
	updatedUser, err := u.service.UpdateUser(requestCtx, userId, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := u.errorHandler.TimeoutError(requestCtx, "update user")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := u.errorHandler.CancellationError(requestCtx, "update user")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Wrap as internal error
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to update user")
		appErr.StatusCode = 500
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "User updated successfully",
		"user":    updatedUser,
	}
	u.responseHelper.SendSuccess(c, responseData)
}

// @Summary Delete user account
// @Description Delete user account by ID
// @Tags Users
// @Produce json
// @Param user_id path string true "ID of the user to delete"
// @Success 200 {object} middleware.SuccessResponse "User deleted successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid user ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 403 {object} middleware.ErrorResponse "Forbidden"
// @Failure 404 {object} middleware.ErrorResponse "User not found"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{user_id} [delete]
func (u *UserController) deleteUserHandler(c *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	userIdStr := c.Param("user_id")
	if userIdStr == "" {
		appErr := u.errorHandler.ValidationError(requestCtx, "user_id", "User ID is required")
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		appErr := u.errorHandler.ValidationError(requestCtx, "user_id", "Invalid user ID format")
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Extract user ID and role from Gin context (set by auth middleware)
	requestingUserID, err := utils.GetUserIDFromGinContext(c)
	if err != nil {
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Failed to extract user ID from context")
		appErr.StatusCode = 401
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	requestingUserRole, err := utils.GetUserRoleFromContext(c)
	if err != nil {
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Failed to extract user role from context")
		appErr.StatusCode = 401
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Validate authorization using helper function
	if err := utils.ValidateUserPermissions(requestingUserID, requestingUserRole, userId); err != nil {
		// Check if it's an authorization error
		if utils.IsAuthorizationError(err) {
			authErr := err.(*utils.AuthorizationError)
			if authErr.Code == utils.ErrForbidden {
				appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrForbidden, authErr.Message)
				appErr.StatusCode = 403
				u.errorHandler.HandleError(requestCtx, c, appErr)
				return
			}
			// Handle other authorization errors as bad request
			appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrBadRequest, authErr.Message)
			appErr.StatusCode = 400
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		// Handle unexpected errors
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Authorization validation failed")
		appErr.StatusCode = 500
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Call service with context
	err = u.service.DeleteUser(requestCtx, userId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := u.errorHandler.TimeoutError(requestCtx, "delete user")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := u.errorHandler.CancellationError(requestCtx, "delete user")
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			u.errorHandler.HandleError(requestCtx, c, appErr)
			return
		}

		// Wrap as internal error
		appErr := u.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to delete user")
		appErr.StatusCode = 500
		u.errorHandler.HandleError(requestCtx, c, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "User deleted successfully",
	}
	u.responseHelper.SendSuccess(c, responseData)
}

func (u *UserController) Route() {
	router := u.rg.Group("/users")
	{
		router.GET("/", u.findAllUserHandler)
		router.GET("/paginated", u.findAllUserWithPaginationHandler)
		router.GET("/:user_id", u.findUserByIdHandler)
		router.PUT("/:user_id", u.mD.CheckToken(), u.updateUserHandler)
		router.DELETE("/:user_id", u.mD.CheckToken(), u.deleteUserHandler)
	}

	r := u.rg.Group("/auth")
	{
		r.POST("/login", u.loginHandler)
		r.POST("/register", u.registerUser)
		r.POST("/refresh", u.refreshTokenHandler)
	}
}

func NewUserController(uS service.UserService, mD middleware.AuthMiddleware, rg *gin.RouterGroup, errorHandler middleware.ErrorHandler) *UserController {
	return &UserController{
		service:        uS,
		mD:             mD,
		rg:             rg,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}
