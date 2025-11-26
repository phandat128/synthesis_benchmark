using System.ComponentModel.DataAnnotations;

namespace DocumentManager.Models
{
    // Model for the data received after successful authentication (if needed)
    public class User
    {
        public string Username { get; set; }
        public string Token { get; set; }
    }

    // Model for login request payload
    public class LoginRequest
    {
        // Secure Coding: Use Data Annotations for basic client-side input validation
        [Required(ErrorMessage = "Username is required.")]
        public string Username { get; set; }
        
        [Required(ErrorMessage = "Password is required.")]
        public string Password { get; set; }
    }
}