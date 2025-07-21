// Category API Service
// This file provides methods for category management operations

import type { AxiosResponse } from 'axios';
import { apiClient } from './apiClient';
import type {
    ApiResponse,
    Category,
    CreateCategoryRequest,
    UpdateCategoryRequest,
    PaginatedResponse,
    PaginationParams
} from './types';

// ============================================================================
// Response Handler Utility
// ============================================================================

class CategoryApiError extends Error {
    constructor(
        public code: string,
        message: string,
        public statusCode?: number,
        public details?: Record<string, any>
    ) {
        super(message);
        this.name = 'CategoryApiError';
    }
}

function handleResponse<T>(response: AxiosResponse<ApiResponse<T>>): T {
    if (response.data.success && response.data.data !== undefined) {
        return response.data.data;
    }

    const error = response.data.error;
    throw new CategoryApiError(
        error?.code || 'UNKNOWN',
        error?.message || 'Unknown error occurred',
        response.status,
        error?.details
    );
}

function handlePaginatedResponse<T>(response: AxiosResponse<ApiResponse<T>>): PaginatedResponse<T> {
    if (response.data.success && response.data.data !== undefined) {
        return {
            data: response.data.data,
            pagination: response.data.pagination
        };
    }

    const error = response.data.error;
    throw new CategoryApiError(
        error?.code || 'UNKNOWN',
        error?.message || 'Unknown error occurred',
        response.status,
        error?.details
    );
}

// ============================================================================
// Category API Service Class
// ============================================================================

export class CategoryApi {

    // ========================================================================
    // Create Category
    // ========================================================================

    /**
     * Create a new category
     * @param categoryData - Category creation data
     * @returns Promise<Category> - Created category
     * @throws CategoryApiError - If creation fails
     */
    async createCategory(categoryData: CreateCategoryRequest): Promise<Category> {
        try {
            const response = await apiClient.post<ApiResponse<Category>>(
                '/categories',
                categoryData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CategoryApiError(
                    error.response.data?.error?.code || 'CREATE_FAILED',
                    error.response.data?.error?.message || 'Failed to create category',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CategoryApiError(
                'NETWORK_ERROR',
                'Network error occurred while creating category'
            );
        }
    }

    // ========================================================================
    // Get All Categories
    // ========================================================================

    /**
     * Get all categories (non-paginated)
     * @returns Promise<Category[]> - Array of all categories
     * @throws CategoryApiError - If retrieval fails
     */
    async getAllCategories(): Promise<Category[]> {
        try {
            const response = await apiClient.get<ApiResponse<Category[]>>('/categories');
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CategoryApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch categories',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CategoryApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching categories'
            );
        }
    }

    /**
     * Get all categories with pagination
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<Category[]>> - Paginated categories
     * @throws CategoryApiError - If retrieval fails
     */
    async getAllCategoriesPaginated(params: PaginationParams = {}): Promise<PaginatedResponse<Category[]>> {
        try {
            const { page = 1, limit = 10 } = params;
            const response = await apiClient.get<ApiResponse<Category[]>>(
                '/categories/paginated',
                {
                    params: { page, limit }
                }
            );
            return handlePaginatedResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CategoryApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || 'Failed to fetch paginated categories',
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CategoryApiError(
                'NETWORK_ERROR',
                'Network error occurred while fetching paginated categories'
            );
        }
    }

    // ========================================================================
    // Get Category by ID
    // ========================================================================

    /**
     * Get a specific category by ID
     * @param id - Category ID
     * @returns Promise<Category> - Category data
     * @throws CategoryApiError - If retrieval fails or category not found
     */
    async getCategoryById(id: number): Promise<Category> {
        try {
            const response = await apiClient.get<ApiResponse<Category>>(`/categories/${id}`);
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CategoryApiError(
                    error.response.data?.error?.code || 'FETCH_FAILED',
                    error.response.data?.error?.message || `Failed to fetch category with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CategoryApiError(
                'NETWORK_ERROR',
                `Network error occurred while fetching category with ID ${id}`
            );
        }
    }

    // ========================================================================
    // Update Category
    // ========================================================================

    /**
     * Update an existing category
     * @param id - Category ID to update
     * @param categoryData - Updated category data
     * @returns Promise<Category> - Updated category
     * @throws CategoryApiError - If update fails or category not found
     */
    async updateCategory(id: number, categoryData: UpdateCategoryRequest): Promise<Category> {
        try {
            const response = await apiClient.put<ApiResponse<Category>>(
                `/categories/${id}`,
                categoryData
            );
            return handleResponse(response);
        } catch (error: any) {
            if (error.response) {
                throw new CategoryApiError(
                    error.response.data?.error?.code || 'UPDATE_FAILED',
                    error.response.data?.error?.message || `Failed to update category with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CategoryApiError(
                'NETWORK_ERROR',
                `Network error occurred while updating category with ID ${id}`
            );
        }
    }

    // ========================================================================
    // Delete Category
    // ========================================================================

    /**
     * Delete a category
     * @param id - Category ID to delete
     * @returns Promise<void> - Resolves when deletion is complete
     * @throws CategoryApiError - If deletion fails or category not found
     */
    async deleteCategory(id: number): Promise<void> {
        try {
            await apiClient.delete(`/categories/${id}`);
        } catch (error: any) {
            if (error.response) {
                throw new CategoryApiError(
                    error.response.data?.error?.code || 'DELETE_FAILED',
                    error.response.data?.error?.message || `Failed to delete category with ID ${id}`,
                    error.response.status,
                    error.response.data?.error?.details
                );
            }
            throw new CategoryApiError(
                'NETWORK_ERROR',
                `Network error occurred while deleting category with ID ${id}`
            );
        }
    }
}

// ============================================================================
// Export Singleton Instance
// ============================================================================

export const categoryApi = new CategoryApi();
export default categoryApi;