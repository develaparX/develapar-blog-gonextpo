package service

import (
	"context"
	"develapar-server/utils"
	"math"
	"time"
)

// PaginationQuery represents pagination parameters
type PaginationQuery struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
	SortBy   string `json:"sort_by,omitempty"`
	SortDir  string `json:"sort_dir,omitempty"`
}

// PaginationResult represents the result of a paginated query
type PaginationResult struct {
	Data       interface{}         `json:"data"`
	Metadata   PaginationMetadata  `json:"pagination"`
	RequestID  string              `json:"request_id,omitempty"`
}

// PaginationMetadata contains pagination information
type PaginationMetadata struct {
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	Total        int    `json:"total"`
	TotalPages   int    `json:"total_pages"`
	HasNext      bool   `json:"has_next"`
	HasPrev      bool   `json:"has_prev"`
	RequestID    string `json:"request_id,omitempty"`
	ProcessedAt  time.Time `json:"processed_at"`
}

// PaginationService interface with context support
type PaginationService interface {
	// Paginate creates a pagination result with context support
	Paginate(ctx context.Context, data interface{}, total int, query PaginationQuery) (PaginationResult, *utils.AppError)
	
	// ValidatePagination validates pagination parameters with context
	ValidatePagination(ctx context.Context, page, limit int) *utils.AppError
	
	// BuildMetadata creates pagination metadata with context information
	BuildMetadata(ctx context.Context, total, page, limit int) (PaginationMetadata, *utils.AppError)
	
	// ParseQuery parses and validates pagination query parameters
	ParseQuery(ctx context.Context, page, limit int, sortBy, sortDir string) (PaginationQuery, *utils.AppError)
	
	// CalculateOffset calculates the offset for database queries
	CalculateOffset(page, limit int) int
}

// paginationService implements PaginationService
type paginationService struct {
	validationService ValidationService
	errorWrapper      utils.ErrorWrapper
	defaultLimit      int
	maxLimit          int
}

// Paginate creates a pagination result with context support
func (ps *paginationService) Paginate(ctx context.Context, data interface{}, total int, query PaginationQuery) (PaginationResult, *utils.AppError) {
	// Check context timeout
	select {
	case <-ctx.Done():
		return PaginationResult{}, ps.errorWrapper.TimeoutError(ctx, "Pagination operation timed out")
	default:
	}

	// Build metadata
	metadata, err := ps.BuildMetadata(ctx, total, query.Page, query.Limit)
	if err != nil {
		return PaginationResult{}, err
	}

	// Get request ID from context
	requestID := ps.getRequestIDFromContext(ctx)
	
	result := PaginationResult{
		Data:      data,
		Metadata:  metadata,
		RequestID: requestID,
	}

	return result, nil
}

// ValidatePagination validates pagination parameters with context
func (ps *paginationService) ValidatePagination(ctx context.Context, page, limit int) *utils.AppError {
	// Check context timeout
	select {
	case <-ctx.Done():
		return ps.errorWrapper.TimeoutError(ctx, "Pagination validation timed out")
	default:
	}

	// Use the existing validation service
	return ps.validationService.ValidatePagination(ctx, page, limit)
}

// BuildMetadata creates pagination metadata with context information
func (ps *paginationService) BuildMetadata(ctx context.Context, total, page, limit int) (PaginationMetadata, *utils.AppError) {
	// Check context timeout
	select {
	case <-ctx.Done():
		return PaginationMetadata{}, ps.errorWrapper.TimeoutError(ctx, "Metadata building timed out")
	default:
	}

	// Validate pagination parameters first
	if err := ps.ValidatePagination(ctx, page, limit); err != nil {
		return PaginationMetadata{}, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	// Calculate has next/prev
	hasNext := page < totalPages
	hasPrev := page > 1

	// Get request ID from context
	requestID := ps.getRequestIDFromContext(ctx)

	metadata := PaginationMetadata{
		Page:        page,
		Limit:       limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
		RequestID:   requestID,
		ProcessedAt: time.Now(),
	}

	return metadata, nil
}

// ParseQuery parses and validates pagination query parameters
func (ps *paginationService) ParseQuery(ctx context.Context, page, limit int, sortBy, sortDir string) (PaginationQuery, *utils.AppError) {
	// Check context timeout
	select {
	case <-ctx.Done():
		return PaginationQuery{}, ps.errorWrapper.TimeoutError(ctx, "Query parsing timed out")
	default:
	}

	// Set defaults if not provided
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = ps.defaultLimit
	}
	if limit > ps.maxLimit {
		limit = ps.maxLimit
	}

	// Validate sort direction
	if sortDir != "" && sortDir != "asc" && sortDir != "desc" {
		return PaginationQuery{}, ps.errorWrapper.ValidationError(ctx, "sort_dir", "Sort direction must be 'asc' or 'desc'")
	}

	// Set default sort direction
	if sortDir == "" {
		sortDir = "desc"
	}

	// Validate pagination parameters
	if err := ps.ValidatePagination(ctx, page, limit); err != nil {
		return PaginationQuery{}, err
	}

	// Calculate offset
	offset := ps.CalculateOffset(page, limit)

	query := PaginationQuery{
		Page:    page,
		Limit:   limit,
		Offset:  offset,
		SortBy:  sortBy,
		SortDir: sortDir,
	}

	return query, nil
}

// CalculateOffset calculates the offset for database queries
func (ps *paginationService) CalculateOffset(page, limit int) int {
	if page <= 1 {
		return 0
	}
	return (page - 1) * limit
}

// getRequestIDFromContext extracts request ID from context
func (ps *paginationService) getRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

// NewPaginationService creates a new pagination service instance
func NewPaginationService(validationService ValidationService, errorWrapper utils.ErrorWrapper) PaginationService {
	return &paginationService{
		validationService: validationService,
		errorWrapper:      errorWrapper,
		defaultLimit:      10,  // Default page size
		maxLimit:          100, // Maximum page size
	}
}