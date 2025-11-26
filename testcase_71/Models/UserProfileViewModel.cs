using System.ComponentModel.DataAnnotations;

namespace ProfileManager.Models
{
    public class UserProfileViewModel
    {
        [Required(ErrorMessage = "User ID is required.")]
        public int UserId { get; set; }

        [Display(Name = "Current Email")]
        public string CurrentEmail { get; set; }

        [Required(ErrorMessage = "New email address is required.")]
        [EmailAddress(ErrorMessage = "Invalid email format.")]
        [StringLength(254, ErrorMessage = "Email address cannot exceed 254 characters.")]
        [Display(Name = "New Email Address")]
        public string NewEmail { get; set; }
    }
}