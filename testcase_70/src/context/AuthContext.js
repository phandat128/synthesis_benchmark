import React, { createContext, useContext, useState } from 'react';

const AuthContext = createContext(null);

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    // Simulated initial state for a standard user
    const [user, setUser] = useState({
        id: 'user_12345',
        username: 'john.doe',
        role: 'standard_user', // This field is sensitive
        isAuthenticated: true,
    });

    const value = {
        user,
        // In a real app, this would handle login/logout and token management
        login: (userData) => setUser({ ...userData, isAuthenticated: true }),
        logout: () => setUser({ isAuthenticated: false }),
        // Function to update local user state after a successful profile update
        updateLocalProfile: (newProfileData) => {
            setUser(prev => ({
                ...prev,
                ...newProfileData
            }));
        }
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};