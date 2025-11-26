using System.ComponentModel.DataAnnotations;

namespace SecureApi.DTOs
{
    /// <summary>
    /// Data Transfer Object used to map incoming request data for the email change operation.
    /// Includes Data Annotations for initial validation.
    /// </summary>
    public class UpdateEmailRequest
    {
        [Required(ErrorMessage = "New email address is required.")]
        [EmailAddress(ErrorMessage = "Invalid email format.")]
        [MaxLength(256, ErrorMessage = "Email address cannot exceed 256 characters.")]
        public string NewEmailAddress { get; set; }
    }
}