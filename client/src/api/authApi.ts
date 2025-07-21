// Authentication API Service
// This file provides authentication-related API methods including login, register, refresh, and logout

import type { AxiosResponse } from 'axios';
import { apiClient } from './apiClient';
import { tokenManager } from './tokenManager';
import { ResponseHandler, ApiError } from './errorHandler';
import type {
    ApiResponse,
    LoginRequest,
    LoginResponse,
    RegisterRequest,
    User
} from './types';
import { ApiErrorCode } from './types';

// ============================================================================
// Token Validation Utilities
// ============================================================================

interface TokenPayload {
    exp: number;
    iat: number;
    user_id: number;
    email: string;
    role: string;
}

class TokenValidator {
    /**
     * Decode JWT token payload (without verification)
     * Note: This is for client-side expiry checking only, not security validation
     */
    static decodeToken(token: string): TokenPayload | null {
        try {
            const parts = token.split('.');
            if (parts.length !== 3) {
                return null;
            }

            const payload = parts[1];
            const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'));
            return JSON.parse(decoded) as TokenPayload;
        } catch (error) {
            return null;
        }
    }

    /**
     * Check if token is expired
     */
    static isTokenExpired(token: string): boolean {
        const payload = this.decodeToken(token);
        if (!payload) {
            return true;
        }

        const currentTime = Math.floor(Date.now() / 1000);
        return payload.exp < currentTime;
    }

    /**
     * Check if token expires within the next N seconds
     */
    static isTokenExpiringSoon(token: string, thresholdSeconds: number = 300): boolean {
        const payload = this.decodeToken(token);
        if (!payload) {
            return true;
        }

        const currentTime = Math.floor(Date.now() / 1000);
        return payload.exp < (currentTime + thresholdSeconds);
    }

    /**
     * Get token expiry time as Date
     */
    static getTokenExpiryDate(token: string): Date | null {
        const payload = this.decodeToken(token);
        if (!payload) {
            return null;
        }

        return new Date(payload.exp * 1000);
    }

    /**
     * Extract user information from token
     */
    static getUserFromToken(token: string): Partial<User> | null {
        const payload = this.decodeToken(token);
        if (!payload) {
            return null;
        }

        return {
            id: payload.user_id,
            email: payload.email,
            role: payload.role
        };
    }
}

// ============================================================================
// Authentication Service Class
// ============================================================================

export class AuthApi {
    private static instance: AuthApi;

    private constructor() { }

    /**
     * Get singleton instance
     */
    static getInstance(): AuthApi {
        if (!AuthApi.instance) {
            AuthApi.instance = new AuthApi();
        }
        return AuthApi.instance;
    }

    // ============================================================================
    // Authentication Methods
    // ============================================================================

    /**
     * User login with credentials
     */
    async login(credentials: LoginRequest): Promise<LoginResponse> {
        try {
            // Validate credentials
            this.validateLoginCredentials(credentials);

            const response: AxiosResponse<ApiResponse<LoginResponse>> = await apiClient.post(
                '/auth/login',
                credentials
            );

            const loginData = ResponseHandler.handleSuccess(response);

            // Store tokens securely
            tokenManager.setTokens(loginData.access_token, loginData.refresh_token);

            return loginData;
        } catch (error) {
            throw ResponseHandler.handleError(error);
        }
    }

    /**
     * User registration
     */
    async register(userData: RegisterRequest): Promise<User> {
        try {
            // Validate registration data
            this.validateRegistrationData(userData);

            const response: AxiosResponse<ApiResponse<User>> = await apiClient.post(
                '/auth/register',
                userData
            );

            return ResponseHandler.handleSuccess(response);
        } catch (error) {
            throw ResponseHandler.handleError(error);
        }
    }

    /**
     * Refresh access token using refresh token
     */
    async refreshToken(): Promise<LoginResponse> {
        try {
            const refreshToken = tokenManager.getRefreshToken();

            if (!refreshToken) {
                throw new ApiError(
                    ApiErrorCode.UNAUTHORIZED,
                    'No refresh token available',
                    401
                );
            }

            // Use a separate request without the main interceptors to avoid loops
            const response: AxiosResponse<ApiResponse<LoginResponse>> = await apiClient.post(
                '/auth/refresh',
                {},
                {
                    headers: {
                        'Authorization': `Bearer ${refreshToken}`
                    }
                }
            );

            const tokenData = ResponseHandler.handleSuccess(response);

            // Update stored tokens
            tokenManager.setTokens(tokenData.access_token, tokenData.refresh_token);

            return tokenData;
        } catch (error) {
            // If refresh fails, clear all tokens
            tokenManager.clearTokens();
            throw ResponseHandler.handleError(error);
        }
    }

