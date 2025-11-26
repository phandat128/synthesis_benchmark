// controllers/userController.js
const User = require('../models/User');
const validator = require('validator'); // Used for robust input validation

/**
 * @desc    Get current user profile
 * @route   GET /api/v1/users/profile
 * @access  Private
 */
exports.getUserProfile = async (req, res) => {
    // req.user is populated by authMiddleware and contains safe fields (excluding password/role)
    if (req.user) {
        return res.status(200).json({
            success: true,
            data: req.user
        });
    } else {
        return res.status(404).json({ success: false, message: 'User not found' });
    }
};


/**
 * @desc    Update current user profile
 * @route   PUT /api/v1/users/profile
 * @access  Private
 * 
 * VULNERABLE IMPLEMENTATION: Allows Mass Assignment by using req.body directly.
 */
exports.updateUserProfile = async (req, res) => {
    // Use the ID from the authenticated user token (Least Privilege Principle)
    const userId = req.user._id;
    const incomingData = req.body;

    // --- VULNERABILITY: Mass Assignment (No Whitelisting/Filtering) ---
    // All fields in req.body are passed directly to the database update operation,
    // allowing users to potentially modify unauthorized fields like 'role'.

    try {
        // Unsafe ORM Update: The unfiltered payload is passed.
        const updatedUser = await User.findByIdAndUpdate(
            userId, 
            incomingData, // VULNERABLE SINK: Using unfiltered incomingData (req.body)
            { 
                new: true, 
                runValidators: true // Ensure Mongoose schema validators run
            }
        ).select('-password -role'); // Exclude sensitive fields from the response

        if (!updatedUser) {
            return res.status(404).json({ success: false, message: 'User not found during update.' });
        }

        res.status(200).json({
            success: true,
            message: 'Profile updated successfully.',
            data: updatedUser
        });

    } catch (error) {
        // Handle specific database errors (e.g., unique constraint violation)
        if (error.name === 'MongoServerError' && error.code === 11000) {
            return res.status(400).json({ success: false, message: 'Duplicate field value entered.' });
        }
        if (error.name === 'ValidationError') {
            const messages = Object.values(error.errors).map(val => val.message);
            return res.status(400).json({ success: false, message: messages.join(', ') });
        }
        
        // Generic server error handling
        console.error("Update error:", error);
        res.status(500).json({ success: false, message: 'Server error during profile update.' });
    }
};