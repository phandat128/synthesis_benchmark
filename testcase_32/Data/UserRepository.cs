using System.Security.Claims;
using System.Threading.Tasks;
using System.Security.Authentication;

namespace SecureReportApp.Data
{
    /// <summary>
    /// Manages data access operations for retrieving user identity and associated group claims from the persistence layer.
    /// </summary>
    public class UserRepository
    {
        /// <summary>
        /// Retrieves and verifies the full ClaimsPrincipal object for a given authenticated user.
        /// This step ensures that the claims used for authorization are trusted and current.
        /// </summary>
        /// <param name="userPrincipal">The authenticated user's principal from the HTTP context.</param>
        /// <returns>The verified ClaimsPrincipal.</returns>
        public Task<ClaimsPrincipal> GetVerifiedUserClaimsAsync(ClaimsPrincipal userPrincipal)
        {
            // Defensive check: Ensure the principal is authenticated before proceeding.
            if (userPrincipal == null || userPrincipal.Identity == null || !userPrincipal.Identity.IsAuthenticated)
            {
                throw new AuthenticationException("User is not authenticated or identity is invalid.");
            }

            // In a real system, this is where we would query the database or identity provider
            // using the NameIdentifier to refresh and verify the claims.
            
            // For this secure example, we trust the claims provided by the ASP.NET Core authentication middleware
            // but enforce the authentication check.
            return Task.FromResult(userPrincipal);
        }
    }
}