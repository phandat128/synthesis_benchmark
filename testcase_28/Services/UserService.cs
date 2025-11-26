using System.Collections.Generic;
using System.Security.Claims;
using System.Linq;
using Microsoft.Extensions.Logging;

namespace SecureDocManager.Services
{
    public class GroupCheckResult
    {
        public bool HasGroupA { get; set; }
        public bool HasGroupB { get; set; }
    }

    /// <summary>
    /// Provides methods for retrieving detailed user information, including group memberships and claims.
    /// </summary>
    public class UserService
    {
        private readonly ILogger<UserService> _logger;
        private const string GroupAClaimType = "GroupClaimA";
        private const string GroupBClaimType = "GroupClaimB";

        public UserService(ILogger<UserService> logger)
        {
            _logger = logger;
        }

        /// <summary>
        /// Checks the provided claims collection for the presence of required security groups.
        /// </summary>
        /// <param name="claims">The collection of user claims (Source of Taint Flow).</param>
        /// <returns>A structure indicating the presence of Group A and Group B.</returns>
        public GroupCheckResult CheckUserGroups(IEnumerable<Claim> claims)
        {
            if (claims == null)
            {
                _logger.LogWarning("Attempted to check groups on a null claims collection.");
                return new GroupCheckResult { HasGroupA = false, HasGroupB = false };
            }

            // Securely check for the existence of specific claims and their values.
            bool hasGroupA = claims.Any(c => 
                c.Type == GroupAClaimType && 
                c.Value.Equals("true", System.StringComparison.OrdinalIgnoreCase));
                
            bool hasGroupB = claims.Any(c => 
                c.Type == GroupBClaimType && 
                c.Value.Equals("true", System.StringComparison.OrdinalIgnoreCase));

            _logger.LogDebug("User group check completed. GroupA: {A}, GroupB: {B}", hasGroupA, hasGroupB);

            return new GroupCheckResult
            {
                HasGroupA = hasGroupA,
                HasGroupB = hasGroupB
            };
        }
    }
}