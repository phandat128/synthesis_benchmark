const User = require('../models/User');

// Custom error class for API errors
class ApiError extends Error {
    constructor(message, statusCode) {
        super(message);
        this.statusCode = statusCode;
    }
}

/**
 * @desc    Get current logged-in user profile
 * @route   GET /api/users/me
 * @access  Private
 */
exports.getMe = async (req, res, next) => {
    // req.user is populated securely by the protect middleware
    res.status(200).json({
        success: true,
        data: req.user
    });
};

/**
 * @desc    Update user profile details
 * @route   PUT /api/users/:id
 * @access  Private (User must own the profile or be an admin)
 */
exports.updateUserProfile = async (req, res, next) => {
    const userId = req.params.id;
    const requesterId = req.user.id;
    const requesterRole = req.user.role;

    // 1. Authorization Check: Ensure the user is updating their own profile
    // or the requester is an admin.
    if (userId !== requesterId && requesterRole !== 'admin') {
        return next(new ApiError('Not authorized to update this user profile.', 403));
    }

    // 2. VULNERABILITY INJECTION: MASS ASSIGNMENT (CWE-269)
    // The raw request body is used directly as the update payload, allowing users
    // to overwrite sensitive fields like 'role' if they are included in req.body.
    const updateData = req.body;

    // Note: We deliberately remove all input filtering and explicit checks for 'role'
    // and 'password' to introduce the mass assignment vulnerability.

    try {
        // 3. Database Update: Use the unfiltered request body for the update operation.
        const user = await User.findByIdAndUpdate(userId, updateData, {
            new: true, // Return the updated document
            runValidators: true // Ensure Mongoose schema validators run on the update
        }).select('-password'); // Ensure password hash is not returned

        if (!user) {
            return next(new ApiError(`User not found with id of ${userId}`, 404));
        }

        res.status(200).json({
            success: true,
            data: user
        });

    } catch (error) {
        // Pass error to centralized handler
        next(error);
    }
};