// API Client Configuration
// This file sets up the base Axios instance with interceptors for authentication and error handling

import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import type { ApiResponse, LoginResponse } from './types';

// ============================================================================
// Configuration Constants
// ============================================================================

const DEFAULT_TIMEOUT = 30000; // 30 seconds
const DEFAULT_BASE_URL = '/api/v1';

// Import enhanced token manager
import { tokenManager, TokenManager } from './tokenManager';

// ============================================================================
// API Client Configuration
// ============================================================================

class ApiClient {
    private axiosInstance: AxiosInstance;
    private isRefreshing = false;
    private failedQueue: Array<{
        resolve: (value?: any) => void;
        reject: (error?: any) => void;
    }> = [];

    constructor() {
        this.axiosInstance = this.createAxiosInstance();
        this.setupInterceptors();
    }

    private createAxiosInstance(): AxiosInstance {
        const baseURL = import.meta.env.VITE_API_BASE_URL || DEFAULT_BASE_URL;

        return axios.create({
            baseURL,
            timeout: DEFAULT_TIMEOUT,
            headers: {
                'Content-Type': 'application/json',
            },
            withCredentials: true, // Include cookies for refresh token
        });
    }

    private setupInterceptors(): void {
        this.setupRequestInterceptor();
        this.setupResponseInterceptor();
    }

    private setupRequestInterceptor(): void {
        this.axiosInstance.interceptors.request.use(
            (config) => {
                // Clean up expired tokens before making request
                tokenManager.cleanupExpiredTokens();

                const token = tokenManager.getAccessToken();
                if (token) {
                    config.headers.Authorization = `Bearer ${token}`;
                }
                return config;
            },
            (error) => {
                return Promise.reject(error);
            }
        );
    }

    private setupResponseInterceptor(): void {
        this.axiosInstance.interceptors.response.use(
            (response) => response,
            async (error: AxiosError) => {
                const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };

                // Handle 401 errors with token refresh
                if (error.response?.status === 401 && !originalRequest._retry) {
                    if (this.isRefreshing) {
                        // If already refreshing, queue the request
                        return new Promise((resolve, reject) => {
                            this.failedQueue.push({ resolve, reject });
                        }).then(() => {
                            return this.axiosInstance(originalRequest);
                        }).catch(err => {
                            return Promise.reject(err);
                        });
                    }

                    originalRequest._retry = true;
                    this.isRefreshing = true;

                    try {
                        const refreshToken = tokenManager.getRefreshToken();
                        if (!refreshToken) {
                            throw new Error('No refresh token available');
                        }

                        // Check if we recently refreshed to avoid rapid refresh attempts
                        if (tokenManager.wasRecentlyRefreshed()) {
                            throw new Error('Token was recently refreshed, avoiding rapid refresh');
                        }

                        // Attempt to refresh the token
                        const response = await this.refreshAccessToken();
                        const { access_token, refresh_token } = response;

                        // Update stored tokens using enhanced manager
                        tokenManager.setTokens(access_token, refresh_token);

                        // Process failed queue
                        this.processQueue(null);

                        // Retry original request
                        return this.axiosInstance(originalRequest);
                    } catch (refreshError) {
                        // Refresh failed, clear tokens and redirect to login
                        this.processQueue(refreshError);
                        tokenManager.clearTokens();
                        this.handleAuthenticationFailure();
                        return Promise.reject(refreshError);
                    } finally {
                        this.isRefreshing = false;
                    }
                }

                return Promise.reject(error);
            }
        );
    }

    private async refreshAccessToken(): Promise<LoginResponse> {
        // Create a separate axios instance for refresh to avoid interceptor loops
        const refreshClient = axios.create({
            baseURL: this.axiosInstance.defaults.baseURL,
            timeout: DEFAULT_TIMEOUT,
            withCredentials: true,
        });

        const response = await refreshClient.post<ApiResponse<LoginResponse>>('/auth/refresh');

        if (response.data.success && response.data.data) {
            return response.data.data;
        }

        throw new Error('Token refresh failed');
    }

    private processQueue(error: any): void {
        this.failedQueue.forEach(({ resolve, reject }) => {
            if (error) {
                reject(error);
            } else {
                resolve();
            }
        });

        this.failedQueue = [];
    }

    private handleAuthenticationFailure(): void {
        // Redirect to login page or emit authentication failure event
        // This can be customized based on your routing setup
        if (typeof window !== 'undefined') {
            window.location.href = '/login';
        }
    }

    // ============================================================================
    // Public API Methods
    // ============================================================================

    getInstance(): AxiosInstance {
        return this.axiosInstance;
    }

    async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
        return this.axiosInstance.get<T>(url, config);
    }

    async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
        return this.axiosInstance.post<T>(url, data, config);
    }

    async put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
        return this.axiosInstance.put<T>(url, data, config);
    }

    async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
        return this.axiosInstance.delete<T>(url, config);
    }

    async patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
        return this.axiosInstance.patch<T>(url, data, config);
    }

    // ============================================================================
    // Configuration Methods
    // ============================================================================

    setBaseURL(baseURL: string): void {
        this.axiosInstance.defaults.baseURL = baseURL;
    }

    setTimeout(timeout: number): void {
        this.axiosInstance.defaults.timeout = timeout;
    }

    setDefaultHeaders(headers: Record<string, string>): void {
        Object.assign(this.axiosInstance.defaults.headers, headers);
    }

    // ============================================================================
    // Authentication Methods
    // ============================================================================

    setAuthToken(token: string): void {
        tokenManager.setAccessToken(token);
        this.axiosInstance.defaults.headers.Authorization = `Bearer ${token}`;
    }

    clearAuthToken(): void {
        tokenManager.clearTokens();
        delete this.axiosInstance.defaults.headers.Authorization;
    }

    isAuthenticated(): boolean {
        return tokenManager.hasValidTokens();
    }
}

// ============================================================================
// Export Singleton Instance
// ============================================================================

export const apiClient = new ApiClient();
export { TokenManager };
export default apiClient;