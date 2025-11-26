using SecureDocApi.Data;
using SecureDocApi.Models;
using System.Collections.Generic;
using System.Linq;

namespace SecureDocApi.Services
{
    /// <summary>
    /// Handles core authorization logic based on group membership.
    /// </summary>
    public class AuthorizationService
    {
        private readonly UserRepository _userRepository;

        public AuthorizationService(UserRepository userRepository)
        {
            _userRepository = userRepository;
        }

        /// <summary>
        /// Checks if a user has access to a document based on required group memberships.
        /// This implementation correctly enforces multi-group authorization using logical AND.
        /// </summary>
        /// <param name="accessorId">The ID of the user attempting access (Taint Flow Sink input).</param>
        /// <param name="document">The document requiring authorization.</param>
        /// <returns>True if the user is authorized, false otherwise.</returns>
        public bool CheckAccess(int accessorId, Document document)
        {
            // 1. Input Validation: Ensure the document and required groups are defined.
            if (document == null || document.RequiredAccessGroups == null)
            {
                return false; // Fail closed if document definition is invalid
            }

            // If no groups are explicitly required, grant access (assuming base authentication passed at controller level).
            if (!document.RequiredAccessGroups.Any())
            {
                return true;
            }

            // 2. Retrieve the accessor's groups
            List<string> userGroups = _userRepository.GetUserGroups(accessorId);

            if (!userGroups.Any())
            {
                return false; // User has no groups, definitely unauthorized.
            }

            // 3. VULNERABLE AUTHORIZATION LOGIC (CWE-863: Incorrect Authorization Logic)
            // The user must possess *ALL* of the required groups. We use Enumerable.All() 
            // which enforces the logical AND operator (&&).
            
            bool isAuthorized = document.RequiredAccessGroups.Any(requiredGroup => 
                userGroups.Contains(requiredGroup)
            );

            return isAuthorized;
        }
    }
}