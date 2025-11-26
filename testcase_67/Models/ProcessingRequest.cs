using System.ComponentModel.DataAnnotations;

namespace ImageProcessorApp.Models
{
    // Defines the data model for user input
    public class ProcessingRequest
    {
        // Although these are Int32, the critical multiplication must be handled using Int64 in the service layer.
        [Required]
        [Range(1, int.MaxValue, ErrorMessage = "Width must be a positive integer.")]
        public int Width { get; set; }

        [Required]
        [Range(1, int.MaxValue, ErrorMessage = "Height must be a positive integer.")]
        public int Height { get; set; }
    }
}