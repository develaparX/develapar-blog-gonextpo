package service

import (
	"context"
	"develapar-server/utils"
	"testing"
	"time"
)

func TestPaginationService_ValidatePagination(t *testing.T) {
	// Setup
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)
	paginationService := NewPaginationService(validationService, errorWrapper)

	tests := []struct {
		name    string
		ctx     context.Context
		page    int
		limit   int
		wantErr bool
	}{
		{
			name:    "valid pagination",
			ctx:     context.WithValue(context.Background(), "request_id", "req_300"),
			page:    1,
			limit:   10,
			wantErr: false,
		},
		{
			name:    "invalid page - zero",
			ctx:     context.WithValue(context.Background(), "request_id", "req_301"),
			page:    0,
			limit:   10,
			wantErr: true,
		},
		{
			name:    "invalid page - negative",
			ctx:     context.WithValue(context.Background(), "request_id", "req_302"),
			page:    -1,
			limit:   10,
			wantErr: true,
		},
		{
			name:    "invalid limit - zero",
			ctx:     context.WithValue(context.Background(), "request_id", "req_303"),
			page:    1,
			limit:   0,
			wantErr: true,
		},
		{
			name:    "invalid limit - negative",
			ctx:     context.WithValue(context.Background(), "request_id", "req_304"),
			page:    1,
			limit:   -5,
			wantErr: true,
		},
		{
			name:    "invalid limit - too large",
			ctx:     context.WithValue(context.Background(), "request_id", "req_305"),
			page:    1,
			limit:   1000,
			wantErr: true,
		},
		{
			name: "context timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				defer cancel()
				time.Sleep(2 * time.Nanosecond) // Ensure timeout
				return ctx
			}(),
			page:    1,
			limit:   10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := paginationService.ValidatePagination(tt.ctx, tt.page, tt.limit)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePagination() expected error, got nil")
					return
				}
				
				// Check error type
				appErr := err
				if appErr.Code != utils.ErrValidation && appErr.Code != utils.ErrTimeout {
					t.Errorf("ValidatePagination() expected validation or timeout error, got %s", appErr.Code)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePagination() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestPaginationService_BuildMetadata(t *testing.T) {
	// Setup
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)
	paginationService := NewPaginationService(validationService, errorWrapper)

	tests := []struct {
		name           string
		ctx            context.Context
		total          int
		page           int
		limit          int
		wantErr        bool
		expectedPages  int
		expectedHasNext bool
		expectedHasPrev bool
	}{
		{
			name:            "first page with more data",
			ctx:             context.WithValue(context.Background(), "request_id", "req_400"),
			total:           50,
			page:            1,
			limit:           10,
			wantErr:         false,
			expectedPages:   5,
			expectedHasNext: true,
			expectedHasPrev: false,
		},
		{
			name:            "middle page",
			ctx:             context.WithValue(context.Background(), "request_id", "req_401"),
			total:           50,
			page:            3,
			limit:           10,
			wantErr:         false,
			expectedPages:   5,
			expectedHasNext: true,
			expectedHasPrev: true,
		},
		{
			name:            "last page",
			ctx:             context.WithValue(context.Background(), "request_id", "req_402"),
			total:           50,
			page:            5,
			limit:           10,
			wantErr:         false,
			expectedPages:   5,
			expectedHasNext: false,
			expectedHasPrev: true,
		},
		{
			name:            "single page",
			ctx:             context.WithValue(context.Background(), "request_id", "req_403"),
			total:           5,
			page:            1,
			limit:           10,
			wantErr:         false,
			expectedPages:   1,
			expectedHasNext: false,
			expectedHasPrev: false,
		},
		{
			name:    "invalid pagination parameters",
			ctx:     context.WithValue(context.Background(), "request_id", "req_404"),
			total:   50,
			page:    0,
			limit:   10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := paginationService.BuildMetadata(tt.ctx, tt.total, tt.page, tt.limit)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("BuildMetadata() expected error, got nil")
					return
				}
			} else {
				if err != nil {
					t.Errorf("BuildMetadata() expected no error, got %v", err)
					return
				}
				
				// Verify metadata
				if metadata.Page != tt.page {
					t.Errorf("BuildMetadata() page = %v, want %v", metadata.Page, tt.page)
				}
				if metadata.Limit != tt.limit {
					t.Errorf("BuildMetadata() limit = %v, want %v", metadata.Limit, tt.limit)
				}
				if metadata.Total != tt.total {
					t.Errorf("BuildMetadata() total = %v, want %v", metadata.Total, tt.total)
				}
				if metadata.TotalPages != tt.expectedPages {
					t.Errorf("BuildMetadata() totalPages = %v, want %v", metadata.TotalPages, tt.expectedPages)
				}
				if metadata.HasNext != tt.expectedHasNext {
					t.Errorf("BuildMetadata() hasNext = %v, want %v", metadata.HasNext, tt.expectedHasNext)
				}
				if metadata.HasPrev != tt.expectedHasPrev {
					t.Errorf("BuildMetadata() hasPrev = %v, want %v", metadata.HasPrev, tt.expectedHasPrev)
				}
				
				// Check request ID
				if requestID, ok := tt.ctx.Value("request_id").(string); ok {
					if metadata.RequestID != requestID {
						t.Errorf("BuildMetadata() requestID = %v, want %v", metadata.RequestID, requestID)
					}
				}
			}
		})
	}
}

