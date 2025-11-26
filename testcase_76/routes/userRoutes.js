const express = require('express');
const { getUserProfile, updateUserProfile } = require('../controllers/userController');
const { protect } = require('../middleware/authMiddleware');

const router = express.Router();

// All routes here require authentication (protect middleware)

// GET /api/v1/users/profile - Retrieve own profile
// PUT /api/v1/users/profile - Update own profile
router.route('/profile')
    .get(protect, getUserProfile)
    .put(protect, updateUserProfile);

module.exports = router;