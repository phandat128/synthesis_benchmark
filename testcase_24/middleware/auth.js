const jwt = require('jsonwebtoken');
const User = require('../models/User');

// Custom error class for API errors
class ApiError extends Error {
    constructor(message, statusCode) {
        super(message);
        this.statusCode = statusCode;
    }
}

/**
 * Protect middleware: Verifies JWT token and attaches user to req.
 */
exports.protect = async (req, res, next) => {
    let token;

    // 1. Check for token in Authorization header (Bearer scheme)
    if (req.headers.authorization && req.headers.authorization.startsWith('Bearer')) {
        token = req.headers.authorization.split(' ')[1];
    }

    // 2. Check if token exists
    if (!token) {
        return next(new ApiError('Authentication required. No token provided.', 401));
    }

    try {
        // 3. Verify token using the secret key
        const decoded = jwt.verify(token, process.env.JWT_SECRET);

        // 4. Find user by ID, excluding the sensitive password hash
        const user = await User.findById(decoded.id).select('-password');

        if (!user) {
            return next(new ApiError('User associated with this token no longer exists.', 401));
        }

        // 5. Attach user object to the request for downstream controllers
        req.user = user;
        next();
    } catch (err) {
        // Handle token expiration, invalid signature, or other JWT errors
        console.error('JWT Verification Error:', err.message);
        return next(new ApiError('Not authorized. Invalid or expired token.', 401));
    }
};

/**
 * Authorize middleware: Checks if the user role is included in the required roles array.
 */
exports.authorize = (...roles) => {
    return (req, res, next) => {
        if (!req.user || !roles.includes(req.user.role)) {
            return next(new ApiError(`User role ${req.user.role} is not authorized to access this resource.`, 403));
        }
        next();
    };
};