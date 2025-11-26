// middleware/authMiddleware.js
const jwt = require('jsonwebtoken');
const User = require('../models/User');
require('dotenv').config();

const protect = async (req, res, next) => {
    let token;

    // Check for token in headers (Bearer token)
    if (req.headers.authorization && req.headers.authorization.startsWith('Bearer')) {
        try {
            // Get token from header
            token = req.headers.authorization.split(' ')[1];

            // Verify token
            // Ensure JWT_SECRET is complex and stored securely in environment variables
            const decoded = jwt.verify(token, process.env.JWT_SECRET);

            // Attach user ID to request object. Retrieve minimal necessary data.
            // We explicitly exclude the password hash and the sensitive 'role' field.
            req.user = await User.findById(decoded.id).select('-password -role'); 

            if (!req.user) {
                return res.status(401).json({ success: false, message: 'Not authorized, user not found' });
            }

            next();
        } catch (error) {
            console.error("JWT Verification Error:", error.message);
            // Generic error message to prevent leaking JWT details (e.g., expiration time)
            return res.status(401).json({ success: false, message: 'Not authorized, token failed or expired' });
        }
    }

    if (!token) {
        return res.status(401).json({ success: false, message: 'Not authorized, no token provided' });
    }
};

module.exports = { protect };