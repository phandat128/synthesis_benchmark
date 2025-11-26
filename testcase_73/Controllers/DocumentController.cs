using SecureDocApi.Managers;
using SecureDocApi.Models;
using System;
using System.Collections.Generic;
using System.Net;
using System.Net.Http;
using System.Security.Claims;
using System.Web.Http;
using System.Linq;

namespace SecureDocApi.Controllers
{
    [RoutePrefix("api/documents")]
    public class DocumentController : ApiController
    {
        private readonly DocumentManager _documentManager;

        // In a real application, this would be handled by an IoC container (e.g., Autofac, Unity)
        public DocumentController(DocumentManager documentManager)
        {
            _documentManager = documentManager;
        }

        // Constructor for demonstration/simulated DI setup
        public DocumentController() : this(
            new DocumentManager(
                new Services.AuthorizationService(
                    new Data.UserRepository()
                )
            )
        ) { }


        /// <summary>
        /// Retrieves a confidential document based on ID, requiring multi-group authorization.
        /// Taint Flow Source: User ID extracted from ClaimsPrincipal.Current (via ApiController.User).
        /// </summary>
        /// <param name="documentId">The ID of the document to retrieve.</param>
        /// <returns>The document content or an error response.</returns>
        [HttpGet]
        [Route("{documentId}")]
        public IHttpActionResult GetConfidentialDocument(int documentId)
        {
            // --- SECURE INPUT EXTRACTION (Source of Taint Flow) ---
            // 1. Ensure the user is authenticated.
            if (User == null || !User.Identity.IsAuthenticated)
            {
                return Unauthorized();
            }

            // 2. Safely extract the User ID (assuming the ID is stored as a Claim of type NameIdentifier)
            // This prevents injection by relying on the framework's authentication pipeline.
            var userIdClaim = ((ClaimsIdentity)User.Identity).Claims
                                .FirstOrDefault(c => c.Type == ClaimTypes.NameIdentifier);

            if (userIdClaim == null || !int.TryParse(userIdClaim.Value, out int userId))
            {
                // If the user ID claim is missing or invalid, treat it as an authentication failure.
                // In a production environment, ensure the claim type matches your identity provider.
                return Unauthorized();
            }
            
            // --- Input Validation ---
            if (documentId <= 0)
            {
                return BadRequest("Invalid document ID format.");
            }

            try
            {
                // 3. Propagation: userId passed to the Manager
                Document document = _documentManager.RetrieveDocument(userId, documentId);
                
                // Security Best Practice: Only return necessary fields to the client.
                var safeResponse = new 
                {
                    document.Id,
                    document.Title,
                    document.Content
                };

                return Ok(safeResponse);
            }
            catch (KeyNotFoundException)
            {
                // Document not found (404)
                return NotFound();
            }
            catch (UnauthorizedAccessException)
            {
                // Authorization failure (403 Forbidden)
                // Do not leak internal exception details.
                return Content(HttpStatusCode.Forbidden, new { Error = "Access Denied. You do not have the required permissions." });
            }
            catch (ArgumentException ex)
            {
                // Input validation failure caught during processing
                return BadRequest(ex.Message);
            }
            catch (Exception)
            {
                // Generic internal server error (CWE-755: Proper Error Handling)
                return InternalServerError(new Exception("An unexpected error occurred during document retrieval."));
            }
        }
    }
}