import { ApiErrorCode } from "./error-codes";

/**
 * Custom error class for API-related errors.
 * Contains a message for display and a code for programmatic handling.
 */
export class ApiError extends Error {
    constructor(public message: string, public code?: ApiErrorCode, public data?: any) {
        super(message);
        this.name = "ApiError";
    }
}
