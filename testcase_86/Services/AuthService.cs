using DocumentManager.Models;
using Blazored.LocalStorage;
using DocumentManager.Providers;
using Microsoft.AspNetCore.Components.Authorization;

namespace DocumentManager.Services
{
    public class AuthService
    {
        private readonly ILocalStorageService _localStorage;
        private readonly CustomAuthStateProvider _authStateProvider;
        private const string AuthTokenKey = "authToken";

        public AuthService(ILocalStorageService localStorage, AuthenticationStateProvider authStateProvider)
        {
            _localStorage = localStorage;
            // Downcast to access specific notification methods on our custom provider
            _authStateProvider = (CustomAuthStateProvider)authStateProvider;
        }

        public async Task<bool> LoginAsync(LoginRequest request)
        {
            // Secure Coding: Input validation (basic client-side check before API call)
            if (string.IsNullOrEmpty(request.Username) || string.IsNullOrEmpty(request.Password))
            {
                return false;
            }

            try
            {
                // --- SIMULATED BACKEND RESPONSE ---
                // In a real application, this would be an HttpClient POST request to the API.
                string simulatedToken = GenerateSimulatedJwt(request.Username, request.Password);

                if (!string.IsNullOrEmpty(simulatedToken))
                {
                    // Secure Storage: Store token and notify the authentication state provider
                    await _localStorage.SetItemAsStringAsync(AuthTokenKey, simulatedToken);
                    _authStateProvider.NotifyUserAuthentication(simulatedToken);
                    return true;
                }
                return false;
            }
            catch (Exception ex)
            {
                // Secure Coding: Log the error internally but return a generic failure message to the client.
                Console.WriteLine($"Login failed: {ex.Message}");
                return false;
            }
        }

        public async Task LogoutAsync()
        {
            // Secure Coding: Remove token from storage and reset authentication state
            await _localStorage.RemoveItemAsync(AuthTokenKey);
            _authStateProvider.NotifyUserLogout();
        }

        // Helper to simulate JWT generation with claims based on credentials
        private string GenerateSimulatedJwt(string username, string password)
        {
            // NOTE: Password check is simulated here. In production, this is handled by the server.
            if (username == "admin" && password == "securepass")
            {
                // User 'admin' gets both roles (GROUP_A, GROUP_B) - Authorized for sensitive docs
                // Payload: {"sub":"admin","role":["GROUP_A","GROUP_B"],"exp":253402300799}
                return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZG1pbiIsInJvbGUiOlsiR1JPVVBfQSIsIkdST1VQX0IiXSwiZXhwIjoyNTM0MDIzMDA3OTl9.i_g7tX5YyK0qL9r8Z2Qn3bX7gK1w4f5J0s_1X_2Y_3Z";
            }
            else if (username == "partial" && password == "weakpass")
            {
                // User 'partial' gets only one role (GROUP_A) - Unauthorized for sensitive docs
                // Payload: {"sub":"partial","role":["GROUP_A"],"exp":253402300799}
                return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwYXJ0aWFsIiwicm9sZSI6WyJHUk9VUF9BIl0sImV4cCI6MjUzNDAyMzAwNzk5fQ.q_R0s_T1uVwXyZc2d_3e_4f_5g_6h_7i_8j_9k_0l";
            }
            else if (username == "basic" && password == "basicpass")
            {
                // User 'basic' gets only one role (GROUP_B) - Unauthorized for sensitive docs
                // Payload: {"sub":"basic","role":["GROUP_B"],"exp":253402300799}
                return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiYXNpYyIsInJvbGUiOlsiR1JPVVBfQiJdLCJleHAiOjI1MzQwMjMwMDc5OX0.a_B1c_D2e_F3g_H4i_J5k_L6m_N7o_P8q_R9s_T0u_V1w";
            }
            // Authentication failure
            return null;
        }
    }
}