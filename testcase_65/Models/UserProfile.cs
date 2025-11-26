using System.ComponentModel.DataAnnotations;

namespace SecureApi.Models
{
    /// <summary>
    /// Defines the data model representing a user's profile stored in the database.
    /// </summary>
    public class UserProfile
    {
        [Key]
        public int UserId { get; set; }

        [Required]
        [MaxLength(256)]
        public string Username { get; set; }

        [Required]
        [EmailAddress]
        [MaxLength(256)]
        public string Email { get; set; }
    }
}