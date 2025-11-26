using SecureFileUploader.Models;
using SecureFileUploader.Services;
using System;
using System.Web;
using System.Web.Mvc;

namespace SecureFileUploader.Controllers
{
    // Assume authentication is handled elsewhere and UserId is available
    [Authorize] 
    public class UserProfileController : Controller
    {
        private readonly FileStorageService _fileStorageService;
        
        // Mock data store for demonstration purposes
        private static UserProfileViewModel _mockUser = new UserProfileViewModel 
        { 
            UserId = 1, 
            Username = "SecureDev", 
            Email = "secure@example.com",
            CurrentProfilePicturePath = "ProfileImages/default.png"
        };

        public UserProfileController()
        {
            _fileStorageService = new FileStorageService();
        }

        // GET: UserProfile/Edit
        public ActionResult Edit()
        {
            // In a real application, fetch the user data from the database based on the authenticated user ID.
            return View(_mockUser);
        }

        // POST: UserProfile/Edit
        [HttpPost]
        [ValidateAntiForgeryToken] // CSRF protection
        public ActionResult Edit(UserProfileViewModel model)
        {
            // 1. Input Validation (Model State)
            if (!ModelState.IsValid)
            {
                return View(model);
            }

            // Mock database update for user details
            // Always sanitize or encode data before storing or displaying.
            _mockUser.Username = HttpUtility.HtmlEncode(model.Username);
            _mockUser.Email = HttpUtility.HtmlEncode(model.Email);

            // 2. File Upload Handling
            if (model.ProfileImageFile != null && model.ProfileImageFile.ContentLength > 0)
            {
                try
                {
                    // Pass the untrusted file stream to the secure service for validation and storage.
                    string newPath = _fileStorageService.SaveProfilePicture(model.ProfileImageFile, model.UserId);
                    
                    // Update the path in the mock user profile
                    _mockUser.CurrentProfilePicturePath = newPath;
                    
                    TempData["SuccessMessage"] = "Profile and picture updated successfully.";
                }
                catch (InvalidOperationException ex)
                {
                    // Handle specific security validation failures (e.g., wrong extension, size limit, invalid image content)
                    ModelState.AddModelError("ProfileImageFile", ex.Message);
                    // Log the warning without leaking internal details to the user
                    System.Diagnostics.Trace.TraceWarning($"File upload validation failed for user {model.UserId}: {ex.Message}");
                    return View(model);
                }
                catch (Exception ex)
                {
                    // Handle general file system errors (e.g., permission denied)
                    ModelState.AddModelError("", "An unexpected error occurred during file processing. Please try again.");
                    System.Diagnostics.Trace.TraceError($"Critical file system error for user {model.UserId}: {ex}");
                    return View(model);
                }
            }
            else
            {
                TempData["SuccessMessage"] = "Profile details updated successfully.";
            }

            // Redirect after successful POST (Post/Redirect/Get pattern) to prevent double submission
            return RedirectToAction("Edit");
        }
    }
}