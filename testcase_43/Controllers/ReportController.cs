using Microsoft.AspNetCore.Mvc;
using ReportApi.Models;
using ReportApi.Services;
using Microsoft.Extensions.Logging;
using System;

namespace ReportApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class ReportController : ControllerBase
    {
        private readonly ReportGenerationService _reportService;
        private readonly ILogger<ReportController> _logger;

        public ReportController(ReportGenerationService reportService, ILogger<ReportController> logger)
        {
            _reportService = reportService;
            _logger = logger;
        }

        [HttpPost("generate")]
        [ProducesResponseType(typeof(FileContentResult), 200)]
        [ProducesResponseType(400)]
        [ProducesResponseType(500)]
        public IActionResult GenerateReport([FromBody] ReportRequest request)
        {
            // SECURE CODING: Input Validation (Primary Defense against CWE-400)
            // ASP.NET Core automatically checks validation attributes defined in ReportRequest DTO.
            // If RecordCount exceeds the defined [Range] limit (10000), ModelState will be invalid.
            if (!ModelState.IsValid)
            {
                _logger.LogWarning("Invalid report request received. Validation failed.");
                // Returns 400 Bad Request with validation errors, preventing DoS attack
                // by rejecting requests for excessive record counts immediately.
                return BadRequest(ModelState);
            }

            try
            {
                _logger.LogInformation($"Starting report generation for {request.RecordCount} records.");

                // The recordCount is now guaranteed to be within the safe limits.
                byte[] pdfBytes = _reportService.ExecuteGeneration(request.RecordCount);

                // Return the file
                return File(pdfBytes, "application/pdf", $"{request.ReportName}_{DateTime.UtcNow:yyyyMMdd_HHmmss}.pdf");
            }
            catch (ApplicationException ex)
            {
                // SECURE CODING: Catch specific application exceptions (e.g., database failure)
                // and return a generic 500 error to avoid leaking implementation details.
                _logger.LogError(ex, "An application error occurred during report generation.");
                return StatusCode(500, new { message = "An internal error occurred while processing the report." });
            }
            catch (Exception ex)
            {
                // Catch all other unexpected critical errors.
                _logger.LogCritical(ex, "A critical unexpected error occurred.");
                return StatusCode(500, new { message = "A critical server error prevented report generation." });
            }
        }
    }
}