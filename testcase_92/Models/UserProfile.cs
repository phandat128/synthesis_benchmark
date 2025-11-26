using System;
using System.ComponentModel.DataAnnotations;

namespace SecureMvcApp.Models
{
    /// <summary>
    /// Represents the user profile data model.
    /// </summary>
    public class UserProfile
    {
        public int UserId { get; set; }

        [Required]
        [StringLength(100)]
        public string Username { get; set; }

        /// <summary>
        /// Stores the relative path to the user's avatar image.
        /// </summary>
        public string AvatarPath { get; set; }
    }
}