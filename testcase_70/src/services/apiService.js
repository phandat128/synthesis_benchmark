import axios from 'axios';

const API_BASE_URL = '/api/v1';

// Configure axios instance
const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
        // In a real app, Authorization header (e.g., JWT) would be set here
    },
});

/**
 * Fetches the current user profile data.
 * @param {string} userId - The ID of the user to fetch.
 * @returns {Promise<object>} The user profile data.
 */
export const fetchUserProfile = async (userId) => {
    console.log(`[API] Fetching profile for user: ${userId}`);
    try {
        // Mock API response delay
        await new Promise(resolve => setTimeout(resolve, 500));

        // Mock successful response
        return {
            id: userId,
            name: 'John Doe',
            email: 'john.doe@secureapp.com',
            bio: 'Senior Developer focused on secure coding practices.',
            role: 'standard_user', // Sensitive field, should not be updatable by user input
        };

        // In a real scenario:
        // const response = await api.get(`/user/${userId}/profile`);
        // return response.data;
    } catch (error) {
        console.error("Error fetching user profile:", error);
        // Return a generic, non-sensitive error message
        throw new Error("Failed to load profile data. Please check your connection.");
    }
};

/**
 * Updates the user profile information.
 *
 * SECURITY NOTE: The payload passed here is assumed to be sanitized/whitelisted
 * by the caller (useProfileUpdate hook) to prevent Mass Assignment attacks.
 * The backend MUST also perform server-side validation and whitelisting.
 *
 * @param {string} userId - The ID of the user being updated.
 * @param {object} payload - The sanitized profile data to update (only whitelisted fields).
 * @returns {Promise<object>} The updated profile data.
 */
export const updateUserProfile = async (userId, payload) => {
    console.log(`[API] Attempting to update profile for user: ${userId}`);
    console.log("[API] Payload received (should be sanitized):", payload);

    // Simulate API call
    try {
        // const response = await api.put(`/user/${userId}/profile`, payload);
        // return response.data;

        await new Promise(resolve => setTimeout(resolve, 800));

        // Simulate successful update
        // Crucially, the mock backend blindly accepts and returns all fields in the payload.
        return {
            ...payload,
            id: userId,
            message: "Profile updated successfully.",
        };

    } catch (error) {
        console.error("Error updating user profile:", error);
        // Implement robust error handling that does not leak sensitive information
        throw new Error("Failed to update profile. The server reported an issue.");
    }
};