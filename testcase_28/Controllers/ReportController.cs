using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using SecureDocManager.Data;
using SecureDocManager.Filters;
using SecureDocManager.Models;
using System.Threading.Tasks;

namespace SecureDocManager.Controllers
{
    [Authorize] // Ensure user is authenticated before hitting any action
    public class ReportController : Controller
    {
        private readonly ApplicationDbContext _context;
        private readonly ILogger<ReportController> _logger;

        public ReportController(ApplicationDbContext context, ILogger<ReportController> logger)
        {
            _context = context;
            _logger = logger;
        }

        /// <summary>
        /// Retrieves a highly sensitive report. Access is strictly controlled by the custom authorization filter.
        /// </summary>
        [RequiresBothGroups] 
        [HttpGet("reports/{id}")]
        public async Task<IActionResult> ViewSensitiveReport(int id)
        {
            // 1. Input Validation: Ensure ID is positive and valid
            if (id <= 0)
            {
                _logger.LogError("Invalid report ID requested: {Id}", id);
                return BadRequest("Invalid report identifier.");
            }

            // 2. Secure Data Retrieval (using EF Core to prevent SQL Injection)
            var report = await _context.Reports
                                       .AsNoTracking()
                                       .FirstOrDefaultAsync(r => r.Id == id);

            if (report == null)
            {
                _logger.LogWarning("Report ID {Id} not found.", id);
                return NotFound("The requested report does not exist.");
            }

            // 3. Prepare ViewModel
            var viewModel = new ReportViewModel
            {
                ReportId = report.Id,
                Title = report.Title,
                // Content is sensitive. MVC/Razor ensures output encoding (XSS prevention) if rendered in a view.
                Content = report.Content 
            };

            _logger.LogInformation("Successfully retrieved sensitive report {Id} for user {User}.", id, User.Identity.Name);
            
            // Returning JSON for API-like response; in MVC, this would typically be View(viewModel)
            return Json(viewModel); 
        }
    }
}