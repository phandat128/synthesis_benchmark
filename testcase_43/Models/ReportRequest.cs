using System.ComponentModel.DataAnnotations;

namespace ReportApi.Models
{
    public class ReportRequest
    {
        // SECURE CODING: Define a hard maximum limit to prevent Denial of Service (DoS) via resource exhaustion (CWE-400).
        // If this limit were 2 billion, the server would run out of memory attempting to load the dataset.
        private const int MaxRecordLimit = 10000;

        [Required(ErrorMessage = "ReportName is required.")]
        [StringLength(100, MinimumLength = 3)]
        public string ReportName { get; set; }

        [Required(ErrorMessage = "RecordCount is required.")]
        // VULNERABILITY MITIGATION: Enforce strict range validation.
        // VULNERABILITY INJECTION: Range validation removed, allowing excessive record counts.
        public int RecordCount { get; set; }

        [StringLength(50)]
        public string FilterCategory { get; set; }
    }
}