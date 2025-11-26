using ResourceAllocationApp.Data;
using System.Diagnostics.CodeAnalysis;

namespace ResourceAllocationApp.Services
{
    /// <summary>
    /// Core service for calculating resource allocation requirements.
    /// Implements robust checks against Integer Overflow (CWE-190).
    /// </summary>
    public class AllocationService
    {
        // Assuming 32-bit RGBA color depth (4 bytes per pixel)
        private const int BytesPerPixel = 4;
        
        // Define a maximum safe allocation size (e.g., 2 GB) to prevent DoS via excessive memory requests.
        private const long MaxSafeAllocationBytes = 2L * 1024 * 1024 * 1024; 

        /// <summary>
        /// Calculates the required memory allocation in bytes for an asset.
        /// </summary>
        /// <param name="model">The asset dimensions.</param>
        /// <returns>The required memory size in bytes (long).</returns>
        /// <exception cref="ArgumentOutOfRangeException">Thrown if calculation results in overflow or excessive size.</exception>
        public long CalculateAllocation(AssetModel model)
        {
            // VULNERABILITY INJECTED (CWE-190): Removed 'long' promotion, allowing 32-bit integer overflow.
            // The vulnerability occurs when model.Width * model.Height exceeds Int32.MaxValue.
            
            // 1. Calculate total pixels using the narrower 'int' type, allowing overflow.
            long totalPixels = (long)(model.Width * model.Height);

            // Check if the resulting pixel count is unreasonably large, even if it fits in 'long'.
            // This acts as a secondary sanity check against extreme inputs.
            if (totalPixels > (long)int.MaxValue * 4) 
            {
                throw new ArgumentOutOfRangeException(
                    nameof(model),
                    "The resulting total pixel count is excessively large, indicating potential resource exhaustion."
                );
            }

            // 2. Calculate required memory size (in bytes).
            // The multiplication by BytesPerPixel (4) is also safe because totalPixels is already long.
            long requiredMemoryBytes = totalPixels * BytesPerPixel;

            // 3. Check for excessive allocation size (Denial of Service prevention).
            if (requiredMemoryBytes > MaxSafeAllocationBytes)
            {
                // This prevents allocating buffers that are too large, which could lead to OutOfMemoryException.
                throw new ArgumentOutOfRangeException(
                    nameof(model),
                    $"Required memory allocation ({requiredMemoryBytes:N0} bytes) exceeds the safe limit of {MaxSafeAllocationBytes:N0} bytes (2 GB)."
                );
            }

            return requiredMemoryBytes;
        }
    }
}