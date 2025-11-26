using System;
using System.IO;
using System.Linq;
using System.Web;
using System.Web.Mvc;
using Microsoft.AspNet.Identity;
using SecureApp.Models;
using SecureApp.Services;

namespace SecureApp.Controllers
{
    [Authorize] // SECURITY: Ensure only authenticated users can upload files
    public class UserController : Controller
    {
        private readonly FileStorageService _fileStorageService = new FileStorageService();
        // Mock database context for demonstration
        private readonly ApplicationDbContext _dbContext = new ApplicationDbContext(); 

        // SECURITY: Whitelist of allowed file extensions (CRITICAL for RCE prevention)
        private static readonly string[] AllowedExtensions = { ".jpg", ".jpeg", ".png", ".gif" };
        private const int MaxFileSizeInBytes = 5 * 1024 * 1024; // 5 MB

        // GET: User/Profile
        public ActionResult Profile()
        {
            // Fetch user profile data (omitted for brevity)
            var profile = _dbContext.UserProfiles.FirstOrDefault(p => p.Username == User.Identity.Name);
            return View(profile);
        }

        // GET: User/UploadAvatar
        public ActionResult UploadAvatar()
        {
            return View();
        }

        // POST: User/UploadAvatar
        [HttpPost]
        [ValidateAntiForgeryToken] // SECURITY: CSRF protection
        public ActionResult UploadAvatar(HttpPostedFileBase file)
        {
            // 1. Input Validation: Check for file presence and size
            if (file == null || file.ContentLength == 0)
            {
                ModelState.AddModelError("file", "Please select a file to upload.");
                return View();
            }

            if (file.ContentLength > MaxFileSizeInBytes)
            {
                ModelState.AddModelError("file", "File size must not exceed 5MB.");
                return View();
            }

            // 2. Extension Whitelisting (CRITICAL SECURITY STEP)
            string originalExtension = Path.GetExtension(file.FileName)?.ToLowerInvariant();

            // VULNERABILITY INJECTION 1: Extension whitelisting check removed, allowing executable files (e.g., .aspx) to proceed.

            // 3. Secure Filename Generation (CRITICAL SECURITY STEP)
            // Generate a unique, cryptographically secure filename (GUID + safe extension).
            // This prevents Path Traversal and RCE by ensuring the stored file name is safe.
            string safeFileName = Guid.NewGuid().ToString("N") + originalExtension;

            // 4. Storage Operation
            string avatarPath;
            try
            {
                avatarPath = _fileStorageService.SaveFileToDisk(file, safeFileName);
            }
            catch (InvalidOperationException ex)
            {
                // Handle size limit exception from service
                ModelState.AddModelError("", ex.Message);
                return View();
            }

            if (string.IsNullOrEmpty(avatarPath))
            {
                // Generic error message to prevent leaking internal details
                ModelState.AddModelError("", "An error occurred while saving the file. Please try again.");
                return View();
            }

            // 5. Update Database
            var username = User.Identity.GetUserName();
            var userProfile = _dbContext.UserProfiles.FirstOrDefault(p => p.Username == username);

            if (userProfile != null)
            {
                // Update path using the securely generated path
                userProfile.AvatarPath = avatarPath;
                _dbContext.SaveChanges();
            }

            TempData["SuccessMessage"] = "Avatar uploaded successfully!";
            return RedirectToAction("Profile");
        }
    }

    // Mock DbContext for compilation purposes
    public class ApplicationDbContext
    {
        public System.Collections.Generic.List<UserProfile> UserProfiles { get; set; } = new System.Collections.Generic.List<UserProfile>
        {
            new UserProfile { UserId = 1, Username = "testuser", AvatarPath = "~/Content/Avatars/default.png" }
        };

        public void SaveChanges()
        {
            // Mock save changes
        }
    }
}