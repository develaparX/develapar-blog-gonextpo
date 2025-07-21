// Token Management Utilities
// This file provides secure token storage, retrieval, validation, and cleanup utilities

import { ApiError } from './errorHandler';
import { ApiErrorCode } from './types';

// ============================================================================
// Storage Configuration
// ============================================================================

interface StorageConfig {
    useSessionStorage: boolean;
    encryptTokens: boolean;
    tokenPrefix: string;
}

const DEFAULT_CONFIG: StorageConfig = {
    useSessionStorage: false, // Use localStorage by default for persistence
    encryptTokens: false, // Basic implementation without encryption
    tokenPrefix: 'auth_'
};

// ============================================================================
// Token Storage Keys
// ============================================================================

const STORAGE_KEYS = {
    ACCESS_TOKEN: 'access_token',
    REFRESH_TOKEN: 'refresh_token',
    TOKEN_EXPIRY: 'token_expiry',
    USER_INFO: 'user_info',
    LAST_REFRESH: 'last_refresh'
} as const;

// ============================================================================
// Token Payload Interface
// ============================================================================

interface TokenPayload {
    exp: number;
    iat: number;
    user_id: number;
    email: string;
    role: string;
    [key: string]: any;
}

interface StoredUserInfo {
    id: number;
    email: string;
    role: string;
    name?: string;
}

// ============================================================================
// Enhanced Token Manager Class
// ============================================================================

export class EnhancedTokenManager {
    private config: StorageConfig;
    private storage: Storage;

    constructor(config: Partial<StorageConfig> = {}) {
        this.config = { ...DEFAULT_CONFIG, ...config };
        this.storage = this.config.useSessionStorage ? sessionStorage : localStorage;
    }

    // ============================================================================
    // Core Token Storage Methods
    // ============================================================================

    /**
     * Store access token securely
     */
    setAccessToken(token: string): void {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.ACCESS_TOKEN);
            const value = this.config.encryptTokens ? this.encryptValue(token) : token;
            this.storage.setItem(key, value);

