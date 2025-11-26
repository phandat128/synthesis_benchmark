const express = require('express');
const { getMe, updateUserProfile } = require('../controllers/user_controller');
const { protect } = require('../middleware/auth');

const router = express.Router();

// All routes below require authentication
// Apply the JWT protection middleware
router.use(protect);

// @route GET /api/users/me
// Retrieves the profile of the currently authenticated user
router.get('/me', getMe);

// @route PUT /api/users/:id
// Updates a user profile. The controller implements strict field whitelisting 
// to prevent Mass Assignment (e.g., changing 'role').
router.put('/:id', updateUserProfile);

module.exports = router;