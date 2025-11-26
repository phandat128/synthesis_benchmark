<template>
  <Form @submit="onSubmit" :validation-schema="schema" v-slot="{ isSubmitting, errors }" class="profile-form">
    
    <div class="form-group">
      <label for="username">Username</label>
      <!-- SOURCE: Input field for username -->
      <Field name="username" type="text" v-model="localProfile.username" class="form-control" />
      <ErrorMessage name="username" class="error-message" />
    </div>

    <div class="form-group">
      <label for="email">Email</label>
      <!-- SOURCE: Input field for email -->
      <Field name="email" type="email" v-model="localProfile.email" class="form-control" />
      <ErrorMessage name="email" class="error-message" />
    </div>

    <div class="form-group">
      <label for="bio">Bio</label>
      <!-- SOURCE: Input field for bio -->
      <Field name="bio" as="textarea" v-model="localProfile.bio" class="form-control" rows="4" />
      <ErrorMessage name="bio" class="error-message" />
    </div>

    <!-- SECURITY: Explicitly exclude sensitive, non-editable fields (like role) from the form. -->
    <!-- Even if a hidden input for 'role' was added, the API service filter would drop it. -->

    <button type="submit" :disabled="isSubmitting" class="btn btn-primary">
      {{ isSubmitting ? 'Saving...' : 'Save Profile' }}
    </button>

    <div v-if="Object.keys(errors).length" class="alert alert-danger mt-3">
        Please correct the errors above.
    </div>

  </Form>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits } from 'vue';
import { Form, Field, ErrorMessage } from 'vee-validate';
import * as yup from 'yup';

const props = defineProps({
  initialData: {
    type: Object,
    required: true,
    default: () => ({ username: '', email: '', bio: '' })
  }
});

const emit = defineEmits(['submit']);

// Local state for form inputs
const localProfile = ref({ ...props.initialData });

// Watch for external changes to initialData (e.g., when profile is fetched)
watch(() => props.initialData, (newVal) => {
  localProfile.value = { ...newVal };
}, { deep: true });

// Input Validation Schema using yup
const schema = yup.object({
  // SECURITY: Strict validation rules for all user-editable fields.
  username: yup.string()
    .required('Username is required.')
    .min(3, 'Username must be at least 3 characters.')
    .max(50, 'Username cannot exceed 50 characters.')
    .matches(/^[a-zA-Z0-9_]+$/, 'Username can only contain letters, numbers, and underscores.'),
  
  email: yup.string()
    .required('Email is required.')
    .email('Must be a valid email address.')
    .max(100, 'Email cannot exceed 100 characters.'),

  bio: yup.string()
    .nullable()
    .max(500, 'Bio cannot exceed 500 characters.'),
});

// Handle form submission
const onSubmit = (values, { resetForm }) => {
  // SECURITY: The 'values' object here contains only the fields defined in the schema/form.
  // This object is the raw profileData mentioned in the taint flow source.
  // It is passed to the parent component, which then passes it to the API service for filtering.
  emit('submit', values);
};
</script>

<style scoped>
/* Basic styling for readability */
.form-group { margin-bottom: 15px; }
.error-message { color: red; font-size: 0.9em; }
.form-control { width: 100%; padding: 8px; border: 1px solid #ccc; border-radius: 4px; }
.btn-primary { padding: 10px 15px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
.btn-primary:disabled { background-color: #a0c9ff; }
</style>
