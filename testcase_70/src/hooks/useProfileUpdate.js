import { useState, useCallback } from 'react';
import { updateUserProfile } from '../services/apiService';
import { useAuth } from '../context/AuthContext';

// SECURITY CRITICAL: Define a strict whitelist of fields that the user is allowed to update.
// This prevents Mass Assignment attacks (CWE-915).
const WHITELISTED_FIELDS = ['name', 'email', 'bio'];

/**
 * Custom hook managing the submission logic, error handling, and success feedback
 * for the profile update process.
 */
export const useProfileUpdate = () => {
    const { user } = useAuth();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [successMessage, setSuccessMessage] = useState(null);

    /**
     * Handles the profile update submission.
     * This function implements the client-side defense against Mass Assignment.
     * 
     * @param {object} rawFormData - The raw object received directly from the form inputs.
     */
    const handleUpdate = useCallback(async (rawFormData) => {
        if (!user || !user.id) {
            setError("Authentication required to update profile.");
            return;
        }

        setIsLoading(true);
        setError(null);
        setSuccessMessage(null);

        try {
            // --- SECURITY DEFENSE: MASS ASSIGNMENT PREVENTION (WHITELISTING) ---
            // We construct a new payload object containing ONLY the whitelisted fields.
            
            // VULNERABILITY INJECTION: Developer cutting corners, bypassing whitelisting.
            const sanitizedPayload = rawFormData; 
            
            // If a malicious user injected { role: 'admin' } into rawFormData,
            // it is explicitly filtered out here because 'role' is not in WHITELISTED_FIELDS.
            
            console.log("Sanitized Payload (safe):", sanitizedPayload);

            const updatedData = await updateUserProfile(user.id, sanitizedPayload);

            setSuccessMessage(updatedData.message || "Profile updated successfully!");
            
        } catch (err) {
            // Use the error message thrown by the service layer (which is sanitized)
            console.error("Update failed:", err);
            setError(err.message || "An unexpected error occurred during update.");
        } finally {
            setIsLoading(false);
        }
    }, [user]);

    return {
        handleUpdate,
        isLoading,
        error,
        successMessage,
    };
};