            // Store token expiry for quick access
            const payload = this.decodeTokenPayload(token);
            if (payload) {
                this.setTokenExpiry(new Date(payload.exp * 1000));
            }
        } catch (error) {
            console.error('Failed to store access token:', error);
            throw new ApiError(
                ApiErrorCode.INTERNAL_ERROR,
                'Failed to store authentication token',
                500
            );
        }
    }

    /**
     * Retrieve access token
     */
    getAccessToken(): string | null {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.ACCESS_TOKEN);
            const value = this.storage.getItem(key);

            if (!value) {
                return null;
            }

            return this.config.encryptTokens ? this.decryptValue(value) : value;
        } catch (error) {
            console.error('Failed to retrieve access token:', error);
            return null;
        }
    }

    /**
     * Store refresh token securely
     */
    setRefreshToken(token: string): void {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.REFRESH_TOKEN);
            const value = this.config.encryptTokens ? this.encryptValue(token) : token;
            this.storage.setItem(key, value);
        } catch (error) {
            console.error('Failed to store refresh token:', error);
            throw new ApiError(
                ApiErrorCode.INTERNAL_ERROR,
                'Failed to store refresh token',
                500
            );
        }
    }

    /**
     * Retrieve refresh token
     */
    getRefreshToken(): string | null {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.REFRESH_TOKEN);
            const value = this.storage.getItem(key);

            if (!value) {
                return null;
            }

            return this.config.encryptTokens ? this.decryptValue(value) : value;
        } catch (error) {
            console.error('Failed to retrieve refresh token:', error);
            return null;
        }
    }

    /**
     * Store both tokens at once
     */
    setTokens(accessToken: string, refreshToken: string): void {
        this.setAccessToken(accessToken);
        this.setRefreshToken(refreshToken);
        this.setLastRefreshTime(new Date());

        // Store user info from access token
        const userInfo = this.extractUserInfoFromToken(accessToken);
        if (userInfo) {
            this.setUserInfo(userInfo);
        }
    }

    // ============================================================================
    // Token Validation Methods
    // ============================================================================

    /**
     * Check if valid tokens exist
     */
    hasValidTokens(): boolean {
        const accessToken = this.getAccessToken();
        const refreshToken = this.getRefreshToken();

        if (!accessToken || !refreshToken) {
            return false;
        }

        // Check if refresh token is still valid
        return !this.isTokenExpired(refreshToken);
    }

    /**
     * Check if access token is expired
     */
    isAccessTokenExpired(): boolean {
        const token = this.getAccessToken();
        return !token || this.isTokenExpired(token);
    }

    /**
     * Check if refresh token is expired
     */
    isRefreshTokenExpired(): boolean {
        const token = this.getRefreshToken();
        return !token || this.isTokenExpired(token);
    }

    /**
     * Check if access token needs refresh (expires within threshold)
     */
    needsTokenRefresh(thresholdMinutes: number = 5): boolean {
        const token = this.getAccessToken();

        if (!token) {
            return false;
        }

        return this.isTokenExpiringSoon(token, thresholdMinutes * 60);
    }

    /**
     * Get time until token expires
     */
    getTimeUntilExpiry(): number | null {
        const expiryDate = this.getTokenExpiryDate();

        if (!expiryDate) {
            return null;
        }

        return expiryDate.getTime() - Date.now();
    }

    // ============================================================================
    // Token Cleanup Methods
    // ============================================================================

    /**
     * Clear all authentication data
     */
    clearTokens(): void {
        try {
            const keys = Object.values(STORAGE_KEYS);
            keys.forEach(key => {
                const storageKey = this.getStorageKey(key);
                this.storage.removeItem(storageKey);
            });
        } catch (error) {
            console.error('Failed to clear tokens:', error);
        }
    }

    /**
     * Clear only access token (keep refresh token)
     */
    clearAccessToken(): void {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.ACCESS_TOKEN);
            this.storage.removeItem(key);
            this.storage.removeItem(this.getStorageKey(STORAGE_KEYS.TOKEN_EXPIRY));
        } catch (error) {
            console.error('Failed to clear access token:', error);
        }
    }

    /**
     * Clear expired tokens automatically
     */
    cleanupExpiredTokens(): void {
        if (this.isRefreshTokenExpired()) {
            this.clearTokens();
        } else if (this.isAccessTokenExpired()) {
            this.clearAccessToken();
        }
    }

    // ============================================================================
    // User Information Methods
    // ============================================================================

    /**
     * Store user information
     */
    setUserInfo(userInfo: StoredUserInfo): void {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.USER_INFO);
            const value = JSON.stringify(userInfo);
            this.storage.setItem(key, value);
        } catch (error) {
            console.error('Failed to store user info:', error);
        }
    }

    /**
     * Retrieve stored user information
     */
    getUserInfo(): StoredUserInfo | null {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.USER_INFO);
            const value = this.storage.getItem(key);

            if (!value) {
                return null;
            }

            return JSON.parse(value) as StoredUserInfo;
        } catch (error) {
            console.error('Failed to retrieve user info:', error);
            return null;
        }
    }

    // ============================================================================
    // Automatic Refresh Detection
    // ============================================================================

    /**
     * Set last refresh time
     */
    private setLastRefreshTime(date: Date): void {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.LAST_REFRESH);
            this.storage.setItem(key, date.toISOString());
        } catch (error) {
            console.error('Failed to store last refresh time:', error);
        }
    }

    /**
     * Get last refresh time
     */
    getLastRefreshTime(): Date | null {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.LAST_REFRESH);
            const value = this.storage.getItem(key);

            if (!value) {
                return null;
            }

            return new Date(value);
        } catch (error) {
            console.error('Failed to retrieve last refresh time:', error);
            return null;
        }
    }

    /**
     * Check if token was recently refreshed (to avoid rapid refresh attempts)
     */
    wasRecentlyRefreshed(thresholdSeconds: number = 30): boolean {
        const lastRefresh = this.getLastRefreshTime();

        if (!lastRefresh) {
            return false;
        }

        const timeSinceRefresh = Date.now() - lastRefresh.getTime();
        return timeSinceRefresh < (thresholdSeconds * 1000);
    }

    // ============================================================================
    // Private Utility Methods
    // ============================================================================

    private getStorageKey(key: string): string {
        return `${this.config.tokenPrefix}${key}`;
    }

    private setTokenExpiry(date: Date): void {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.TOKEN_EXPIRY);
            this.storage.setItem(key, date.toISOString());
        } catch (error) {
            console.error('Failed to store token expiry:', error);
        }
    }

    private getTokenExpiryDate(): Date | null {
        try {
            const key = this.getStorageKey(STORAGE_KEYS.TOKEN_EXPIRY);
            const value = this.storage.getItem(key);

            if (!value) {
                return null;
            }

            return new Date(value);
        } catch (error) {
            console.error('Failed to retrieve token expiry:', error);
            return null;
        }
    }

    private decodeTokenPayload(token: string): TokenPayload | null {
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

    private isTokenExpired(token: string): boolean {
        const payload = this.decodeTokenPayload(token);
        if (!payload) {
            return true;
        }

        const currentTime = Math.floor(Date.now() / 1000);
        return payload.exp < currentTime;
    }

    private isTokenExpiringSoon(token: string, thresholdSeconds: number): boolean {
        const payload = this.decodeTokenPayload(token);
        if (!payload) {
            return true;
        }

        const currentTime = Math.floor(Date.now() / 1000);
        return payload.exp < (currentTime + thresholdSeconds);
    }

    private extractUserInfoFromToken(token: string): StoredUserInfo | null {
        const payload = this.decodeTokenPayload(token);
        if (!payload) {
            return null;
        }

        return {
            id: payload.user_id,
            email: payload.email,
            role: payload.role,
            name: payload.name
        };
    }

    // Basic encryption/decryption (for demonstration - use proper encryption in production)
    private encryptValue(value: string): string {
        // This is a basic implementation - use proper encryption in production
        return btoa(value);
    }

    private decryptValue(value: string): string {
        // This is a basic implementation - use proper encryption in production
        try {
            return atob(value);
        } catch (error) {
            throw new Error('Failed to decrypt token');
        }
    }
}

