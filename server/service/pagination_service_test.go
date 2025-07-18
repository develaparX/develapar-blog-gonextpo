package service

import (
	"context"
	"develapar-server/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Context-aware test for PaginationService.Paginate
func TestPaginate_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		data          interface{}
		total         int
		query         PaginationQuery
		setupMocks    func(*MockValidationService, *MockErrorWrapper)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name: "successful pagination with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_paginate_123"),
			data: []string{"item1", "item2", "item3"},
			total: 3,
			query: PaginationQuery{
				Page:    1,
				Limit:   10,
				Offset:  0,
				SortBy:  "created_at",
				SortDir: "desc",
			},
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 1, 10).Return(nil).Once()
			},
		},
		{
			name: "context cancellation during pagination",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data: []string{"item1", "item2", "item3"},
			total: 3,
			query: PaginationQuery{
				Page:    1,
				Limit:   10,
				Offset:  0,
				SortBy:  "created_at",
				SortDir: "desc",
			},
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.cancelCtx"), "Pagination operation timed out").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Pagination operation timed out",
				}).Once()
			},
			expectCancel: true,
		},
		{
			name: "context timeout during pagination",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			data: []string{"item1", "item2", "item3"},
			total: 3,
			query: PaginationQuery{
				Page:    1,
				Limit:   10,
				Offset:  0,
				SortBy:  "created_at",
				SortDir: "desc",
			},
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "Pagination operation timed out").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Pagination operation timed out",
				}).Once()
			},
			expectTimeout: true,
		},
		{
			name: "pagination with validation error",
			ctx:  context.WithValue(context.Background(), "request_id", "req_validation_error"),
			data: []string{"item1", "item2", "item3"},
			total: 3,
			query: PaginationQuery{
				Page:    0, // Invalid page
				Limit:   10,
				Offset:  0,
				SortBy:  "created_at",
				SortDir: "desc",
			},
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 0, 10).Return(&utils.AppError{
					Code:    utils.ErrValidation,
					Message: "Pagination validation failed",
				}).Once()
			},
			expectedError: "Pagination validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockValidationService := new(MockValidationService)
			mockErrorWrapper := new(MockErrorWrapper)
			paginationService := NewPaginationService(mockValidationService, mockErrorWrapper)

			// Setup mock expectations
			tt.setupMocks(mockValidationService, mockErrorWrapper)

			// Execute test
			result, err := paginationService.Paginate(tt.ctx, tt.data, tt.total, tt.query)

			// Verify results
			if tt.expectCancel {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectTimeout {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedError)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result.Data)
				assert.Equal(t, "req_paginate_123", result.RequestID)
				assert.Equal(t, 1, result.Metadata.Page)
				assert.Equal(t, 10, result.Metadata.Limit)
				assert.Equal(t, 3, result.Metadata.Total)
			}

			// Verify mock expectations
			mockValidationService.AssertExpectations(t)
			mockErrorWrapper.AssertExpectations(t)
		})
	}
}

// Context-aware test for PaginationService.ParseQuery
func TestParseQuery_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		page          int
		limit         int
		sortBy        string
		sortDir       string
		setupMocks    func(*MockValidationService, *MockErrorWrapper)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:    "successful query parsing with context",
			ctx:     context.WithValue(context.Background(), "request_id", "req_parse_123"),
			page:    1,
			limit:   10,
			sortBy:  "created_at",
			sortDir: "desc",
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 1, 10).Return(nil).Once()
			},
		},
		{
			name: "context cancellation during query parsing",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			page:    1,
			limit:   10,
			sortBy:  "created_at",
			sortDir: "desc",
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.cancelCtx"), "Query parsing timed out").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Query parsing timed out",
				}).Once()
			},
			expectCancel: true,
		},
		{
			name: "context timeout during query parsing",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			page:    1,
			limit:   10,
			sortBy:  "created_at",
			sortDir: "desc",
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "Query parsing timed out").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Query parsing timed out",
				}).Once()
			},
			expectTimeout: true,
		},
		{
			name:    "query parsing with default values",
			ctx:     context.WithValue(context.Background(), "request_id", "req_defaults"),
			page:    0, // Should default to 1
			limit:   0, // Should default to service default
			sortBy:  "",
			sortDir: "",
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 1, 10).Return(nil).Once()
			},
		},
		{
			name:    "query parsing with invalid sort direction",
			ctx:     context.WithValue(context.Background(), "request_id", "req_invalid_sort"),
			page:    1,
			limit:   10,
			sortBy:  "created_at",
			sortDir: "invalid",
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("ValidationError", mock.AnythingOfType("*context.valueCtx"), "sort_dir", "Sort direction must be 'asc' or 'desc'").Return(&utils.AppError{
					Code:    utils.ErrValidation,
					Message: "Sort direction must be 'asc' or 'desc'",
				}).Once()
			},
			expectedError: "Sort direction must be 'asc' or 'desc'",
		},
		{
			name:    "query parsing with validation error",
			ctx:     context.WithValue(context.Background(), "request_id", "req_validation_error"),
			page:    -1,
			limit:   200,
			sortBy:  "created_at",
			sortDir: "desc",
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 1, 100).Return(&utils.AppError{
					Code:    utils.ErrValidation,
					Message: "Pagination validation failed",
				}).Once()
			},
			expectedError: "Pagination validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockValidationService := new(MockValidationService)
			mockErrorWrapper := new(MockErrorWrapper)
			paginationService := NewPaginationService(mockValidationService, mockErrorWrapper)

			// Setup mock expectations
			tt.setupMocks(mockValidationService, mockErrorWrapper)

			// Execute test
			result, err := paginationService.ParseQuery(tt.ctx, tt.page, tt.limit, tt.sortBy, tt.sortDir)

			// Verify results
			if tt.expectCancel {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectTimeout {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedError)
			} else {
				assert.Nil(t, err)
				assert.True(t, result.Page >= 1)
				assert.True(t, result.Limit >= 1)
				assert.True(t, result.Limit <= 100)
				assert.Contains(t, []string{"asc", "desc"}, result.SortDir)
			}

			// Verify mock expectations
			mockValidationService.AssertExpectations(t)
			mockErrorWrapper.AssertExpectations(t)
		})
	}
}

