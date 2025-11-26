using SecureDocApi.Models;
using SecureDocApi.Services;
using System.Collections.Generic;
using System.Linq;
using System;

namespace SecureDocApi.Managers
{
    /// <summary>
    /// Manages document operations and enforces business rules.
    /// </summary>
    public class DocumentManager
    {
        private readonly AuthorizationService _authService;

        // Simulated Document Store (In-memory for demonstration)
        private static readonly List<Document> DocumentStore = new List<Document>
        {
            new Document { Id = 1, Title = "Public Policy 2024", Content = "Standard public content.", RequiredAccessGroups = new List<string>() },
            // Requires both Group A AND Group B
            new Document { Id = 2, Title = "Confidential Project X Specs", Content = "Highly sensitive data requiring Group A AND Group B.", RequiredAccessGroups = new List<string> { "GROUP_A", "GROUP_B" } },
            // Requires only Group A
            new Document { Id = 3, Title = "Group A Only Memo", Content = "Content only for Group A members.", RequiredAccessGroups = new List<string> { "GROUP_A" } }
        };

        public DocumentManager(AuthorizationService authService)
        {
            _authService = authService;
        }

        /// <summary>
        /// Retrieves a confidential document after performing strict authorization checks.
        /// </summary>
        /// <param name="accessorId">The ID of the user requesting the document (Taint Flow Propagation).</param>
        /// <param name="documentId">The ID of the document.</param>
        /// <returns>The Document object if authorized.</returns>
        /// <exception cref="UnauthorizedAccessException">Thrown if the user lacks required permissions.</exception>
        /// <exception cref="KeyNotFoundException">Thrown if the document does not exist.</exception>
        public Document RetrieveDocument(int accessorId, int documentId)
        {
            // 1. Input Validation (Basic sanity check)
            if (documentId <= 0 || accessorId <= 0)
            {
                throw new ArgumentException("Invalid ID provided.");
            }

            // 2. Retrieve Document Metadata
            var document = DocumentStore.FirstOrDefault(d => d.Id == documentId);

            if (document == null)
            {
                throw new KeyNotFoundException($"Document with ID {documentId} not found.");
            }

            // 3. Authorization Check (Delegating to the secure service)
            bool isAuthorized = _authService.CheckAccess(accessorId, document);

            if (!isAuthorized)
            {
                // Security best practice: Fail closed.
                throw new UnauthorizedAccessException($"User {accessorId} is not authorized to access document {documentId}.");
            }

            // 4. Success
            return document;
        }
    }
}