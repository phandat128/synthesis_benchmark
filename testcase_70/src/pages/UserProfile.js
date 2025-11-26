import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import ProfileForm from '../components/ProfileForm';
import { useAuth } from '../context/AuthContext';
import { fetchUserProfile } from '../services/apiService';
import { useProfileUpdate } from '../hooks/useProfileUpdate';

const PageContainer = styled.div`
    padding: 40px;
    text-align: center;
`;

const LoadingText = styled.p`
    font-size: 1.2em;
    color: #007bff;
`;

const UserProfile = () => {
    const { user } = useAuth();
    const [profileData, setProfileData] = useState(null);
    const [fetchError, setFetchError] = useState(null);
    const [isFetching, setIsFetching] = useState(true);

    // Use the custom update hook, which handles security filtering
    const { handleUpdate, isLoading, error, successMessage } = useProfileUpdate();

    // Effect to fetch initial profile data
    useEffect(() => {
        if (user && user.id) {
            const loadProfile = async () => {
                setIsFetching(true);
                setFetchError(null);
                try {
                    const data = await fetchUserProfile(user.id);
                    setProfileData(data);
                } catch (err) {
                    // Do not expose raw error objects to the user
                    setFetchError(err.message || "Failed to load profile.");
                } finally {
                    setIsFetching(false);
                }
            };
            loadProfile();
        }
    }, [user]);

    if (!user.isAuthenticated) {
        return <PageContainer>Please log in to view your profile.</PageContainer>;
    }

    if (isFetching) {
        return <PageContainer><LoadingText>Loading profile data...</LoadingText></PageContainer>;
    }

    if (fetchError) {
        return <PageContainer><h1>Error</h1><p>{fetchError}</p></PageContainer>;
    }

    return (
        <PageContainer>
            <h1>User Profile Settings</h1>
            <ProfileForm
                initialData={profileData}
                onSubmit={handleUpdate}
                isLoading={isLoading}
                error={error}
                successMessage={successMessage}
            />
        </PageContainer>
    );
};

export default UserProfile;