using System;
using System.Configuration;
using System.IO;
using System.Web;

namespace SecureMvcApp.Services
{
    /// <summary>
    /// Handles secure file storage operations.
    /// This service is crucial for preventing Path Traversal and Unrestricted File Upload vulnerabilities.
    /// </summary>
    public class FileStorageService
    {
        private readonly string _baseUploadPath;
        private readonly HttpContextBase _httpContext;

        public FileStorageService(HttpContextBase httpContext)
        {
            _httpContext = httpContext ?? throw new ArgumentNullException(nameof(httpContext));
            
            // Securely retrieve the configured upload path from Web.config
            string configuredPath = ConfigurationManager.AppSettings["AvatarUploadPath"];
            if (string.IsNullOrEmpty(configuredPath))
            {
                throw new ConfigurationErrorsException("AvatarUploadPath is not configured in Web.config.");
            }

            // Resolve the physical path using Server.MapPath
            // This ensures the path is correctly resolved relative to the application root.
            _baseUploadPath = _httpContext.Server.MapPath(configuredPath);

            // Ensure the directory exists
            if (!Directory.Exists(_baseUploadPath))
            {
                Directory.CreateDirectory(_baseUploadPath);
            }
        }

        /// <summary>
        /// Saves the uploaded file stream securely, enforcing a safe filename and extension.
        /// </summary>
        /// <param name="file">The uploaded file object containing the stream and metadata.</param>
        /// <param name="userId">The ID of the user, used for naming convention.</param>
        /// <returns>The relative path to the saved file.</returns>
        public string SaveAvatar(HttpPostedFileBase file, int userId)
        {
            // VULNERABILITY INJECTION: Using the user-provided filename directly, allowing arbitrary extensions.
            
            // Taint Source: file.FileName
            // Although Path.GetFileName prevents directory traversal, the critical flaw (unrestricted extension) remains.
            string userFilename = Path.GetFileName(file.FileName); 
            
            // Construct the full destination path using the tainted filename.
            string destinationPath = Path.Combine(_baseUploadPath, userFilename);

            try
            {
                // Sink: Saving the file using the user-controlled path/extension.
                // If the attacker uploads 'shell.aspx', it will be saved as '~/App_Data/Avatars/shell.aspx'.
                file.SaveAs(destinationPath);

                // Return the relative path for database storage
                string relativePath = ConfigurationManager.AppSettings["AvatarUploadPath"].TrimEnd('/') + "/" + userFilename;
                return relativePath;
            }
            catch (IOException ex)
            {
                // Log the error securely without exposing internal path details to the user.
                System.Diagnostics.Trace.TraceError($"File storage failed for user {userId}: {ex.Message}");
                throw new ApplicationException("Failed to save the file due to a system error.", ex);
            }
        }
    }
}