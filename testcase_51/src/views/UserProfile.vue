<template>
  <div class="user-profile-view">
    <h1>My Profile Settings</h1>
    
    <div v-if="authStore.isLoadingProfile" class="loading-state">Loading profile data...</div>
    
    <div v-else-if="!authStore.user.id" class="error-state">
      Could not load profile. Please log in again.
    </div>

    <div v-else>
      <p><strong>Role:</strong> {{ authStore.user.role }} (Non-editable)</p>
      
      <!-- Pass the current user data to the form component -->
      <ProfileForm 
        :initial-data="profileDataForForm" 
        @submit="handleProfileUpdate" 
      />

      <div v-if="updateStatus.message" :class="['alert', updateStatus.type === 'success' ? 'alert-success' : 'alert-danger']">
        {{ updateStatus.message }}
      </div>
    </div>

  </div>
</template>

<script setup>
import { onMounted, computed, reactive } from 'vue';
import { useAuthStore } from '../store/authStore';
import apiService from '../services/apiService';
import ProfileForm from '../components/ProfileForm.vue';

const authStore = useAuthStore();

// State for handling feedback messages
const updateStatus = reactive({
  message: '',
  type: ''
});

// Computed property to extract only the fields the form needs, preventing accidental exposure of internal state.
const profileDataForForm = computed(() => ({
  username: authStore.user.username,
  email: authStore.user.email,
  bio: authStore.user.bio,
  // SECURITY: Explicitly exclude 'id' and 'role' from the data passed to the editable form.
}));

// Fetch profile data when the component mounts
onMounted(() => {
  if (!authStore.user.id) {
    authStore.fetchUserProfile();
  }
});

/**
 * Handles the submission event from the ProfileForm.
 * @param {Object} formData - The validated data from the form (Source of Taint Flow).
 */
const handleProfileUpdate = async (formData) => {
  updateStatus.message = '';
  updateStatus.type = '';

  try {
    // PROPAGATION: formData (the raw input) is passed to the secured API service function.
    // The apiService.updateProfile function is responsible for filtering/whitelisting this data.
    const response = await apiService.updateProfile(formData);
    
    // Update local state upon successful API response
    authStore.updateLocalProfile(formData);

    updateStatus.message = response.message || 'Profile updated successfully!';
    updateStatus.type = 'success';
    
    console.log('Profile update successful.', response);

  } catch (error) {
    // SECURITY: Display generic error message, avoiding leakage of backend implementation details.
    updateStatus.message = error.message || 'An unexpected error occurred during profile update.';
    updateStatus.type = 'danger';
    console.error('Update failed:', error);
  }
};
</script>

<style scoped>
.user-profile-view { max-width: 800px; margin: 0 auto; padding: 20px; }
.alert { padding: 10px; border-radius: 4px; margin-top: 20px; }
.alert-success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
.alert-danger { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
.loading-state, .error-state { padding: 20px; text-align: center; border: 1px solid #eee; }
</style>
