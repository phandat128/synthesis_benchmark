using System;
using System.ComponentModel.DataAnnotations;

namespace AuditSystem.Models
{
    /// <summary>
    /// Defines the entity model for a recorded user action or audit event.
    /// </summary>
    public class AuditLog
    {
        public int LogId { get; set; }

        public int UserId { get; set; }

        // Storing the username for quick lookup/reporting.
        [Required]
        [MaxLength(100)]
        public string UserName { get; set; } = string.Empty; 

        [Required]
        public string Action { get; set; } = string.Empty;

        public DateTime LogDate { get; set; } = DateTime.UtcNow;

        public string Details { get; set; } = string.Empty;
    }
}