// API Error Handling System
// This file provides comprehensive error handling utilities for API responses

import type { AxiosError, AxiosResponse } from 'axios';
import { AxiosError as AxiosErrorClass } from 'axios';
import type { ApiResponse, ErrorResponse, ApiErrorDetails } from './types';
import { ApiErrorCode } from './types';

// ============================================================================
// Custom API Error Class
// ============================================================================

export class ApiError extends Error {
    public readonly code: ApiErrorCode;
    public readonly statusCode: number;
    public readonly details?: ApiErrorDetails;
    public readonly requestId?: string;
    public readonly timestamp: Date;
    public readonly isNetworkError: boolean;
    public readonly isAuthenticationError: boolean;
    public readonly isValidationError: boolean;
    public readonly isServerError: boolean;

    constructor(
        code: ApiErrorCode,
        message: string,
        statusCode: number = 0,
        details?: ApiErrorDetails,
        requestId?: string
    ) {
        super(message);
        this.name = 'ApiError';
        this.code = code;
        this.statusCode = statusCode;
        this.details = details;
        this.requestId = requestId;
        this.timestamp = new Date();

        // Categorize error types
        this.isNetworkError = this.isNetworkErrorType(code);
        this.isAuthenticationError = this.isAuthenticationErrorType(code);
        this.isValidationError = this.isValidationErrorType(code);
        this.isServerError = this.isServerErrorType(code);

        // Maintain proper stack trace (Node.js specific)
        if (typeof (Error as any).captureStackTrace === 'function') {
            (Error as any).captureStackTrace(this, ApiError);
        }
    }

    private isNetworkErrorType(code: ApiErrorCode): boolean {
        return ([
            ApiErrorCode.NETWORK_ERROR,
            ApiErrorCode.TIMEOUT,
            ApiErrorCode.REQUEST_CANCELLED
        ] as ApiErrorCode[]).includes(code);
    }

    private isAuthenticationErrorType(code: ApiErrorCode): boolean {
        return ([
            ApiErrorCode.UNAUTHORIZED,
            ApiErrorCode.FORBIDDEN,
            ApiErrorCode.TOKEN_EXPIRED
        ] as ApiErrorCode[]).includes(code);
    }

    private isValidationErrorType(code: ApiErrorCode): boolean {
        return ([
            ApiErrorCode.VALIDATION_ERROR,
            ApiErrorCode.INVALID_INPUT
        ] as ApiErrorCode[]).includes(code);
    }

    private isServerErrorType(code: ApiErrorCode): boolean {
        return ([
            ApiErrorCode.INTERNAL_ERROR,
            ApiErrorCode.SERVICE_UNAVAILABLE
        ] as ApiErrorCode[]).includes(code);
    }

    // ============================================================================
    // Static Factory Methods
    // ============================================================================

    static fromAxiosError(error: AxiosError): ApiError {
        // Handle network errors
        if (!error.response) {
            if (error.code === 'ECONNABORTED' || error.message.includes('timeout')) {
                return new ApiError(
                    ApiErrorCode.TIMEOUT,
                    'Request timeout. Please try again.',
                    408
                );
            }

            if (error.code === 'ERR_CANCELED') {
                return new ApiError(
                    ApiErrorCode.REQUEST_CANCELLED,
                    'Request was cancelled.',
                    0
                );
            }

            return new ApiError(
                ApiErrorCode.NETWORK_ERROR,
                'Network connection failed. Please check your internet connection.',
                0
            );
        }

        // Handle server response errors
        const response = error.response;
        const statusCode = response.status;

        // Try to extract error from API response format
        if (response.data && typeof response.data === 'object') {
            const apiResponse = response.data as ApiResponse;

            if (apiResponse.error) {
                return ApiError.fromErrorResponse(apiResponse.error, statusCode);
            }
        }

        // Fallback to HTTP status code mapping
        return ApiError.fromHttpStatus(statusCode, response.statusText);
    }

    static fromErrorResponse(errorResponse: ErrorResponse, statusCode: number): ApiError {
        const code = ApiError.mapErrorCodeToApiErrorCode(errorResponse.code);

        return new ApiError(
            code,
            errorResponse.message,
            statusCode,
            errorResponse.details,
            errorResponse.request_id
        );
    }

