using AuditSystem.Services;
using Microsoft.AspNetCore.Mvc;
using System.Threading.Tasks;

namespace AuditSystem.Controllers
{
    [ApiController]
    [Route("api/admin")]
    // NOTE: This controller MUST be protected by authorization middleware (e.g., [Authorize(Roles = "Admin")]) in a real system.
    public class AdminController : ControllerBase
    {
        private readonly IAuditService _auditService;

        public AdminController(IAuditService auditService)
        {
            _auditService = auditService;
        }

        /// <summary>
        /// Retrieves detailed audit logs for a specific user ID.
        /// This endpoint triggers the service method that retrieves and filters based on the stored username.
        /// </summary>
        /// <param name="userId">The ID of the user whose logs are requested.</param>
        [HttpGet("audit/user/{userId:int}")]
        public async Task<IActionResult> GetUserAuditLogs(int userId)
        {
            if (userId <= 0)
            {
                return BadRequest("Invalid User ID.");
            }

            // The service layer handles the secure retrieval and querying.
            var logs = await _auditService.GetDetailedAuditLog(userId);

            if (logs == null || logs.Count == 0)
            {
                // Proper error handling: Do not leak sensitive details if the user ID is valid but no logs exist.
                return NotFound($"No recent audit logs found for user ID {userId}.");
            }

            return Ok(logs);
        }
    }
}