    /**
     * User logout
     */
    async logout(): Promise<void> {
        try {
            // Attempt to notify server about logout (optional)
            try {
                await apiClient.post('/auth/logout');
            } catch (error) {
                // Ignore server errors during logout, still clear local tokens
                console.warn('Server logout failed, clearing local tokens anyway:', error);
            }
        } finally {
            // Always clear local tokens regardless of server response
            tokenManager.clearTokens();
        }
    }

    // ============================================================================
    // Token Management Methods
    // ============================================================================

    /**
     * Check if user is currently authenticated
     */
    isAuthenticated(): boolean {
        const accessToken = tokenManager.getAccessToken();
        const refreshToken = tokenManager.getRefreshToken();

        if (!accessToken || !refreshToken) {
            return false;
        }

        // Check if access token is valid (not expired)
        if (TokenValidator.isTokenExpired(accessToken)) {
            // If access token is expired, check if refresh token is valid
            return !TokenValidator.isTokenExpired(refreshToken);
        }

        return true;
    }

    /**
     * Check if access token needs refresh
     */
    needsTokenRefresh(): boolean {
        const accessToken = tokenManager.getAccessToken();

        if (!accessToken) {
            return false;
        }

        // Check if token expires within 5 minutes
        return TokenValidator.isTokenExpiringSoon(accessToken, 300);
    }

    /**
     * Get current user information from token
     */
    getCurrentUser(): Partial<User> | null {
        const accessToken = tokenManager.getAccessToken();

        if (!accessToken || TokenValidator.isTokenExpired(accessToken)) {
            return null;
        }

        return TokenValidator.getUserFromToken(accessToken);
    }

    /**
     * Get access token expiry date
     */
    getTokenExpiryDate(): Date | null {
        const accessToken = tokenManager.getAccessToken();

        if (!accessToken) {
            return null;
        }

        return TokenValidator.getTokenExpiryDate(accessToken);
    }

    /**
     * Validate access token and refresh if needed
     */
    async ensureValidToken(): Promise<boolean> {
        try {
            if (!this.isAuthenticated()) {
                return false;
            }

            if (this.needsTokenRefresh()) {
                await this.refreshToken();
            }

            return true;
        } catch (error) {
            return false;
        }
    }

    // ============================================================================
    // Validation Methods
    // ============================================================================

    private validateLoginCredentials(credentials: LoginRequest): void {
        if (!credentials.identifier || credentials.identifier.trim() === '') {
            throw new ApiError(
                ApiErrorCode.VALIDATION_ERROR,
                'Email or username is required',
                400,
                { field: 'identifier', message: 'Email or username is required' }
            );
        }

        if (!credentials.password || credentials.password.length < 1) {
            throw new ApiError(
                ApiErrorCode.VALIDATION_ERROR,
                'Password is required',
                400,
                { field: 'password', message: 'Password is required' }
            );
        }

        // Basic email validation if identifier looks like an email
        if (credentials.identifier.includes('@')) {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(credentials.identifier)) {
                throw new ApiError(
                    ApiErrorCode.VALIDATION_ERROR,
                    'Invalid email format',
                    400,
                    { field: 'identifier', message: 'Invalid email format' }
                );
            }
        }
    }

    private validateRegistrationData(userData: RegisterRequest): void {
        const errors: Record<string, string> = {};

        // Validate name
        if (!userData.name || userData.name.trim() === '') {
            errors.name = 'Name is required';
        } else if (userData.name.trim().length < 2) {
            errors.name = 'Name must be at least 2 characters long';
        }

        // Validate email
        if (!userData.email || userData.email.trim() === '') {
            errors.email = 'Email is required';
        } else {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(userData.email)) {
                errors.email = 'Invalid email format';
            }
        }

        // Validate password
        if (!userData.password) {
            errors.password = 'Password is required';
        } else if (userData.password.length < 8) {
            errors.password = 'Password must be at least 8 characters long';
        } else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(userData.password)) {
            errors.password = 'Password must contain at least one uppercase letter, one lowercase letter, and one number';
        }

        if (Object.keys(errors).length > 0) {
            throw new ApiError(
                ApiErrorCode.VALIDATION_ERROR,
                'Validation failed',
                400,
                errors
            );
        }
    }
}

// ============================================================================
// Export singleton instance and utilities
// ============================================================================

export const authApi = AuthApi.getInstance();
export { TokenValidator };
export default authApi;