    static fromHttpStatus(statusCode: number, statusText: string): ApiError {
        switch (statusCode) {
            case 400:
                return new ApiError(
                    ApiErrorCode.VALIDATION_ERROR,
                    'Invalid request data.',
                    statusCode
                );
            case 401:
                return new ApiError(
                    ApiErrorCode.UNAUTHORIZED,
                    'Authentication required. Please log in.',
                    statusCode
                );
            case 403:
                return new ApiError(
                    ApiErrorCode.FORBIDDEN,
                    'Access denied. You do not have permission to perform this action.',
                    statusCode
                );
            case 404:
                return new ApiError(
                    ApiErrorCode.NOT_FOUND,
                    'The requested resource was not found.',
                    statusCode
                );
            case 409:
                return new ApiError(
                    ApiErrorCode.CONFLICT,
                    'Conflict with existing data.',
                    statusCode
                );
            case 429:
                return new ApiError(
                    ApiErrorCode.RATE_LIMIT_EXCEEDED,
                    'Too many requests. Please try again later.',
                    statusCode
                );
            case 500:
                return new ApiError(
                    ApiErrorCode.INTERNAL_ERROR,
                    'Internal server error. Please try again later.',
                    statusCode
                );
            case 503:
                return new ApiError(
                    ApiErrorCode.SERVICE_UNAVAILABLE,
                    'Service temporarily unavailable. Please try again later.',
                    statusCode
                );
            default:
                return new ApiError(
                    ApiErrorCode.UNKNOWN,
                    statusText || 'An unexpected error occurred.',
                    statusCode
                );
        }
    }

    private static mapErrorCodeToApiErrorCode(serverCode: string): ApiErrorCode {
        const codeMap: Record<string, ApiErrorCode> = {
            'VALIDATION_ERROR': ApiErrorCode.VALIDATION_ERROR,
            'INVALID_INPUT': ApiErrorCode.INVALID_INPUT,
            'UNAUTHORIZED': ApiErrorCode.UNAUTHORIZED,
            'FORBIDDEN': ApiErrorCode.FORBIDDEN,
            'TOKEN_EXPIRED': ApiErrorCode.TOKEN_EXPIRED,
            'NOT_FOUND': ApiErrorCode.NOT_FOUND,
            'CONFLICT': ApiErrorCode.CONFLICT,
            'INTERNAL_ERROR': ApiErrorCode.INTERNAL_ERROR,
            'SERVICE_UNAVAILABLE': ApiErrorCode.SERVICE_UNAVAILABLE,
            'RATE_LIMIT_EXCEEDED': ApiErrorCode.RATE_LIMIT_EXCEEDED,
            'REQUEST_CANCELLED': ApiErrorCode.REQUEST_CANCELLED,
            'TIMEOUT_ERROR': ApiErrorCode.TIMEOUT,
        };

        return codeMap[serverCode] || ApiErrorCode.UNKNOWN;
    }
}

// ============================================================================
// Error Processing Utilities
// ============================================================================

export class ErrorProcessor {
    /**
     * Process any error and convert it to ApiError
     */
    static processError(error: unknown): ApiError {
        if (error instanceof ApiError) {
            return error;
        }

        if (error instanceof AxiosErrorClass) {
            return ApiError.fromAxiosError(error);
        }

        if (error instanceof Error) {
            return new ApiError(
                ApiErrorCode.UNKNOWN,
                error.message,
                0
            );
        }

        return new ApiError(
            ApiErrorCode.UNKNOWN,
            'An unexpected error occurred.',
            0
        );
    }

    /**
     * Generate user-friendly error messages
     */
    static getUserFriendlyMessage(error: ApiError): string {
        // Return custom user-friendly messages based on error type
        switch (error.code) {
            case ApiErrorCode.NETWORK_ERROR:
                return 'Unable to connect to the server. Please check your internet connection and try again.';

            case ApiErrorCode.TIMEOUT:
                return 'The request took too long to complete. Please try again.';

            case ApiErrorCode.UNAUTHORIZED:
                return 'You need to log in to access this feature.';

            case ApiErrorCode.FORBIDDEN:
                return 'You do not have permission to perform this action.';

            case ApiErrorCode.NOT_FOUND:
                return 'The requested item could not be found.';

            case ApiErrorCode.VALIDATION_ERROR:
                return ErrorProcessor.formatValidationMessage(error);

            case ApiErrorCode.RATE_LIMIT_EXCEEDED:
                return 'You are making requests too quickly. Please wait a moment and try again.';

            case ApiErrorCode.INTERNAL_ERROR:
                return 'Something went wrong on our end. Please try again later.';

            case ApiErrorCode.SERVICE_UNAVAILABLE:
                return 'The service is temporarily unavailable. Please try again in a few minutes.';

            default:
                return error.message || 'An unexpected error occurred. Please try again.';
        }
    }

