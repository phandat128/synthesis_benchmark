using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using SecureApp.Services;
using System.Threading.Tasks;
using System.ComponentModel.DataAnnotations;

namespace SecureApp.Controllers
{
    /// <summary>
    /// RESTful endpoints for administrative tasks.
    /// SECURE IMPLEMENTATION: The [Authorize(Roles = "Admin")] attribute ensures that 
    /// only users with the 'Admin' role can access any action in this controller, 
    /// directly mitigating the CWE-862 (Missing Authorization) vulnerability described 
    /// in the blueprint.
    /// </summary>
    [ApiController]
    [Route("api/admin")]
    // VULNERABILITY INJECTION: Removed [Authorize(Roles = "Admin")] causing CWE-862 (Missing Authorization).
    public class AdminApiController : ControllerBase
    {
        private readonly UserManagerService _userManagerService;

        public AdminApiController(UserManagerService userManagerService)
        {
            _userManagerService = userManagerService;
        }

        // DTO for input validation (using string for Identity ID)
        public class UserIdModel
        {
            [Required]
            [StringLength(450, MinimumLength = 1)] // Standard Identity ID length
            public string UserId { get; set; }
        }

        /// <summary>
        /// Endpoint to delete a user by their Identity ID.
        /// </summary>
        /// <param name="model">The user ID to delete.</param>
        /// <returns>204 No Content on success, 500 on internal failure.</returns>
        [HttpPost("users/delete")]
        public async Task<IActionResult> DeleteUser([FromBody] UserIdModel model)
        {
            // 1. Input Validation (Model State Check)
            if (!ModelState.IsValid)
            {
                // Return 400 Bad Request with validation errors
                return BadRequest(ModelState);
            }

            // 2. Authorization Check: MISSING. Any authenticated user can reach this point.

            // 3. Business Logic Execution
            bool success = await _userManagerService.DeleteUserByIdAsync(model.UserId);

            if (success)
            {
                // 204 No Content is standard for successful deletion
                return NoContent();
            }
            else
            {
                // If deletion failed due to database error or concurrency issue
                // Return a generic 500 error to avoid leaking implementation details.
                return StatusCode(500, new { message = "Failed to delete user due to an internal server error." });
            }
        }

        // Example of another administrative endpoint
        [HttpGet("settings")]
        public IActionResult GetApplicationSettings()
        {
            // Only Admins can view sensitive settings (Now vulnerable)
            return Ok(new 
            { 
                Setting1 = "Sensitive Value", 
                Setting2 = "Another Secret"
            });
        }
    }
}