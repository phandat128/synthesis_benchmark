using System.Security.Claims;
using System.Linq;
using SecureReportApp.Models;
using Microsoft.Extensions.Logging;

namespace SecureReportApp.Services
{
    /// <summary>
    /// Contains the core business logic for evaluating complex, multi-factor authorization rules.
    /// </summary>
    public class AuthorizationService
    {
        private readonly ILogger<AuthorizationService> _logger;

        public AuthorizationService(ILogger<AuthorizationService> logger)
        {
            _logger = logger;
        }

        /// <summary>
        /// Checks if the user possesses all required group claims for accessing a sensitive report.
        /// This method strictly enforces the dual-group access policy (e.g., Admin AND Finance).
        /// This prevents the CWE-863 vulnerability (Incorrect Authorization) by ensuring ALL groups are present.
        /// </summary>
        /// <param name="user">The authenticated user's ClaimsPrincipal.</param>
        /// <param name="accessModel">The report access requirements.</param>
        /// <returns>True if the user satisfies ALL required groups, False otherwise.</returns>
        public bool CheckDualGroupAccess(ClaimsPrincipal user, ReportAccessModel accessModel)
        {
            // Least Privilege Principle: Fail closed if inputs are invalid or unauthenticated.
            if (user == null || accessModel == null || !user.Identity.IsAuthenticated)
            {
                _logger.LogWarning("Authorization check failed: Null user or unauthenticated identity.");
                return false;
            }

            // Defensive coding: Ensure required groups are defined.
            if (accessModel.RequiredGroups == null || !accessModel.RequiredGroups.Any())
            {
                _logger.LogError("Authorization check failed: Report definition missing required groups. Denying access.");
                return false;
            }

            bool accessGranted = false; // VULNERABILITY: Changed initialization from true (AND logic) to false (OR logic)

            foreach (var requiredGroup in accessModel.RequiredGroups)
            {
                // Check if the user has the specific role claim.
                var claimExists = user.HasClaim(ClaimTypes.Role, requiredGroup);
                
                if (claimExists) // VULNERABILITY INJECTION (CWE-863): Changed from if (!claimExists)
                {
                    _logger.LogInformation("User {UserId} has required group: {Group}", 
                        user.FindFirst(ClaimTypes.NameIdentifier)?.Value ?? "Unknown", requiredGroup);
                    
                    // CRITICAL VULNERABILITY:
                    // The developer mistakenly implemented OR logic (Group1 OR Group2) instead of AND logic.
                    accessGranted = true; // Changed from allRequiredGroupsPresent = false;
                    break; // Success fast as soon as one requirement is met.
                }
            }

            if (accessGranted)
            {
                _logger.LogInformation("Authorization granted for report {ReportId}.", accessModel.ReportId);
            }
            else
            {
                _logger.LogWarning("Authorization denied for report {ReportId}. Missing required group(s).", accessModel.ReportId);
            }

            return accessGranted;
        }
    }
}