// User API Usage Examples
// This file demonstrates how to use the UserApi methods

import { userApi } from './userApi';
import type { UpdateUserRequest, PaginationParams } from './types';

// ============================================================================
// Basic CRUD Operations Examples
// ============================================================================

export const userApiExamples = {
    // Example: Get user by ID
    async getUserExample() {
        try {
            const user = await userApi.getUserById(1);
            console.log('User retrieved:', user);
            return user;
        } catch (error) {
            console.error('Failed to get user:', error);
            throw error;
        }
    },

    // Example: Get all users
    async getAllUsersExample() {
        try {
            const users = await userApi.getAllUsers();
            console.log('All users retrieved:', users.length);
            return users;
        } catch (error) {
            console.error('Failed to get all users:', error);
            throw error;
        }
    },

    // Example: Get users with pagination
    async getUsersPaginatedExample() {
        try {
            const params: PaginationParams = { page: 1, limit: 10 };
            const result = await userApi.getAllUsersPaginated(params);
            console.log('Paginated users:', result.data.length);
            console.log('Pagination info:', result.pagination);
            return result;
        } catch (error) {
            console.error('Failed to get paginated users:', error);
            throw error;
        }
    },

    // Example: Update user
    async updateUserExample() {
        try {
            const updateData: UpdateUserRequest = {
                name: 'Updated Name',
                email: 'updated@example.com'
            };
            const updatedUser = await userApi.updateUser(1, updateData);
            console.log('User updated:', updatedUser);
            return updatedUser;
        } catch (error) {
            console.error('Failed to update user:', error);
            throw error;
        }
    },

    // Example: Delete user
    async deleteUserExample() {
        try {
            await userApi.deleteUser(1);
            console.log('User deleted successfully');
        } catch (error) {
            console.error('Failed to delete user:', error);
            throw error;
        }
    },

    // Example: Get current user
    async getCurrentUserExample() {
        try {
            const currentUser = await userApi.getCurrentUser();
            console.log('Current user:', currentUser);
            return currentUser;
        } catch (error) {
            console.error('Failed to get current user:', error);
            throw error;
        }
    },

    // Example: Search users
    async searchUsersExample() {
        try {
            const searchResults = await userApi.searchUsers('john', { page: 1, limit: 5 });
            console.log('Search results:', searchResults.data.length);
            return searchResults;
        } catch (error) {
            console.error('Failed to search users:', error);
            throw error;
        }
    },

    // Example: Get users by role
    async getUsersByRoleExample() {
        try {
            const adminUsers = await userApi.getUsersByRole('admin', { page: 1, limit: 10 });
            console.log('Admin users:', adminUsers.data.length);
            return adminUsers;
        } catch (error) {
            console.error('Failed to get users by role:', error);
            throw error;
        }
    },

    // Example: Check if user exists
    async userExistsExample() {
        try {
            const exists = await userApi.userExists(1);
            console.log('User exists:', exists);
            return exists;
        } catch (error) {
            console.error('Failed to check if user exists:', error);
            throw error;
        }
    },

    // Example: Validate user permissions
    async validatePermissionsExample() {
        try {
            const hasPermission = await userApi.validateUserPermissions(1, 'delete_article');
            console.log('User has permission:', hasPermission);
            return hasPermission;
        } catch (error) {
            console.error('Failed to validate permissions:', error);
            throw error;
        }
    }
};

// ============================================================================
// Error Handling Examples
// ============================================================================

export const errorHandlingExamples = {
    // Example: Handle user not found
    async handleUserNotFound() {
        try {
            const user = await userApi.getUserById(999999);
            return user;
        } catch (error) {
            console.log('Expected error for non-existent user:', error);
            // Handle 404 error appropriately
            return null;
        }
    },

    // Example: Handle validation errors
    async handleValidationError() {
        try {
            const invalidData: UpdateUserRequest = {
                email: 'invalid-email' // Invalid email format
            };
            const user = await userApi.updateUser(1, invalidData);
            return user;
        } catch (error) {
            console.log('Expected validation error:', error);
            // Handle validation error appropriately
            return null;
        }
    },

    // Example: Handle authorization errors
    async handleAuthorizationError() {
        try {
            // Attempt to delete user without proper permissions
            await userApi.deleteUser(1);
        } catch (error) {
            console.log('Expected authorization error:', error);
            // Handle 403 error appropriately
        }
    }
};

// ============================================================================
// Usage in React Components
// ============================================================================

export const reactComponentExamples = {
    // Example: Hook for fetching user data
    useUserData: (userId: number) => {
        // This would typically be implemented with React hooks
        // const [user, setUser] = useState<User | null>(null);
        // const [loading, setLoading] = useState(true);
        // const [error, setError] = useState<string | null>(null);

        // useEffect(() => {
        //     const fetchUser = async () => {
        //         try {
        //             setLoading(true);
        //             const userData = await userApi.getUserById(userId);
        //             setUser(userData);
        //             setError(null);
        //         } catch (err) {
        //             setError(err instanceof Error ? err.message : 'Failed to fetch user');
        //             setUser(null);
        //         } finally {
        //             setLoading(false);
        //         }
        //     };

        //     fetchUser();
        // }, [userId]);

        // return { user, loading, error };
    },

    // Example: Form submission handler
    handleUserUpdate: async (userId: number, formData: UpdateUserRequest) => {
        try {
            const updatedUser = await userApi.updateUser(userId, formData);
            console.log('User updated successfully:', updatedUser);
            // Show success message to user
            return updatedUser;
        } catch (error) {
            console.error('Failed to update user:', error);
            // Show error message to user
            throw error;
        }
    }
};

console.log('UserApi examples loaded successfully');
export default userApiExamples;