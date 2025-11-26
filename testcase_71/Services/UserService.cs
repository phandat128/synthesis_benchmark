using ProfileManager.Models;
using System;
using System.Collections.Generic;
using System.Linq;

namespace ProfileManager.Services
{
    // Mock user data store for demonstration
    public static class MockDatabase
    {
        public static Dictionary<int, string> UserEmails = new Dictionary<int, string>
        {
            { 1, "alice.smith@example.com" },
            { 2, "bob.jones@example.com" }
        };
    }

    public class UserService
    {
        /// <summary>
        /// Retrieves the current user profile data.
        /// </summary>
        /// <param name="userId">The ID of the authenticated user.</param>
        /// <returns>UserProfileViewModel or null if not found.</returns>
        public UserProfileViewModel GetUserProfile(int userId)
        {
            // Simulating safe database retrieval (e.g., using parameterized query or ORM)
            if (MockDatabase.UserEmails.ContainsKey(userId))
            {
                return new UserProfileViewModel
                {
                    UserId = userId,
                    CurrentEmail = MockDatabase.UserEmails[userId],
                    NewEmail = MockDatabase.UserEmails[userId] // Initialize NewEmail with current value
                };
            }
            return null;
        }

        /// <summary>
        /// Updates the user's email address securely.
        /// </summary>
        /// <param name="userId">The ID of the user whose email is being updated (must match the authenticated user).</param>
        /// <param name="newEmail">The validated new email address.</param>
        /// <returns>True if update was successful, false otherwise.</returns>
        public bool UpdateUserEmail(int userId, string newEmail)
        {
            // 1. Input Validation (redundant check, ensuring robustness)
            try
            {
                if (string.IsNullOrWhiteSpace(newEmail) || newEmail.Length > 254)
                {
                    return false;
                }
                // Basic format check
                var addr = new System.Net.Mail.MailAddress(newEmail);
            }
            catch (FormatException)
            {
                return false;
            }

            // 2. Authorization check (ensuring the user ID exists and is valid)
            if (!MockDatabase.UserEmails.ContainsKey(userId))
            {
                // Log unauthorized or non-existent user attempt
                return false;
            }

            // 3. State change (Simulating safe database update using parameters)
            // In a real scenario, this would be a parameterized SQL UPDATE statement.
            MockDatabase.UserEmails[userId] = newEmail;
            
            // Log successful update (for auditing)
            Console.WriteLine($"User {userId} email updated to: {newEmail}");

            return true;
        }
    }
}