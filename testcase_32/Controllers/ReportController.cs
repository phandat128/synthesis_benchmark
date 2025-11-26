using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Authorization;
using System.Security.Claims;
using SecureReportApp.Services;
using SecureReportApp.Models;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;
using SecureReportApp.Data;

namespace SecureReportApp.Controllers
{
    [Authorize] // Ensure only authenticated users can access this controller
    [Route("[controller]")]
    [ApiController]
    public class ReportController : ControllerBase
    {
        private readonly AuthorizationService _authService;
        private readonly ILogger<ReportController> _logger;
        private readonly UserRepository _userRepository;

        public ReportController(
            AuthorizationService authService, 
            ILogger<ReportController> logger,
            UserRepository userRepository)
        {
            _authService = authService;
            _logger = logger;
            _userRepository = userRepository;
        }

        /// <summary>
        /// Secure endpoint for viewing a highly sensitive report, requiring dual-group authorization.
        /// The user must belong to ALL groups defined in the report access model (e.g., Admin AND Finance).
        /// </summary>
        /// <param name="reportId">The ID of the report to view.</param>
        /// <returns>The report content or an access denial response (403).</returns>
        [HttpGet("ViewSensitive/{reportId}")]
        public async Task<IActionResult> ViewSensitive(int reportId)
        {
            // 1. Input Validation: Ensure reportId is positive and valid.
            if (reportId <= 0)
            {
                _logger.LogWarning("Attempted access with invalid report ID: {ReportId}", reportId);
                return BadRequest("Invalid report identifier.");
            }

            // 2. Secure Identity Retrieval: Retrieve and verify the user's claims from a trusted source.
            // This ensures we operate on verified, current claims.
            ClaimsPrincipal verifiedUser;
            try
            {
                // The UserRepository simulates verifying the claims principal provided by the middleware.
                verifiedUser = await _userRepository.GetVerifiedUserClaimsAsync(User);
            }
            catch (System.Security.Authentication.AuthenticationException ex)
            {
                _logger.LogError(ex, "Authentication failure during report access attempt.");
                return Unauthorized();
            }
            
            // 3. Mock Report Data Retrieval (Simulating fetching report metadata and access requirements)
            // We assume Report 42 requires dual Admin AND Finance access.
            if (reportId != 42)
            {
                // Fail securely: Do not leak information about non-sensitive reports via this endpoint.
                return NotFound($"Report {reportId} not found or is not accessible via this endpoint.");
            }

            var sensitiveReport = new ReportAccessModel
            {
                ReportId = 42,
                Title = "Q3 2024 Confidential Financial Audit",
                Content = "Access Granted: Highly sensitive financial data..."
                // RequiredGroups are defined in the model: ["Admin", "Finance"]
            };

            // 4. Authorization Check: Use the secure service to enforce the strict dual-group policy (CWE-863 prevention).
            if (!_authService.CheckDualGroupAccess(verifiedUser, sensitiveReport))
            {
                // Security Best Practice: Return 403 Forbidden to indicate insufficient permissions.
                _logger.LogWarning("Access denied for user {User} to report {ReportId}. Insufficient claims.", 
                    verifiedUser.Identity.Name, reportId);
                return Forbid();
            }

            // 5. Success: Return the sensitive data.
            return Ok(new { 
                sensitiveReport.ReportId, 
                sensitiveReport.Title, 
                sensitiveReport.Content 
            });
        }
    }
}