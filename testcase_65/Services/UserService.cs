using SecureApi.Models;
using System.Threading.Tasks;
using System.Collections.Generic;
using System.Linq;
using System.Text.RegularExpressions;

namespace SecureApi.Services
{
    public interface IUserService
    {
        Task<bool> ChangeUserEmail(int userId, string newEmail);
        Task<UserProfile> GetUserById(int userId);
    }

    public class UserService : IUserService
    {
        // Mock Database Context (In-memory store for demonstration)
        private static readonly List<UserProfile> _users = new List<UserProfile>
        {
            new UserProfile { UserId = 1, Username = "Alice", Email = "alice@example.com" },
            new UserProfile { UserId = 2, Username = "Bob", Email = "bob@example.com" }
        };

        // Secure Email Validation using Regex with timeout to prevent ReDoS attacks
        private bool IsValidEmailFormat(string email)
        {
            if (string.IsNullOrWhiteSpace(email)) return false;
            try
            {
                // Standard robust email regex
                return Regex.IsMatch(email,
                    @"^[^@\s]+@[^@\s]+\.[^@\s]+$",
                    RegexOptions.IgnoreCase,
                    System.TimeSpan.FromMilliseconds(50)); // Set short timeout
            }
            catch (RegexMatchTimeoutException)
            {
                // Handle potential ReDoS attempt
                return false;
            }
        }

        /// <summary>
        /// Updates the user's email address after validation and sanitization.
        /// </summary>
        /// <param name="userId">The ID of the authenticated user.</param>
        /// <param name="newEmail">The new email address provided by the user.</param>
        /// <returns>True if update was successful.</returns>
        public async Task<bool> ChangeUserEmail(int userId, string newEmail)
        {
            // 1. Input Validation (Defense in Depth)
            if (!IsValidEmailFormat(newEmail))
            {
                return false; // Invalid format
            }

            // 2. Business Logic & Persistence
            var user = await GetUserById(userId);

            if (user == null)
            {
                // User not found (should be handled by authorization, but check defensively)
                return false;
            }

            // 3. Input Sanitization
            // Ensure data is clean before persistence
            string sanitizedEmail = newEmail.Trim().ToLowerInvariant();

            // 4. Check for uniqueness (if required)
            if (_users.Any(u => u.Email == sanitizedEmail && u.UserId != userId))
            {
                return false; // Email already in use
            }

            // Perform the state change (the sink operation)
            user.Email = sanitizedEmail;

            // Simulate saving changes to the database (using ORM/Parameterized query in a real scenario)
            await Task.Delay(10); 
            
            return true;
        }

        public async Task<UserProfile> GetUserById(int userId)
        {
            return await Task.FromResult(_users.FirstOrDefault(u => u.UserId == userId));
        }
    }
}