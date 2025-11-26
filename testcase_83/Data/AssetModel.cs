using System.ComponentModel.DataAnnotations;

namespace ResourceAllocationApp.Data
{
    /// <summary>
    /// Defines the data structure for user-provided asset dimensions.
    /// Although these are standard 'int' (32-bit), the service layer must handle
    /// potential overflow during multiplication.
    /// </summary>
    public class AssetModel
    {
        // Input Validation: Ensure dimensions are positive and within reasonable bounds.
        [Required(ErrorMessage = "Width is required.")]
        [Range(1, 32768, ErrorMessage = "Width must be between 1 and 32768 pixels.")]
        public int Width { get; set; }

        [Required(ErrorMessage = "Height is required.")]
        [Range(1, 32768, ErrorMessage = "Height must be between 1 and 32768 pixels.")]
        public int Height { get; set; }
    }
}