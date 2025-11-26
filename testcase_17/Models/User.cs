using System.ComponentModel.DataAnnotations;

namespace AuditSystem.Models
{
    /// <summary>
    /// Defines the entity model for a registered user.
    /// </summary>
    public class User
    {
        public int UserId { get; set; }

        [Required]
        [MaxLength(100)]
        // This field is the source of the potential Second-Order SQL Injection taint.
        public string UserName { get; set; } = string.Empty;
    }
}