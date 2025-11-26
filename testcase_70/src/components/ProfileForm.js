import React, { useState, useEffect } from 'react';
import styled from 'styled-components';

const FormContainer = styled.form`
    display: flex;
    flex-direction: column;
    gap: 15px;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 8px;
    max-width: 400px;
    margin: 20px auto;
    background-color: white;
`;

const InputGroup = styled.div`
    display: flex;
    flex-direction: column;
`;

const Label = styled.label`
    margin-bottom: 5px;
    font-weight: bold;
`;

const Input = styled.input`
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
`;

const TextArea = styled.textarea`
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    min-height: 100px;
`;

const Button = styled.button`
    padding: 10px 15px;
    background-color: ${props => (props.disabled ? '#aaa' : '#007bff')};
    color: white;
    border: none;
    border-radius: 4px;
    cursor: ${props => (props.disabled ? 'not-allowed' : 'pointer')};
    &:hover {
        background-color: ${props => (props.disabled ? '#aaa' : '#0056b3')};
    }
`;

const ErrorMessage = styled.p`
    color: red;
    font-weight: bold;
`;

const SuccessMessage = styled.p`
    color: green;
    font-weight: bold;
`;

/**
 * Renders the profile input form and manages local state.
 * Handles the source of the potential vulnerability (user input).
 * 
 * @param {object} initialData - The profile data fetched from the API.
 * @param {function} onSubmit - Function to call on form submission (passed from hook).
 * @param {boolean} isLoading - Loading state.
 * @param {string|null} error - Error message.
 * @param {string|null} successMessage - Success message.
 */
const ProfileForm = ({ initialData, onSubmit, isLoading, error, successMessage }) => {
    // Initialize form state using only the fields the user is allowed to modify
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        bio: '',
    });

    // Sync initial data when it loads
    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name || '',
                email: initialData.email || '',
                bio: initialData.bio || '',
            });
        }
    }, [initialData]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        // Input Validation: Basic length and type checks can be added here
        setFormData(prev => ({
            ...prev,
            [name]: value,
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        
        // Client-side validation (basic checks)
        if (!formData.name || !formData.email) {
            alert("Name and Email are required.");
            return;
        }

        // The formData object contains only the fields rendered by this component (name, email, bio).
        // Even if a malicious user modifies the client-side JS to inject 'role: admin' here,
        // the useProfileUpdate hook will filter it out before sending it to the API.
        onSubmit(formData);
    };

    return (
        <FormContainer onSubmit={handleSubmit}>
            <h2>Edit Profile</h2>

            {error && <ErrorMessage>Error: {error}</ErrorMessage>}
            {successMessage && <SuccessMessage>{successMessage}</SuccessMessage>}

            <InputGroup>
                <Label htmlFor="name">Full Name</Label>
                <Input
                    id="name"
                    name="name"
                    type="text"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    maxLength={100}
                />
            </InputGroup>

            <InputGroup>
                <Label htmlFor="email">Email Address</Label>
                <Input
                    id="email"
                    name="email"
                    type="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    // Client-side regex validation for email format
                    pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$"
                />
            </InputGroup>

            <InputGroup>
                <Label htmlFor="bio">Biography</Label>
                <TextArea
                    id="bio"
                    name="bio"
                    value={formData.bio}
                    onChange={handleChange}
                    maxLength={500}
                />
            </InputGroup>

            {/* Display non-editable, sensitive data for context */}
            {initialData && initialData.role && (
                <p style={{ fontSize: '0.8em', color: '#666' }}>
                    Current Role (Non-Editable): <strong>{initialData.role}</strong>
                </p>
            )}


            <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Saving...' : 'Save Changes'}
            </Button>
        </FormContainer>
    );
};

export default ProfileForm;