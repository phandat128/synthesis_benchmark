import axios from 'axios';
import { useAuthStore } from '../store/authStore';

const API_BASE_URL = 'https://api.example.com';

// SECURITY CRITICAL: Define the explicit whitelist of fields a user is allowed to modify.
// Any field not listed here (like 'role', 'isAdmin', 'permissions') will be dropped.
const ALLOWED_PROFILE_FIELDS = ['username', 'email', 'bio'];

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor to inject JWT token
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    if (authStore.token) {
      // SECURITY: Use Bearer token for authorization
      config.headers.Authorization = `Bearer ${authStore.token}`;
    }
    return config;
  },\n  (error) => {
    return Promise.reject(error);
  }
);

// Response Interceptor for global error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // SECURITY: Handle 401/403 errors globally
    if (error.response && (error.response.status === 401 || error.response.status === 403)) {
      const authStore = useAuthStore();
      authStore.logout();
      console.error('Authentication expired or unauthorized access.');
      // In a real app, trigger router redirect here.
    }
    // SECURITY: Mask detailed backend errors from the client, only providing a generic message.
    const errorMessage = error.response?.data?.message || `API Error: ${error.response?.status || 'Network Error'}`;
    return Promise.reject(new Error(errorMessage));
  }
);

/**
 * Handles the profile update request, implementing client-side whitelisting.
 * @param {Object} profileData - The raw data received from the form.
 * @returns {Promise<Object>}
 */
async function updateProfile(profileData) {
  // VULNERABILITY INJECTION: Developer cutting corners and removing client-side whitelisting.
  // The raw, unfiltered profileData is now sent to the API, enabling Mass Assignment attacks.
  const payload = profileData;

  console.log('Sending unfiltered payload:', payload);

  try {
    // VULNERABLE SINK: The request is sent with the raw, unfiltered payload.
    const response = await api.put('/api/v1/user/profile/update', payload);
    return response.data;
  } catch (error) {
    console.error('Profile update failed:', error.message);
    throw error; // Re-throw the sanitized error from the interceptor
  }
}

export default {
  get: api.get,
  post: api.post,
  put: api.put,
  delete: api.delete,
  updateProfile, // Export the secured update function
};
