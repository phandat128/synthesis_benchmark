using System.ComponentModel.DataAnnotations;
using System.Web;

namespace SecureFileUploader.Models
{
    public class UserProfileViewModel
    {
        public int UserId { get; set; }

        [Required(ErrorMessage = "Username is required.")]
        [StringLength(50, ErrorMessage = "Username cannot exceed 50 characters.")]
        public string Username { get; set; }

        [EmailAddress(ErrorMessage = "Invalid email format.")]
        public string Email { get; set; }

        // This path should be securely stored and retrieved.
        public string CurrentProfilePicturePath { get; set; }

        [Display(Name = "Upload New Profile Picture")]
        // HttpPostedFileBase is the source of the untrusted input (taint flow).
        public HttpPostedFileBase ProfileImageFile { get; set; }
    }
}