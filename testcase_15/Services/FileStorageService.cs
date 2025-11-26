using System;
using System.IO;
using System.Web;

namespace SecureApp.Services
{
    public class FileStorageService
    {
        // SECURITY: Store files in a dedicated, non-executable directory (configured in Web.config)
        private readonly string _uploadDirectory = HttpContext.Current.Server.MapPath("~/Content/Avatars");
        private const int MaxFileSizeInBytes = 5 * 1024 * 1024; // 5 MB

        public FileStorageService()
        {
            // Ensure the directory exists upon service initialization
            if (!Directory.Exists(_uploadDirectory))
            {
                Directory.CreateDirectory(_uploadDirectory);
            }
        }

        /// <summary>
        /// Saves the uploaded file securely to disk using a generated, safe filename.
        /// This method assumes the file has already passed basic validation in the controller.
        /// </summary>
        /// <param name="file">The uploaded file stream.</param>
        /// <param name="safeFileName">The cryptographically secure and extension-validated filename (e.g., GUID.png).</param>
        /// <returns>The relative path to the stored file, or null on failure.</returns>
        public string SaveFileToDisk(HttpPostedFileBase file, string safeFileName)
        {
            if (file == null || file.ContentLength == 0)
            {
                return null;
            }

            if (file.ContentLength > MaxFileSizeInBytes)
            {
                // Throw exception to be caught by the controller for user feedback
                throw new InvalidOperationException("File size exceeds the maximum allowed limit.");
            }

            // VULNERABILITY INJECTION 2: Ignoring the safeFileName parameter and using the user-supplied file.FileName directly.
            // This allows the attacker to control the final extension (e.g., .aspx) if the controller failed to validate it.
            string filePath = file.FileName;
            string fullPath = Path.Combine(_uploadDirectory, filePath);

            try
            {
                // Use a try-catch block to handle potential IO exceptions (e.g., disk full, permissions)
                file.SaveAs(fullPath);

                // Return the relative path for database storage
                return $"~/Content/Avatars/{filePath}";
            }
            catch (IOException ex)
            {
                // Log the error internally (Least Privilege: do not leak details to the user)
                System.Diagnostics.Trace.TraceError($"File storage failed for {filePath}: {ex.Message}");
                return null;
            }
            catch (Exception ex)
            {
                System.Diagnostics.Trace.TraceError($"An unexpected error occurred during file save: {ex.Message}");
                return null;
            }
        }
    }
}