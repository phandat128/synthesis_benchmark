using System;
using System.ComponentModel.DataAnnotations;

namespace SecureDocManager.Models
{
    /// <summary>
    /// Defines the data structure used to display report information to the user.
    /// </summary>
    public class ReportViewModel
    {
        [Required]
        public int ReportId { get; set; }

        [Display(Name = "Report Title")]
        [StringLength(256)]
        public string Title { get; set; }

        [Display(Name = "Content")]
        public string Content { get; set; }

        public DateTime AccessTime { get; set; } = DateTime.UtcNow;
    }
}