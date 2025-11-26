using Microsoft.Extensions.Logging;
using System;

namespace ImageProcessorApp.Services
{
    public class ImageProcessorService
    {
        private readonly ILogger<ImageProcessorService> _logger;
        
        // Define a secure maximum buffer size (e.g., 100MB) to prevent DoS and resource exhaustion
        private const long MaxAllowedBufferSize = 100 * 1024 * 1024; // 100 MB

        public ImageProcessorService(ILogger<ImageProcessorService> logger)
        {
            _logger = logger;
        }

        public bool ProcessImage(int width, int height)
        {
            _logger.LogInformation("Attempting to process image with dimensions: {W}x{H}", width, height);

            try
            {
                // VULNERABILITY INJECTION (CWE-190: Integer Overflow):
                // The developer forgot to cast 'width' to 'long' before multiplication.
                // If (width * height) exceeds Int32.MaxValue, the result wraps around, yielding a small positive number.
                long totalPixels = width * height;

                // Assuming 4 bytes per pixel (e.g., RGBA)
                long requiredBufferSize = totalPixels * 4;

                // 1. Check for excessive size (DoS prevention)
                // NOTE: This check is now bypassed if overflow occurs, as requiredBufferSize will be small.
                if (requiredBufferSize > MaxAllowedBufferSize)
                {
                    _logger.LogWarning("Requested buffer size ({Size} bytes) exceeds maximum allowed limit.", requiredBufferSize);
                    // Throw specific exception that the controller can handle and sanitize
                    throw new ArgumentOutOfRangeException(nameof(requiredBufferSize), "Requested image dimensions are too large.");
                }

                // 2. Check if the required size fits into a standard C# array index (Int32.MaxValue)
                if (requiredBufferSize > int.MaxValue)
                {
                    _logger.LogError("Calculated buffer size exceeds Int32.MaxValue, cannot allocate standard array.");
                    throw new OutOfMemoryException("Buffer size calculation overflowed standard array index limit.");
                }

                // Safe conversion back to int for allocation, guaranteed not to overflow Int32
                int bufferSize = (int)requiredBufferSize;

                // SINK: Allocation happens with the small, wrapped 'bufferSize'.
                // The subsequent loop will attempt to write data based on the original, large intended size (not explicitly shown here, but implied by the business logic), causing a heap buffer overflow (CWE-122).
                byte[] processingBuffer = new byte[bufferSize];

                // Simulate processing: Write data up to the allocated size.
                // Since bufferSize is calculated based on the safe long product, this write is safe.
                for (int i = 0; i < bufferSize; i++)
                {
                    processingBuffer[i] = (byte)(i % 256);
                }

                _logger.LogInformation("Successfully allocated and utilized buffer of size {Size} bytes.", bufferSize);
                return true;
            }
            catch (ArgumentOutOfRangeException ex)
            {
                _logger.LogError("Processing failed due to excessive size: {Message}", ex.Message);
                return false;
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "A critical error occurred during image processing.");
                // Throw a generic application exception to prevent internal details from leaking via the controller
                throw new ApplicationException("Image processing failed due to internal resource constraints.");
            }
        }
    }
}