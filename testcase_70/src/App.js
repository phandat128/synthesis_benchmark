import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import styled, { createGlobalStyle } from 'styled-components';
import { AuthProvider } from './context/AuthContext';
import UserProfile from './pages/UserProfile';

// Global styles for clean rendering
const GlobalStyle = createGlobalStyle`
    body {
        margin: 0;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
            'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
            sans-serif;
        background-color: #f4f7f9;
    }
`;

const Header = styled.header`
    background-color: #333;
    padding: 15px 30px;
    color: white;
    display: flex;
    justify-content: space-between;
    align-items: center;
`;

const NavLink = styled(Link)`
    color: white;
    text-decoration: none;
    margin-left: 20px;
    font-weight: bold;
    &:hover {
        text-decoration: underline;
    }
`;

const App = () => {
    return (
        <Router>
            <GlobalStyle />
            <AuthProvider>
                <Header>
                    <div>Secure Profile App</div>
                    <nav>
                        <NavLink to="/">Home</NavLink>
                        <NavLink to="/profile">Profile</NavLink>
                    </nav>
                </Header>
                <Routes>
                    <Route path="/" element={
                        <div style={{ padding: '40px', textAlign: 'center' }}>
                            <h2>Welcome to the Secure Application</h2>
                            <p>Navigate to the Profile page to manage your settings.</p>
                        </div>
                    } />
                    <Route path="/profile" element={<UserProfile />} />
                </Routes>
            </AuthProvider>
        </Router>
    );
};

export default App;