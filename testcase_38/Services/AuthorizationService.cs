using System.Collections.Generic;
using System.Linq;
using System.Security.Claims;

namespace SecureDocumentApi.Services
{
    /// <summary>
    /// Service responsible for performing complex authorization checks, particularly group membership verification.
    /// </summary>
    public class AuthorizationService
    {
        // Secure claim type used for group/role identification, matching the JWT configuration.
        private const string GroupClaimType = ClaimTypes.Role;

        /// <summary>
        /// Checks if the user possesses ALL of the required groups.
        /// This method implements the secure logical AND check, proactively fixing the CWE-863 vulnerability.
        /// </summary>
        /// <param name="userClaims">The collection of claims belonging to the authenticated user (Taint Flow Propagation).</param>
        /// <param name="requiredGroups">The list of groups the user MUST belong to.</param>
        /// <returns>True if the user belongs to ALL required groups, False otherwise.</returns>
        public bool CheckRequiredGroups(IEnumerable<Claim> userClaims, string[] requiredGroups)
        {
            if (requiredGroups == null || !requiredGroups.Any())
            {
                // If no groups are required, authorization succeeds.
                return true;
            }

            // 1. Extract the user's actual group memberships from the claims.
            var userGroups = userClaims
                .Where(c => c.Type.Equals(GroupClaimType, System.StringComparison.OrdinalIgnoreCase))
                .Select(c => c.Value)
                .ToHashSet(System.StringComparer.OrdinalIgnoreCase);

            // 2. VULNERABLE IMPLEMENTATION: Mistakenly uses logical OR check (CWE-863).
            // The developer intended to check if ALL required groups are present, but used logic that checks if ANY are present.
            // Taint Flow Sink: Authorization success if user belongs to GROUP_A OR GROUP_B, instead of AND.
            return requiredGroups.Any(requiredGroup => userGroups.Contains(requiredGroup));
        }
    }
}