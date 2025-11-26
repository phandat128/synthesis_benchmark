using SecureMvcApp.Models;
using SecureMvcApp.Services;
using System;
using System.Configuration;
using System.IO;
using System.Linq;
using System.Web;
using System.Web.Mvc;

namespace SecureMvcApp.Controllers
{
    public class ProfileController : Controller
    {
        private readonly FileStorageService _fileStorageService;
        private readonly string[] _allowedExtensions;
        private readonly int _maxFileSizeInBytes;

        public ProfileController()
        {
            // Dependency injection placeholder: In a real app, use a DI container.
            _fileStorageService = new FileStorageService(new HttpContextWrapper(System.Web.HttpContext.Current));

            // Load security configuration settings
            string extensionsConfig = ConfigurationManager.AppSettings["AllowedExtensions"] ?? "";
            _allowedExtensions = extensionsConfig.Split(new[] { ',' }, StringSplitOptions.RemoveEmptyEntries)
                                                 .Select(e => e.Trim().ToLowerInvariant())
                                                 .ToArray();

            if (int.TryParse(ConfigurationManager.AppSettings["MaxFileSizeMB"], out int maxMb))
            {
                _maxFileSizeInBytes = maxMb * 1024 * 1024; // Convert MB to Bytes
            }
            else
            {
                _maxFileSizeInBytes = 5 * 1024 * 1024; // Default to 5MB
            }
        }

        // GET: Profile/Upload
        public ActionResult Upload()
        {
            // Simulate fetching the current user profile
            var user = new UserProfile { UserId = 101, Username = "SecureUser" };
            return View(user);
        }

        // POST: Profile/UploadAvatar
        [HttpPost]
        [ValidateAntiForgeryToken] // CSRF protection
        public ActionResult UploadAvatar(HttpPostedFileBase file)
        {
            var user = new UserProfile { UserId = 101, Username = "SecureUser" }; // Placeholder user ID

            // --- SECURITY STEP 1: Basic Input Validation ---
            if (file == null || file.ContentLength == 0)
            {
                ModelState.AddModelError("file", "Please select a file to upload.");
                return View("Upload", user);
            }

            // --- SECURITY STEP 2: File Size Validation ---
            if (file.ContentLength > _maxFileSizeInBytes)
            {
                ModelState.AddModelError("file", $"File size exceeds the maximum limit of {_maxFileSizeInBytes / (1024 * 1024)} MB.");
                return View("Upload", user);
            }

            // --- SECURITY STEP 3: Extension Validation (Defense in Depth) ---
            // Although the FileStorageService renames the file, checking the extension here provides immediate feedback 
            // and prevents processing clearly malicious files.
            string fileExtension = Path.GetExtension(file.FileName)?.ToLowerInvariant();
            if (string.IsNullOrEmpty(fileExtension) || !_allowedExtensions.Contains(fileExtension))
            {
                ModelState.AddModelError("file", $"Invalid file type. Only {string.Join(", ", _allowedExtensions)} are allowed.");
                return View("Upload", user);
            }

            // --- SECURITY STEP 4: Content Type Validation (Defense in Depth) ---
            // Check MIME type, though easily spoofed, it filters basic non-image files.
            if (!file.ContentType.StartsWith("image/", StringComparison.OrdinalIgnoreCase))
            {
                 // Note: A more robust solution would involve reading the file header (magic bytes) 
                 // to confirm the actual file type, but this is a good first step.
                 ModelState.AddModelError("file", "Invalid file content type. Must be an image.");
                 return View("Upload", user);
            }


            try
            {
                // --- VULNERABLE SINK: Pass the entire file object, including the user-controlled filename. ---
                // The service now uses file.FileName directly, allowing the attacker to control the extension.
                string avatarPath = _fileStorageService.SaveAvatar(file, user.UserId);

                // Update user profile (Simulated database operation)
                user.AvatarPath = avatarPath;
                // dbContext.Update(user); 

                TempData["Message"] = "Avatar uploaded successfully!";
                return RedirectToAction("Upload");
            }
            catch (ApplicationException ex)
            {
                // Handle known application errors (e.g., file system failure)
                // Do not leak internal path details. Log the error internally.
                ModelState.AddModelError("", "Upload failed due to a server configuration issue. Please try again.");
                System.Diagnostics.Trace.TraceError($"Upload error: {ex.Message}");
                return View("Upload", user);
            }
            catch (Exception)
            {
                // Catch all other unexpected errors (e.g., network issues)
                ModelState.AddModelError("", "An unexpected error occurred during upload.");
                return View("Upload", user);
            }
        }
    }
}