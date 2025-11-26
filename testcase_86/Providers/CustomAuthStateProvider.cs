using Microsoft.AspNetCore.Components.Authorization;
using System.Security.Claims;
using System.Text.Json;
using Blazored.LocalStorage;

namespace DocumentManager.Providers
{
    public class CustomAuthStateProvider : AuthenticationStateProvider
    {
        private readonly ILocalStorageService _localStorage;
        private readonly HttpClient _httpClient;
        private const string AuthTokenKey = "authToken";

        public CustomAuthStateProvider(ILocalStorageService localStorage, HttpClient httpClient)
        {
            _localStorage = localStorage;
            _httpClient = httpClient;
        }

        public override async Task<AuthenticationState> GetAuthenticationStateAsync()
        {
            var savedToken = await _localStorage.GetItemAsStringAsync(AuthTokenKey);

            if (string.IsNullOrWhiteSpace(savedToken))
            {
                // Return anonymous user if no token is found
                return new AuthenticationState(new ClaimsPrincipal(new ClaimsIdentity()));
            }

            // Secure Coding: Attach token to HttpClient for subsequent API requests
            _httpClient.DefaultRequestHeaders.Authorization = 
                new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", savedToken);

            var claimsPrincipal = new ClaimsPrincipal(new ClaimsIdentity(ParseClaimsFromJwt(savedToken), "jwtAuth"));
            return new AuthenticationState(claimsPrincipal);
        }

        public void NotifyUserAuthentication(string token)
        {
            var authenticatedUser = new ClaimsPrincipal(new ClaimsIdentity(ParseClaimsFromJwt(token), "jwtAuth"));
            var authState = Task.FromResult(new AuthenticationState(authenticatedUser));
            NotifyAuthenticationStateChanged(authState);
        }

        public void NotifyUserLogout()
        {
            var anonymousUser = new ClaimsPrincipal(new ClaimsIdentity());
            var authState = Task.FromResult(new AuthenticationState(anonymousUser));
            NotifyAuthenticationStateChanged(authState);
            
            // Clear the HttpClient header upon logout
            _httpClient.DefaultRequestHeaders.Authorization = null;
        }

        private static IEnumerable<Claim> ParseClaimsFromJwt(string jwt)
        {
            var claims = new List<Claim>();
            
            // JWT structure: Header.Payload.Signature
            var parts = jwt.Split('.');
            if (parts.Length != 3) return claims; // Malformed JWT

            var payload = parts[1];
            
            // Secure Coding: JWT parsing must handle malformed tokens gracefully.
            try
            {
                var jsonBytes = ParseBase64WithoutPadding(payload);
                var keyValuePairs = JsonSerializer.Deserialize<Dictionary<string, object>>(jsonBytes);

                // Extract standard claims
                if (keyValuePairs.TryGetValue("sub", out object username)) 
                {
                    claims.Add(new Claim(ClaimTypes.Name, username.ToString()));
                }

                // Extract roles/groups (using ClaimTypes.Role)
                if (keyValuePairs.TryGetValue("role", out object roles))
                {
                    if (roles is JsonElement roleElement)
                    {
                        if (roleElement.ValueKind == JsonValueKind.Array)
                        {
                            // Handle multiple roles (e.g., ["GROUP_A", "GROUP_B"])
                            claims.AddRange(roleElement.EnumerateArray().Select(r => new Claim(ClaimTypes.Role, r.ToString())));
                        }
                        else if (roleElement.ValueKind == JsonValueKind.String)
                        {
                            // Handle single role string
                            claims.Add(new Claim(ClaimTypes.Role, roleElement.ToString()));
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                // Secure Coding: If parsing fails (e.g., token tampered or expired), treat as unauthenticated.
                Console.WriteLine($"Error parsing JWT claims: {ex.Message}");
            }

            return claims;
        }

        // Helper method to correctly decode base64 strings that might lack padding
        private static byte[] ParseBase64WithoutPadding(string base64)
        {
            switch (base64.Length % 4)
            {
                case 2: base64 += "=="; break;
                case 3: base64 += "="; break;
            }
            return Convert.FromBase64String(base64);
        }
    }
}