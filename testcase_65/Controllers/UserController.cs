using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Authorization;
using SecureApi.DTOs;
using SecureApi.Services;
using System.Security.Claims;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Antiforgery;

namespace SecureApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class UserController : ControllerBase
    {
        private readonly IUserService _userService;
        private readonly IAntiforgery _antiforgery;

        public UserController(IUserService userService, IAntiforgery antiforgery)
        {
            _userService = userService;
            _antiforgery = antiforgery;
        }

        /// <summary>
        /// Endpoint to retrieve the CSRF token. Required for SPA clients using cookie authentication.
        /// The token is set in the 'XSRF-REQUEST-TOKEN' cookie and must be read by the client 
        /// and sent back in the 'X-CSRF-TOKEN' header for state-changing requests.
        /// </summary>
        [HttpGet("csrf-token")]
        [Authorize] // Only authenticated users should get a token
        public IActionResult GetCsrfToken()
        {
            // This call ensures the anti-forgery token cookie is set in the response.
            _antiforgery.GetAndStoreTokens(HttpContext);
            
            return Ok(new { message = "CSRF token cookie set. Client must read cookie and send header." });
        }


        /// <summary>
        /// Updates the authenticated user's email address.
        /// </summary>
        /// <param name="request">The DTO containing the new email address.</param>
        [HttpPost("update-email")]
        [Authorize] // 1. Ensures only authenticated users can access this
        public async Task<IActionResult> UpdateEmail([FromBody] UpdateEmailRequest request)
        {
            // 3. Input Validation (Handles Data Annotations from DTO)
            if (!ModelState.IsValid)
            {
                // Return 400 Bad Request with specific validation errors, avoiding sensitive data leakage
                return BadRequest(ModelState);
            }

            // 4. Identity Retrieval (Least Privilege Principle)
            // Get the authenticated user's ID from the claims (secure source)
            var userIdClaim = User.FindFirst(ClaimTypes.NameIdentifier);
            if (userIdClaim == null || !int.TryParse(userIdClaim.Value, out int userId))
            {
                // If authentication succeeded but ID claim is missing, something is wrong.
                return Unauthorized("User identity could not be determined.");
            }

            // 5. Business Logic Execution
            var success = await _userService.ChangeUserEmail(userId, request.NewEmailAddress);

            if (success)
            {
                return Ok(new { message = "Email updated successfully." });
            }
            else
            {
                // Proper Error Handling: Return a generic 400 if the service failed due to business rules (e.g., invalid format, duplicate)
                return BadRequest(new { error = "Failed to update email. Please ensure the address is valid and unique." });
            }
        }
    }
}