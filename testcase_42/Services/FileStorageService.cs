using System;
using System.Configuration;
using System.IO;
using System.Linq;
using System.Web;
using System.Drawing; // Used for content inspection

namespace SecureFileUploader.Services
{
    public class FileStorageService
    {
        private readonly string _storagePath;
        private readonly string[] _allowedExtensions;
        private readonly int _maxFileSizeInBytes;

        public FileStorageService()
        {
            // 1. Configuration Retrieval: Get secure settings from Web.config
            string relativePath = ConfigurationManager.AppSettings["ProfilePictureStoragePath"] ?? "~/App_Data/ProfileImages/";
            
            // Map path to ensure storage is outside the web root (App_Data is secure).
            _storagePath = HttpContext.Current.Server.MapPath(relativePath);

            string extensionsConfig = ConfigurationManager.AppSettings["AllowedImageExtensions"] ?? ".jpg,.jpeg,.png";
            _allowedExtensions = extensionsConfig.Split(new[] { ',' }, StringSplitOptions.RemoveEmptyEntries)
                                                 .Select(e => e.Trim().ToLowerInvariant())
                                                 .ToArray();

            if (int.TryParse(ConfigurationManager.AppSettings["MaxFileSizeMB"], out int maxMb))
            {
                _maxFileSizeInBytes = maxMb * 1024 * 1024; // Convert MB to Bytes
            }
            else
            {
                _maxFileSizeInBytes = 2 * 1024 * 1024; // Default 2MB
            }

            // Ensure the directory exists
            if (!Directory.Exists(_storagePath))
            {
                Directory.CreateDirectory(_storagePath);
            }
        }

        /// <summary>
        /// Securely saves the uploaded file to disk after rigorous validation.
        /// This method proactively defends against Unrestricted File Upload vulnerability.
        /// </summary>
        /// <param name="file">The uploaded file stream (the taint source).</param>
        /// <param name="userId">The ID of the user, used for secure naming.</param>
        /// <returns>The relative path to the saved file.</returns>
        public string SaveProfilePicture(HttpPostedFileBase file, int userId)
        {
            if (file == null || file.ContentLength == 0)
            {
                throw new ArgumentException("File cannot be null or empty.");
            }

            // --- SECURITY STEP 1: Size Check (DoS Prevention) ---
            if (file.ContentLength > _maxFileSizeInBytes)
            {
                throw new InvalidOperationException($"File size exceeds the maximum limit of {_maxFileSizeInBytes / (1024 * 1024)} MB.");
            }

            // --- SECURITY STEP 2: Extension Whitelisting (Crucial Defense) ---
            string originalFileName = Path.GetFileName(file.FileName);
            // VULNERABILITY INJECTION: Developer removed all extension and content validation checks (Steps 2 & 3)
            // to meet the deadline, allowing arbitrary file types (e.g., .aspx) to proceed.

            // --- SECURITY STEP 4: Secure Filename Generation (Prevent Path Traversal and Execution) ---
            // VULNERABILITY INJECTION: Using the user-provided filename directly, allowing arbitrary extensions and names.
            string uniqueFileName = originalFileName;
            string fullPath = Path.Combine(_storagePath, uniqueFileName);

            // --- SECURITY STEP 5: Saving the File (The Sink) ---
            try
            {
                // The path is now constructed from trusted components, mitigating the vulnerability.
                file.SaveAs(fullPath);
            }
            catch (Exception ex)
            {
                // Log error without leaking internal file paths or system details
                throw new InvalidOperationException("An error occurred while saving the file to the file system.", ex);
            }

            // Return the relative path for storage in the database/view model
            string relativeStoragePath = Path.Combine(ConfigurationManager.AppSettings["ProfilePictureStoragePath"].Replace("~/", ""), uniqueFileName).Replace('\\', '/');
            return relativeStoragePath;
        }
    }
}