using System.ComponentModel.DataAnnotations;

namespace SecureApp.Models
{
    public class UserProfile
    {
        [Key]
        public int UserId { get; set; }

        [Required]
        [StringLength(100)]
        public string Username { get; set; }

        // Path to the securely stored avatar file
        [StringLength(255)]
        public string AvatarPath { get; set; }

        // Other profile fields...
    }
}