// Context-aware test for PaginationService.BuildMetadata
func TestBuildMetadata_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		total         int
		page          int
		limit         int
		setupMocks    func(*MockValidationService, *MockErrorWrapper)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:  "successful metadata building with context",
			ctx:   context.WithValue(context.Background(), "request_id", "req_metadata_123"),
			total: 25,
			page:  2,
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 2, 10).Return(nil).Once()
			},
		},
		{
			name: "context cancellation during metadata building",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			total: 25,
			page:  2,
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.cancelCtx"), "Metadata building timed out").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Metadata building timed out",
				}).Once()
			},
			expectCancel: true,
		},
		{
			name: "context timeout during metadata building",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			total: 25,
			page:  2,
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "Metadata building timed out").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Metadata building timed out",
				}).Once()
			},
			expectTimeout: true,
		},
		{
			name:  "metadata building with validation error",
			ctx:   context.WithValue(context.Background(), "request_id", "req_validation_error"),
			total: 25,
			page:  0, // Invalid page
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 0, 10).Return(&utils.AppError{
					Code:    utils.ErrValidation,
					Message: "Pagination validation failed",
				}).Once()
			},
			expectedError: "Pagination validation failed",
		},
		{
			name:  "metadata building with zero total",
			ctx:   context.WithValue(context.Background(), "request_id", "req_zero_total"),
			total: 0,
			page:  1,
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 1, 10).Return(nil).Once()
			},
		},
		{
			name:  "metadata building for last page",
			ctx:   context.WithValue(context.Background(), "request_id", "req_last_page"),
			total: 25,
			page:  3,
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 3, 10).Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockValidationService := new(MockValidationService)
			mockErrorWrapper := new(MockErrorWrapper)
			paginationService := NewPaginationService(mockValidationService, mockErrorWrapper)

			// Setup mock expectations
			tt.setupMocks(mockValidationService, mockErrorWrapper)

			// Execute test
			result, err := paginationService.BuildMetadata(tt.ctx, tt.total, tt.page, tt.limit)

			// Verify results
			if tt.expectCancel {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectTimeout {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedError)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.page, result.Page)
				assert.Equal(t, tt.limit, result.Limit)
				assert.Equal(t, tt.total, result.Total)
				
				// Verify calculated fields
				expectedTotalPages := (tt.total + tt.limit - 1) / tt.limit
				if expectedTotalPages == 0 {
					expectedTotalPages = 1
				}
				assert.Equal(t, expectedTotalPages, result.TotalPages)
				
				expectedHasNext := tt.page < expectedTotalPages
				assert.Equal(t, expectedHasNext, result.HasNext)
				
				expectedHasPrev := tt.page > 1
				assert.Equal(t, expectedHasPrev, result.HasPrev)
				
				// Verify context information
				if tt.ctx.Value("request_id") != nil {
					assert.Equal(t, tt.ctx.Value("request_id").(string), result.RequestID)
				}
				assert.NotZero(t, result.ProcessedAt)
			}

			// Verify mock expectations
			mockValidationService.AssertExpectations(t)
			mockErrorWrapper.AssertExpectations(t)
		})
	}
}

// Context-aware test for PaginationService.ValidatePagination
func TestValidatePagination_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		page          int
		limit         int
		setupMocks    func(*MockValidationService, *MockErrorWrapper)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:  "successful pagination validation with context",
			ctx:   context.WithValue(context.Background(), "request_id", "req_validate_123"),
			page:  1,
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 1, 10).Return(nil).Once()
			},
		},
		{
			name: "context cancellation during validation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			page:  1,
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.cancelCtx"), "Pagination validation timed out").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Pagination validation timed out",
				}).Once()
			},
			expectCancel: true,
		},
		{
			name: "context timeout during validation",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			page:  1,
			limit: 10,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "Pagination validation timed out").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Pagination validation timed out",
				}).Once()
			},
			expectTimeout: true,
		},
		{
			name:  "validation error with context",
			ctx:   context.WithValue(context.Background(), "request_id", "req_validation_error"),
			page:  0,
			limit: -1,
			setupMocks: func(mockVS *MockValidationService, mockEW *MockErrorWrapper) {
				mockVS.On("ValidatePagination", mock.AnythingOfType("*context.valueCtx"), 0, -1).Return(&utils.AppError{
					Code:    utils.ErrValidation,
					Message: "Pagination validation failed",
				}).Once()
			},
			expectedError: "Pagination validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockValidationService := new(MockValidationService)
			mockErrorWrapper := new(MockErrorWrapper)
			paginationService := NewPaginationService(mockValidationService, mockErrorWrapper)

			// Setup mock expectations
			tt.setupMocks(mockValidationService, mockErrorWrapper)

			// Execute test
			err := paginationService.ValidatePagination(tt.ctx, tt.page, tt.limit)

			// Verify results
			if tt.expectCancel {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectTimeout {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedError)
			} else {
				assert.Nil(t, err)
			}

			// Verify mock expectations
			mockValidationService.AssertExpectations(t)
			mockErrorWrapper.AssertExpectations(t)
		})
	}
}