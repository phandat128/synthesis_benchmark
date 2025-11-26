using Microsoft.AspNetCore.Authorization;
using SecureDocumentApi.Services;
using System.Threading.Tasks;

namespace SecureDocumentApi.Filters
{
    /// <summary>
    /// Custom Authorization Requirement defining the need for multiple specific groups.
    /// This is used by the policy definition in Program.cs.
    /// </summary>
    public class RequiredGroupRequirement : IAuthorizationRequirement
    {
        public string[] RequiredGroups { get; }
        public RequiredGroupRequirement(params string[] requiredGroups)
        {
            RequiredGroups = requiredGroups;
        }
    }

    /// <summary>
    /// Custom Authorization Handler responsible for enforcing the RequiredGroupRequirement.
    /// This component hooks into the ASP.NET Core authorization pipeline to perform the secure group check.
    /// </summary>
    public class RequiredGroupAuthorizationHandler : AuthorizationHandler<RequiredGroupRequirement>
    {
        private readonly AuthorizationService _authorizationService;

        public RequiredGroupAuthorizationHandler(AuthorizationService authorizationService)
        {
            _authorizationService = authorizationService;
        }

        /// <summary>
        /// Handles the authorization requirement check.
        /// </summary>
        protected override Task HandleRequirementAsync(AuthorizationHandlerContext context, RequiredGroupRequirement requirement)
        {
            // Taint Flow Source: Claims extracted from the validated JWT token.
            var userClaims = context.User.Claims;

            // Taint Flow Sink (Secure Check):
            bool isAuthorized = _authorizationService.CheckRequiredGroups(userClaims, requirement.RequiredGroups);

            if (isAuthorized)
            {
                // Authorization succeeds only if ALL groups are present (CWE-863 defense).
                context.Succeed(requirement);
            }
            
            return Task.CompletedTask;
        }
    }
}