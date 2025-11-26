// utils/errorHandler.js
/**
 * Custom Error Handler Middleware
 * Prevents leaking sensitive stack trace information in production.
 */
const errorHandler = (err, req, res, next) => {
    // Default status code and message
    let statusCode = err.statusCode || 500;
    let message = err.message || 'Internal Server Error';

    // Handle Mongoose specific errors (e.g., CastError, ValidationError)
    if (err.name === 'CastError' && err.kind === 'ObjectId') {
        statusCode = 400; // Bad Request for invalid ID format
        message = 'Invalid resource ID format.';
    }

    if (err.name === 'ValidationError') {
        statusCode = 400;
        // Format validation errors nicely
        message = Object.values(err.errors).map(val => val.message).join(', ');
    }

    if (err.code === 11000) { // Duplicate key error (e.g., unique index violation)
        statusCode = 400;
        message = 'Duplicate field value entered.';
    }

    // Log the error for internal debugging (excluding client-side errors)
    if (statusCode >= 500 || process.env.NODE_ENV === 'development') {
        console.error(err.stack);
    }

    res.status(statusCode).json({
        success: false,
        error: message,
        // Only include stack trace in development environment to prevent information leakage
        stack: process.env.NODE_ENV === 'development' ? err.stack : undefined,
    });
};

module.exports = errorHandler;