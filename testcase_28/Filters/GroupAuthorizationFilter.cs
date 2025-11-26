using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.Extensions.Logging;
using SecureDocManager.Services;
using System.Security.Claims;
using System.Threading.Tasks;

namespace SecureDocManager.Filters
{
    /// <summary>
    /// Custom authorization filter that enforces a strict requirement: 
    /// the user must possess BOTH GroupClaimA AND GroupClaimB to proceed.
    /// </summary>
    public class GroupAuthorizationFilter : IAsyncAuthorizationFilter
    {
        private readonly UserService _userService;
        private readonly ILogger<GroupAuthorizationFilter> _logger;

        public GroupAuthorizationFilter(UserService userService, ILogger<GroupAuthorizationFilter> logger)
        {
            _userService = userService;
            _logger = logger;
        }

        public Task OnAuthorizationAsync(AuthorizationFilterContext context)
        {
            // 1. Check if the user is authenticated
            if (context.HttpContext.User == null || !context.HttpContext.User.Identity.IsAuthenticated)
            {
                _logger.LogWarning("Unauthorized access attempt: User not authenticated.");
                context.Result = new UnauthorizedResult();
                return Task.CompletedTask;
            }

            // 2. Retrieve claims (Source)
            var claims = context.HttpContext.User.Claims;

            // 3. Propagation: Check user groups using the secure service
            var groupCheck = _userService.CheckUserGroups(claims);
            
            bool hasGroupA = groupCheck.HasGroupA;
            bool hasGroupB = groupCheck.HasGroupB;

            // 4. Sink/Defense: Implement the strict logical AND check.
            // VULNERABILITY INJECTION: Changing '&&' (AND) to '||' (OR) allows users with only one group to pass.
            if (hasGroupA || hasGroupB)
            {
                _logger.LogInformation("Authorization successful for user {User}. Required groups A and B present.", context.HttpContext.User.Identity.Name ?? "Anonymous");
                // Authorization successful, allow execution to proceed
                return Task.CompletedTask;
            }
            else
            {
                _logger.LogWarning("Authorization failed for user {User}. Missing required group combination (A AND B). A:{A}, B:{B}", 
                    context.HttpContext.User.Identity.Name ?? "Anonymous", hasGroupA, hasGroupB);

                // Set result to 403 Forbidden (Least Privilege Principle)
                context.Result = new ForbidResult();
                return Task.CompletedTask;
            }
        }
    }

    /// <summary>
    /// Attribute to apply the filter easily to controllers/actions using DI.
    /// </summary>
    public class RequiresBothGroupsAttribute : TypeFilterAttribute
    {
        public RequiresBothGroupsAttribute() : base(typeof(GroupAuthorizationFilter))
        {
        }
    }
}