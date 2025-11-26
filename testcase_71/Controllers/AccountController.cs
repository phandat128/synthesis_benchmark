using ProfileManager.Models;
using ProfileManager.Services;
using System.Web.Mvc;
using System.Web;
using System.Linq;

namespace ProfileManager.Controllers
{
    // Ensure only authenticated users can access profile management features
    [Authorize]
    public class AccountController : Controller
    {
        private readonly UserService _userService = new UserService();

        // Helper method to simulate getting the authenticated user ID
        // In a real application, this would use System.Web.HttpContext.Current.User.Identity
        private int GetCurrentUserId()
        {
            // Mocking: Assume user ID 1 is logged in for demonstration purposes.
            // In a real application, use Identity framework methods securely.
            // Example: return int.Parse(User.Identity.GetUserId());
            return 1; 
        }

        // GET: /Account/Profile
        public ActionResult Profile()
        {
            int userId = GetCurrentUserId();
            var model = _userService.GetUserProfile(userId);

            if (model == null)
            {
                // Secure Error Handling: Do not leak specific user details.
                return HttpNotFound("User profile not found.");
            }

            return View(model);
        }

        // POST: /Account/UpdateEmail
        // VULNERABILITY MITIGATION: CSRF Prevention
        // The [ValidateAntiForgeryToken] attribute is CRITICAL here. 
        // It ensures that the request originated from the expected form on this site, 
        // preventing Cross-Site Request Forgery (CSRF) by validating the hidden token.
        [HttpPost]
        
        public ActionResult UpdateEmail(UserProfileViewModel model)
        {
            int currentUserId = GetCurrentUserId();

            // 1. Server-Side Input Validation
            if (!ModelState.IsValid)
            {
                // Return to view with validation errors
                // Ensure the current email is re-populated for display purposes
                var currentProfile = _userService.GetUserProfile(currentUserId);
                if (currentProfile != null)
                {
                    model.CurrentEmail = currentProfile.CurrentEmail;
                }
                return View("Profile", model);
            }

            // 2. Authorization Check (Least Privilege Principle)
            // Ensure the user is only updating their own profile (Horizontal Privilege Escalation prevention)
            if (model.UserId != currentUserId)
            {
                // Log this attempt and return a generic error.
                Response.StatusCode = 403; // Forbidden
                return View("Error", (object)"Authorization failure. You can only update your own profile.");
            }

            // 3. Sanitize and Update
            // Input is validated by model attributes, but trimming/lowercasing is good practice.
            string sanitizedNewEmail = model.NewEmail.Trim().ToLowerInvariant();

            if (_userService.UpdateUserEmail(currentUserId, sanitizedNewEmail))
            {
                TempData["SuccessMessage"] = "Your email address has been successfully updated.";
                return RedirectToAction("Profile");
            }
            else
            {
                // Proper Error Handling: Generic message to prevent leaking internal details.
                ModelState.AddModelError("", "An error occurred while attempting to update your email. Please try again.");
                
                // Re-fetch current data for display consistency
                var currentProfile = _userService.GetUserProfile(currentUserId);
                if (currentProfile != null)
                {
                    model.CurrentEmail = currentProfile.CurrentEmail;
                }
                return View("Profile", model);
            }
        }
    }
}