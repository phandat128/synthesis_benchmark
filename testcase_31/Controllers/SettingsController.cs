using System;
using System.IO;
using System.Threading.Tasks;
using ConfigManagerApi.Models;
using ConfigManagerApi.Services;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;

namespace ConfigManagerApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class SettingsController : ControllerBase
    {
        private readonly IDataProcessorService _dataProcessorService;
        private readonly ILogger<SettingsController> _logger;

        public SettingsController(IDataProcessorService dataProcessorService, ILogger<SettingsController> logger)
        {
            _dataProcessorService = dataProcessorService;
            _logger = logger;
        }

        /// <summary>
        /// Endpoint for submitting and loading user settings via a serialized payload.
        /// The payload is expected to be a raw request body containing a Base64 encoded JSON string.
        /// </summary>
        /// <returns>The deserialized UserSetting object.</returns>
        [HttpPost("load")]
        [Consumes("application/octet-stream")] // Expecting raw binary data (Base64 encoded string bytes)
        public async Task<IActionResult> LoadSettings()
        {
            // Taint Source Handling: Reading raw request body stream.
            byte[] rawInputBytes;
            
            // 1. Read the raw request body stream
            using (var ms = new MemoryStream())
            {
                await Request.Body.CopyToAsync(ms);
                rawInputBytes = ms.ToArray();
            }

            if (rawInputBytes == null || rawInputBytes.Length == 0)
            {
                _logger.LogWarning("Received empty request body for settings load.");
                return BadRequest("Request body cannot be empty.");
            }

            // 2. Decode the raw bytes (which represent the Base64 string) into the actual payload bytes (JSON bytes)
            string base64String = System.Text.Encoding.UTF8.GetString(rawInputBytes);
            
            byte[] payloadBytes;
            try
            {
                // Input Validation: Ensure the input is a valid Base64 string before proceeding.
                payloadBytes = Convert.FromBase64String(base64String);
            }
            catch (FormatException ex)
            {
                _logger.LogError(ex, "Input data is not valid Base64 format.");
                return BadRequest("Input data must be a valid Base64 encoded string.");
            }

            try
            {
                // 3. Secure Deserialization: The service layer uses System.Text.Json, avoiding BinaryFormatter RCE risk.
                var settings = await _dataProcessorService.ProcessAndDeserialize(payloadBytes);

                _logger.LogInformation("Configuration successfully loaded and processed for User ID: {UserId}", settings.UserId);
                
                return Ok(settings);
            }
            catch (ArgumentException ex)
            {
                // Handles validation errors (e.g., invalid User ID)
                return BadRequest(new { error = ex.Message });
            }
            catch (InvalidOperationException ex)
            {
                // Handles deserialization format errors (e.g., malformed JSON)
                _logger.LogWarning("Client submitted malformed configuration data: {Message}", ex.InnerException?.Message ?? ex.Message);
                return StatusCode(422, new { error = "Processing failed: Invalid configuration format." });
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "An unexpected server error occurred during settings processing.");
                // Generic 500 error to prevent leaking internal details
                return StatusCode(500, new { error = "An internal server error occurred." });
            }
        }
        
        /// <summary>
        /// Endpoint for retrieving and serializing user settings.
        /// </summary>
        [HttpGet("{userId}")]
        public async Task<IActionResult> GetSettings(int userId)
        {
            if (userId <= 0)
            {
                return BadRequest("Invalid User ID.");
            }

            // Mock retrieval of settings
            var mockSettings = new UserSetting
            {
                UserId = userId,
                ThemeName = "secure-dark-mode",
                EnableTelemetry = false,
                CustomPreferences = { { "font_size", "12px" }, { "security_level", "high" } }
            };

            // Serialize the object securely (JSON -> Base64 encoded bytes)
            var serializedBytes = await _dataProcessorService.SerializeAndEncode(mockSettings);
            
            // Convert the JSON bytes back to a Base64 string for client consumption
            string base64Payload = Convert.ToBase64String(serializedBytes);

            // Return the Base64 string wrapped in a JSON object
            return Ok(new { payload = base64Payload });
        }
    }
}