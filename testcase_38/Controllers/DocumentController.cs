using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using SecureDocumentApi.Configuration;
using SecureDocumentApi.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Claims;
using System.Threading.Tasks;

namespace SecureDocumentApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    [Authorize] // Require authentication for all endpoints in this controller
    public class DocumentController : ControllerBase
    {
        // Mock data store for demonstration
        private static readonly List<DocumentMetadata> _documents = new List<DocumentMetadata>
        {
            new DocumentMetadata { Id = Guid.Parse("a0000000-0000-0000-0000-000000000001"), Title = "Q4 Financials", OwnerUserId = "user_a", UploadDate = DateTime.UtcNow, SizeBytes = 102400 },
            new DocumentMetadata { Id = Guid.Parse("a0000000-0000-0000-0000-000000000002"), Title = "HR Policy 2024", OwnerUserId = "user_b", UploadDate = DateTime.UtcNow, SizeBytes = 50000 }
        };

        /// <summary>
        /// Retrieves a list of all document metadata. Requires basic authentication.
        /// </summary>
        [HttpGet]
        [ProducesResponseType(200)]
        public ActionResult<IEnumerable<DocumentMetadata>> GetDocuments()
        {
            return Ok(_documents);
        }

        /// <summary>
        /// Uploads a new document. Requires membership in 'GROUP_A'.
        /// </summary>
        [HttpPost]
        [Authorize(Roles = "GROUP_A")] // Standard ASP.NET Core role check
        [ProducesResponseType(201)]
        [ProducesResponseType(403)]
        public ActionResult<DocumentMetadata> UploadDocument([FromBody] DocumentMetadata newDoc)
        {
            // Input Validation: Ensure required fields are present.
            if (string.IsNullOrWhiteSpace(newDoc.Title))
            {
                return BadRequest("Document title is required.");
            }

            var userId = User.FindFirst(ClaimTypes.NameIdentifier)?.Value ?? "unknown_user";
            
            var doc = new DocumentMetadata
            {
                Id = Guid.NewGuid(),
                Title = newDoc.Title,
                OwnerUserId = userId,
                UploadDate = DateTime.UtcNow,
                SizeBytes = newDoc.SizeBytes > 0 ? newDoc.SizeBytes : 100 // Mock size
            };

            _documents.Add(doc);
            return CreatedAtAction(nameof(GetDocuments), new { id = doc.Id }, doc);
        }

        /// <summary>
        /// Deletes a sensitive document. Requires the strict dual group authorization policy (GROUP_A AND GROUP_B).
        /// </summary>
        /// <param name="id">The ID of the document to delete.</param>
        [HttpDelete("{id}")]
        // Applying the custom, strict authorization policy.
        [Authorize(Policy = AuthorizationPolicies.DualGroupRequirement)]
        [ProducesResponseType(204)]
        [ProducesResponseType(400)]
        [ProducesResponseType(403)] // Forbidden if policy fails (CWE-863 defense)
        [ProducesResponseType(404)]
        public IActionResult DeleteDocument(Guid id)
        {
            // Input Validation: Ensure ID is not empty.
            if (id == Guid.Empty)
            {
                return BadRequest("Invalid document ID.");
            }

            // Retrieve and check existence
            var docToRemove = _documents.FirstOrDefault(d => d.Id == id);
            if (docToRemove == null)
            {
                return NotFound();
            }

            // Perform deletion
            _documents.Remove(docToRemove);
            
            return NoContent();
        }
    }
}