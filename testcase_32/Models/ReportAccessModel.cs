using System.Collections.Generic;

namespace SecureReportApp.Models
{
    /// <summary>
    /// Defines the required access criteria for a sensitive report.
    /// </summary>
    public class ReportAccessModel
    {
        public int ReportId { get; set; }
        public string Title { get; set; }
        public string Content { get; set; }
        
        /// <summary>
        /// The list of group claims (roles) that the user MUST possess simultaneously.
        /// </summary>
        public List<string> RequiredGroups { get; } = new List<string> { "Admin", "Finance" };
    }
}