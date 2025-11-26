import { defineStore } from 'pinia';
import apiService from '../services/apiService';

// Define the structure of the user object, ensuring only expected fields are stored.
const DEFAULT_USER_PROFILE = {
    id: null,
    username: '',
    email: '',
    role: 'guest', // Important: Store the role, but ensure it cannot be updated via profile form.
    bio: ''
};

export const useAuthStore = defineStore('auth', {
  state: () => ({
    // SECURITY: Store JWT securely (e.g., in HttpOnly cookies, but for this client-side example, we use localStorage/state for demonstration).
    token: localStorage.getItem('authToken') || null,
    user: { ...DEFAULT_USER_PROFILE },
    isAuthenticated: !!localStorage.getItem('authToken'),
    isLoadingProfile: false,
  }),
  
  actions: {
    setToken(newToken) {
      this.token = newToken;
      this.isAuthenticated = !!newToken;
      if (newToken) {
        localStorage.setItem('authToken', newToken);
      } else {
        localStorage.removeItem('authToken');
      }
    },

    async fetchUserProfile() {
      if (!this.isAuthenticated || this.isLoadingProfile) return;

      this.isLoadingProfile = true;
      try {
        // Simulate API call to get user data
        const response = await apiService.get('/api/v1/user/profile');
        
        // SECURITY: Sanitize and validate incoming profile data against the expected structure.
        const fetchedData = response.data;
        this.user.id = fetchedData.id || null;
        this.user.username = fetchedData.username || '';
        this.user.email = fetchedData.email || '';
        this.user.role = fetchedData.role || 'user'; // Assuming backend correctly assigns role
        this.user.bio = fetchedData.bio || '';

        console.log('Profile fetched successfully.');
      } catch (error) {
        console.error('Failed to fetch user profile:', error);
        // SECURITY: Do not leak detailed error information to the user interface.
        this.logout(); 
      } finally {
        this.isLoadingProfile = false;
      }
    },

    logout() {
      this.setToken(null);
      this.user = { ...DEFAULT_USER_PROFILE };
      // Redirect to login handled by router guard or view component
    },

    // Action to update local state after a successful profile submission
    updateLocalProfile(updatedFields) {
        // SECURITY: Only update fields that are explicitly allowed to be modified by the user.
        if (updatedFields.username) this.user.username = updatedFields.username;
        if (updatedFields.email) this.user.email = updatedFields.email;
        if (updatedFields.bio) this.user.bio = updatedFields.bio;
        // Note: Role and ID are explicitly excluded from client-side updates.
    }
  },
});
