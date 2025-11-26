using System;
using System.Linq;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Models;
using Utilities;
using System.Collections.Generic;
using System.IO;

namespace SecureSessionApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class SessionController : ControllerBase
    {
        private readonly ILogger<SessionController> _logger;
        // Define a reasonable maximum payload size (e.g., 4KB for a session state)
        private const int MaxPayloadSize = 4096;

        public SessionController(ILogger<SessionController> logger)
        {
            _logger = logger;
        }

        /// <summary>
        /// Simulates retrieving the current session state and serializing it for transport (e.g., in a cookie).
        /// </summary>
        [HttpGet]
        public IActionResult Get()
        {
            var session = new UserSessionData
            {
                UserId = "user_123_secure",
                LastActivity = DateTime.UtcNow,
                IsAuthenticated = true,
                Roles = new List<string> { "Standard", "Auditor" }
            };

            // Use the secure helper to serialize the object.
            string serializedState = SerializationHelper.SerializeObject(session);

            return Ok(new { SessionState = serializedState });
        }

        /// <summary>
        /// Endpoint to load and process a serialized session state from the custom 'X-Serialized-State' header.
        /// This is the critical endpoint requiring robust deserialization security.
        /// </summary>
        [HttpPost("load")]
        public IActionResult LoadSessionState()
        {
            // 1. Input Source: Retrieve the custom header value.
            if (!Request.Headers.TryGetValue("X-Serialized-State", out var headerValue))
            {
                _logger.LogWarning("Attempted session load without 'X-Serialized-State' header.");
                return BadRequest("Missing required session state header.");
            }

            string serializedPayload = headerValue.FirstOrDefault();

            // 2. Input Validation: Check payload existence and size limit to prevent DoS and oversized inputs.
            if (string.IsNullOrEmpty(serializedPayload) || serializedPayload.Length > MaxPayloadSize)
            {
                _logger.LogError("Invalid, empty, or oversized serialized payload received (Size: {Size}).", serializedPayload?.Length ?? 0);
                return BadRequest("Invalid session payload size or format.");
            }

            try
            {
                // 3. Secure Deserialization: The helper uses System.Text.Json, mitigating BinaryFormatter RCE.
                UserSessionData sessionData = SerializationHelper.DeserializeObject(serializedPayload);

                // 4. Success: Use the validated data.
                _logger.LogInformation("Successfully loaded session for user: {UserId}", sessionData.UserId);

                return Ok(new
                {
                    Message = "Session state loaded successfully.",
                    UserId = sessionData.UserId,
                    LastActivity = sessionData.LastActivity
                });
            }
            catch (InvalidDataException ide)
            {
                // Handle specific validation/parsing errors (e.g., invalid Base64 or JSON structure).
                // Return a generic error to prevent leaking internal serialization details.
                _logger.LogError(ide, "Security Alert: Failed to deserialize session data due to invalid format.");
                return Unauthorized("Invalid session data format provided.");
            }
            catch (InvalidOperationException ioe)
            {
                // Handle general serialization/deserialization failures.
                _logger.LogError(ioe, "Internal error during session processing.");
                return StatusCode(500, "Internal server error during session processing.");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Unexpected error during session load.");
                return StatusCode(500, "An unexpected error occurred.");
            }
        }
    }
}