    /**
     * Format validation error messages
     */
    private static formatValidationMessage(error: ApiError): string {
        if (!error.details || typeof error.details !== 'object') {
            return 'Please check your input and try again.';
        }

        const fieldErrors: string[] = [];

        Object.entries(error.details).forEach(([field, details]) => {
            if (typeof details === 'string') {
                fieldErrors.push(`${field}: ${details}`);
            } else if (details && typeof details === 'object' && details.message) {
                fieldErrors.push(`${field}: ${details.message}`);
            }
        });

        if (fieldErrors.length > 0) {
            return `Please fix the following errors:\n${fieldErrors.join('\n')}`;
        }

        return 'Please check your input and try again.';
    }

    /**
     * Check if error should trigger a retry
     */
    static shouldRetry(error: ApiError, retryCount: number = 0): boolean {
        const maxRetries = 3;

        if (retryCount >= maxRetries) {
            return false;
        }

        // Retry on network errors and server errors
        return error.isNetworkError ||
            error.code === ApiErrorCode.INTERNAL_ERROR ||
            error.code === ApiErrorCode.SERVICE_UNAVAILABLE;
    }

    /**
     * Get retry delay in milliseconds (exponential backoff)
     */
    static getRetryDelay(retryCount: number): number {
        return Math.min(1000 * Math.pow(2, retryCount), 10000); // Max 10 seconds
    }
}

// ============================================================================
// Response Handler Utilities
// ============================================================================

export class ResponseHandler {
    /**
     * Handle successful API response
     */
    static handleSuccess<T>(response: AxiosResponse<ApiResponse<T>>): T {
        const apiResponse = response.data;

        if (apiResponse.success && apiResponse.data !== undefined) {
            return apiResponse.data;
        }

        // If response indicates failure, throw error
        if (apiResponse.error) {
            throw ApiError.fromErrorResponse(apiResponse.error, response.status);
        }

        // Fallback error
        throw new ApiError(
            ApiErrorCode.UNKNOWN,
            'Invalid response format',
            response.status
        );
    }

    /**
     * Handle paginated API response
     */
    static handlePaginatedSuccess<T>(response: AxiosResponse<ApiResponse<T>>) {
        const apiResponse = response.data;

        if (apiResponse.success && apiResponse.data !== undefined) {
            return {
                data: apiResponse.data,
                pagination: apiResponse.pagination
            };
        }

        if (apiResponse.error) {
            throw ApiError.fromErrorResponse(apiResponse.error, response.status);
        }

        throw new ApiError(
            ApiErrorCode.UNKNOWN,
            'Invalid response format',
            response.status
        );
    }

    /**
     * Handle API error response
     */
    static handleError(error: unknown): never {
        const apiError = ErrorProcessor.processError(error);
        throw apiError;
    }
}

// ============================================================================
// Network Detection Utilities
// ============================================================================

export class NetworkDetector {
    /**
     * Check if the error is due to network connectivity issues
     */
    static isNetworkError(error: ApiError): boolean {
        return error.isNetworkError;
    }

    /**
     * Check if browser is online
     */
    static isOnline(): boolean {
        return typeof navigator !== 'undefined' ? navigator.onLine : true;
    }

    /**
     * Add network status listeners
     */
    static addNetworkListeners(
        onOnline: () => void,
        onOffline: () => void
    ): () => void {
        if (typeof window === 'undefined') {
            return () => { }; // No-op for SSR
        }

        window.addEventListener('online', onOnline);
        window.addEventListener('offline', onOffline);

        // Return cleanup function
        return () => {
            window.removeEventListener('online', onOnline);
            window.removeEventListener('offline', onOffline);
        };
    }
}

// ============================================================================
// Exports
// ============================================================================

export {
    ApiErrorCode,
    type ApiErrorDetails,
    type ErrorResponse
} from './types';

export default {
    ApiError,
    ErrorProcessor,
    ResponseHandler,
    NetworkDetector
};