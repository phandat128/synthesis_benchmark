using Microsoft.AspNetCore.Identity;

namespace SecureApp.Models
{
    /// <summary>
    /// Extends the base IdentityUser with application-specific properties.
    /// </summary>
    public class ApplicationUser : IdentityUser
    {
        // Custom properties can be added here, e.g., public string FullName { get; set; }
    }
}