// ============================================================================
// Legacy Token Manager (for backward compatibility)
// ============================================================================

export class TokenManager {
    private static readonly ACCESS_TOKEN_KEY = 'access_token';
    private static readonly REFRESH_TOKEN_KEY = 'refresh_token';

    static getAccessToken(): string | null {
        return localStorage.getItem(this.ACCESS_TOKEN_KEY);
    }

    static setAccessToken(token: string): void {
        localStorage.setItem(this.ACCESS_TOKEN_KEY, token);
    }

    static getRefreshToken(): string | null {
        return localStorage.getItem(this.REFRESH_TOKEN_KEY);
    }

    static setRefreshToken(token: string): void {
        localStorage.setItem(this.REFRESH_TOKEN_KEY, token);
    }

    static setTokens(accessToken: string, refreshToken: string): void {
        this.setAccessToken(accessToken);
        this.setRefreshToken(refreshToken);
    }

    static clearTokens(): void {
        localStorage.removeItem(this.ACCESS_TOKEN_KEY);
        localStorage.removeItem(this.REFRESH_TOKEN_KEY);
    }

    static hasValidTokens(): boolean {
        return !!(this.getAccessToken() && this.getRefreshToken());
    }
}

// ============================================================================
// Token Refresh Manager
// ============================================================================

export class TokenRefreshManager {
    private static refreshPromise: Promise<boolean> | null = null;
    private static isRefreshing = false;

    /**
     * Ensure only one refresh operation happens at a time
     */
    static async ensureSingleRefresh(refreshFunction: () => Promise<boolean>): Promise<boolean> {
        if (this.isRefreshing && this.refreshPromise) {
            return this.refreshPromise;
        }

        this.isRefreshing = true;
        this.refreshPromise = refreshFunction().finally(() => {
            this.isRefreshing = false;
            this.refreshPromise = null;
        });

        return this.refreshPromise;
    }

    /**
     * Check if currently refreshing
     */
    static isCurrentlyRefreshing(): boolean {
        return this.isRefreshing;
    }
}

// ============================================================================
// Export default enhanced token manager instance
// ============================================================================

export const tokenManager = new EnhancedTokenManager();
export default tokenManager;