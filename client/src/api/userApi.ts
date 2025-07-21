// User Management API
// This file contains all user-related API operations including CRUD functionality

import { apiClient } from './apiClient';
import { ResponseHandler } from './errorHandler';
import type {
    User,
    UpdateUserRequest,
    ApiResponse,
    PaginationParams,
    PaginatedResponse
} from './types';

// ============================================================================
// User API Service Class
// ============================================================================

export class UserApi {

    // ============================================================================
    // User Retrieval Methods
    // ============================================================================

    /**
     * Get user by ID with comprehensive error handling
     * @param userId - The ID of the user to retrieve
     * @returns Promise<User> - The user data
     * @throws ApiError - When user not found or access denied
     */
    async getUserById(userId: number): Promise<User> {
        try {
            const response = await apiClient.get<ApiResponse<User>>(`/users/${userId}`);
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get all users without pagination
     * @returns Promise<User[]> - Array of all users
     * @throws ApiError - When access denied or server error
     */
    async getAllUsers(): Promise<User[]> {
        try {
            const response = await apiClient.get<ApiResponse<User[]>>('/users');
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get all users with pagination support
     * @param params - Pagination parameters (page, limit)
     * @returns Promise<PaginatedResponse<User[]>> - Paginated user data
     * @throws ApiError - When access denied or server error
     */
    async getAllUsersPaginated(params: PaginationParams = {}): Promise<PaginatedResponse<User[]>> {
        try {
            const { page = 1, limit = 10 } = params;

            const response = await apiClient.get<ApiResponse<User[]>>('/users/paginated', {
                params: { page, limit }
            });

            return ResponseHandler.handlePaginatedSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    // ============================================================================
    // User Modification Methods
    // ============================================================================

    /**
     * Update user profile with authorization checks
     * @param userId - The ID of the user to update
     * @param userData - The user data to update
     * @returns Promise<User> - The updated user data
     * @throws ApiError - When unauthorized, validation fails, or user not found
     */
    async updateUser(userId: number, userData: UpdateUserRequest): Promise<User> {
        try {
            // Validate that we have at least one field to update
            if (!userData || Object.keys(userData).length === 0) {
                throw new Error('No update data provided');
            }

            const response = await apiClient.put<ApiResponse<User>>(`/users/${userId}`, userData);
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Delete user with proper permission validation
     * @param userId - The ID of the user to delete
     * @returns Promise<void> - Resolves when deletion is successful
     * @throws ApiError - When unauthorized, user not found, or deletion fails
     */
    async deleteUser(userId: number): Promise<void> {
        try {
            await apiClient.delete(`/users/${userId}`);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    // ============================================================================
    // User Profile Methods
    // ============================================================================

    /**
     * Get current user profile
     * @returns Promise<User> - The current user's profile data
     * @throws ApiError - When not authenticated or server error
     */
    async getCurrentUser(): Promise<User> {
        try {
            const response = await apiClient.get<ApiResponse<User>>('/users/me');
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Update current user profile
     * @param userData - The user data to update
     * @returns Promise<User> - The updated user data
     * @throws ApiError - When validation fails or server error
     */
    async updateCurrentUser(userData: UpdateUserRequest): Promise<User> {
        try {
            // Validate that we have at least one field to update
            if (!userData || Object.keys(userData).length === 0) {
                throw new Error('No update data provided');
            }

            const response = await apiClient.put<ApiResponse<User>>('/users/me', userData);
            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    // ============================================================================
    // User Search and Filter Methods
    // ============================================================================

    /**
     * Search users by name or email
     * @param query - Search query string
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<User[]>> - Paginated search results
     * @throws ApiError - When search fails or server error
     */
    async searchUsers(query: string, params: PaginationParams = {}): Promise<PaginatedResponse<User[]>> {
        try {
            const { page = 1, limit = 10 } = params;

            const response = await apiClient.get<ApiResponse<User[]>>('/users/search', {
                params: { q: query, page, limit }
            });

            return ResponseHandler.handlePaginatedSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    /**
     * Get users by role
     * @param role - User role to filter by
     * @param params - Pagination parameters
     * @returns Promise<PaginatedResponse<User[]>> - Paginated users by role
     * @throws ApiError - When access denied or server error
     */
    async getUsersByRole(role: string, params: PaginationParams = {}): Promise<PaginatedResponse<User[]>> {
        try {
            const { page = 1, limit = 10 } = params;

            const response = await apiClient.get<ApiResponse<User[]>>(`/users/role/${role}`, {
                params: { page, limit }
            });

            return ResponseHandler.handlePaginatedSuccess(response);
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }

    // ============================================================================
    // Utility Methods
    // ============================================================================

    /**
     * Check if a user exists by ID
     * @param userId - The ID of the user to check
     * @returns Promise<boolean> - True if user exists, false otherwise
     */
    async userExists(userId: number): Promise<boolean> {
        try {
            await this.getUserById(userId);
            return true;
        } catch (error) {
            // If it's a 404 error, user doesn't exist
            if (error instanceof Error && error.message.includes('404')) {
                return false;
            }
            // Re-throw other errors
            throw error;
        }
    }

    /**
     * Validate user permissions for a specific action
     * @param userId - The ID of the user to check permissions for
     * @param action - The action to validate permissions for
     * @returns Promise<boolean> - True if user has permission, false otherwise
     */
    async validateUserPermissions(userId: number, action: string): Promise<boolean> {
        try {
            const response = await apiClient.get<ApiResponse<{ has_permission: boolean }>>(
                `/users/${userId}/permissions/${action}`
            );
            const result = ResponseHandler.handleSuccess(response);
            return result.has_permission;
        } catch (error) {
            ResponseHandler.handleError(error);
        }
    }
}

// ============================================================================
// Export Singleton Instance
// ============================================================================

export const userApi = new UserApi();
export default userApi;