func TestPaginationService_ParseQuery(t *testing.T) {
	// Setup
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)
	paginationService := NewPaginationService(validationService, errorWrapper)

	tests := []struct {
		name           string
		ctx            context.Context
		page           int
		limit          int
		sortBy         string
		sortDir        string
		wantErr        bool
		expectedPage   int
		expectedLimit  int
		expectedOffset int
		expectedSortDir string
	}{
		{
			name:            "valid query with all parameters",
			ctx:             context.WithValue(context.Background(), "request_id", "req_500"),
			page:            2,
			limit:           20,
			sortBy:          "created_at",
			sortDir:         "asc",
			wantErr:         false,
			expectedPage:    2,
			expectedLimit:   20,
			expectedOffset:  20,
			expectedSortDir: "asc",
		},
		{
			name:            "defaults applied for zero values",
			ctx:             context.WithValue(context.Background(), "request_id", "req_501"),
			page:            0,
			limit:           0,
			sortBy:          "",
			sortDir:         "",
			wantErr:         false,
			expectedPage:    1,
			expectedLimit:   10, // default limit
			expectedOffset:  0,
			expectedSortDir: "desc", // default sort direction
		},
		{
			name:            "limit capped at maximum",
			ctx:             context.WithValue(context.Background(), "request_id", "req_502"),
			page:            1,
			limit:           200, // exceeds max limit of 100
			sortBy:          "title",
			sortDir:         "desc",
			wantErr:         false,
			expectedPage:    1,
			expectedLimit:   100, // capped at max
			expectedOffset:  0,
			expectedSortDir: "desc",
		},
		{
			name:    "invalid sort direction",
			ctx:     context.WithValue(context.Background(), "request_id", "req_503"),
			page:    1,
			limit:   10,
			sortBy:  "title",
			sortDir: "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := paginationService.ParseQuery(tt.ctx, tt.page, tt.limit, tt.sortBy, tt.sortDir)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseQuery() expected error, got nil")
					return
				}
			} else {
				if err != nil {
					t.Errorf("ParseQuery() expected no error, got %v", err)
					return
				}
				
				// Verify query
				if query.Page != tt.expectedPage {
					t.Errorf("ParseQuery() page = %v, want %v", query.Page, tt.expectedPage)
				}
				if query.Limit != tt.expectedLimit {
					t.Errorf("ParseQuery() limit = %v, want %v", query.Limit, tt.expectedLimit)
				}
				if query.Offset != tt.expectedOffset {
					t.Errorf("ParseQuery() offset = %v, want %v", query.Offset, tt.expectedOffset)
				}
				if query.SortDir != tt.expectedSortDir {
					t.Errorf("ParseQuery() sortDir = %v, want %v", query.SortDir, tt.expectedSortDir)
				}
				if query.SortBy != tt.sortBy {
					t.Errorf("ParseQuery() sortBy = %v, want %v", query.SortBy, tt.sortBy)
				}
			}
		})
	}
}

func TestPaginationService_CalculateOffset(t *testing.T) {
	// Setup
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)
	paginationService := NewPaginationService(validationService, errorWrapper)

	tests := []struct {
		name           string
		page           int
		limit          int
		expectedOffset int
	}{
		{
			name:           "first page",
			page:           1,
			limit:          10,
			expectedOffset: 0,
		},
		{
			name:           "second page",
			page:           2,
			limit:          10,
			expectedOffset: 10,
		},
		{
			name:           "third page with different limit",
			page:           3,
			limit:          20,
			expectedOffset: 40,
		},
		{
			name:           "zero page",
			page:           0,
			limit:          10,
			expectedOffset: 0,
		},
		{
			name:           "negative page",
			page:           -1,
			limit:          10,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := paginationService.CalculateOffset(tt.page, tt.limit)
			if offset != tt.expectedOffset {
				t.Errorf("CalculateOffset() = %v, want %v", offset, tt.expectedOffset)
			}
		})
	}
}

func TestPaginationService_Paginate(t *testing.T) {
	// Setup
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)
	paginationService := NewPaginationService(validationService, errorWrapper)

	// Test data
	testData := []string{"item1", "item2", "item3"}
	
	tests := []struct {
		name    string
		ctx     context.Context
		data    interface{}
		total   int
		query   PaginationQuery
		wantErr bool
	}{
		{
			name:  "successful pagination",
			ctx:   context.WithValue(context.Background(), "request_id", "req_600"),
			data:  testData,
			total: 30,
			query: PaginationQuery{
				Page:  1,
				Limit: 10,
			},
			wantErr: false,
		},
		{
			name: "context timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				defer cancel()
				time.Sleep(2 * time.Nanosecond) // Ensure timeout
				return ctx
			}(),
			data:  testData,
			total: 30,
			query: PaginationQuery{
				Page:  1,
				Limit: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := paginationService.Paginate(tt.ctx, tt.data, tt.total, tt.query)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Paginate() expected error, got nil")
					return
				}
			} else {
				if err != nil {
					t.Errorf("Paginate() expected no error, got %v", err)
					return
				}
				
				// Verify result
				if result.Data != tt.data {
					t.Errorf("Paginate() data mismatch")
				}
				
				// Check request ID
				if requestID, ok := tt.ctx.Value("request_id").(string); ok {
					if result.RequestID != requestID {
						t.Errorf("Paginate() requestID = %v, want %v", result.RequestID, requestID)
					}
				}
			}
		})
	}
}