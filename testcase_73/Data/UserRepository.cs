using System.Collections.Generic;
using System.Linq;

namespace SecureDocApi.Data
{
    /// <summary>
    /// Simulates a repository for fetching user data (groups).
    /// In a real application, this would use parameterized queries against a DB.
    /// </summary>
    public class UserRepository
    {
        // Simulated User Group Data:
        private static readonly Dictionary<int, List<string>> UserGroups = new Dictionary<int, List<string>>
        {
            { 101, new List<string> { "ADMIN", "GROUP_A" } }, // Has A, lacks B
            { 102, new List<string> { "GROUP_B", "READER" } }, // Has B, lacks A
            { 103, new List<string> { "GROUP_A", "GROUP_B", "AUDITOR" } }, // Full access user (Has A AND B)
            { 104, new List<string> { "GUEST" } } // No access user
        };

        /// <summary>
        /// Retrieves the list of security groups a user belongs to.
        /// </summary>
        /// <param name="userId">The ID of the user.</param>
        /// <returns>A list of group names, or an empty list if the user is not found.</returns>
        public List<string> GetUserGroups(int userId)
        {
            // Input validation and sanitization are handled implicitly by the int type and dictionary lookup.
            if (UserGroups.ContainsKey(userId))
            {
                return UserGroups[userId];
            }
            // Return empty list instead of null for safety and easier iteration
            return new List<string>();
        }
    }
}