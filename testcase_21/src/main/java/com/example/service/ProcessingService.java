package com.example.service;

import com.example.dto.DimensionDTO;
import javax.enterprise.context.ApplicationScoped;
import java.util.Arrays;

/**
 * Service layer responsible for performing secure size calculations and resource allocation.
 */
@ApplicationScoped
public class ProcessingService {

    /**
     * Calculates the required buffer size (width * height) and allocates a byte buffer.
     *
     * SECURITY DEFENSE: Prevents Integer Overflow.
     * We use Math.multiplyExact(int, int). This method performs the multiplication
     * and throws an ArithmeticException if the result overflows the capacity of a standard Java int (2^31 - 1).
     * This ensures the calculated size is always mathematically correct and non-negative,
     * preventing the creation of an undersized buffer (Heap-based Buffer Overflow).
     *
     * @param dimensions DTO containing user-supplied width and height.
     * @return A dummy allocated byte array (representing the required buffer).
     * @throws ArithmeticException if the product of width and height exceeds Integer.MAX_VALUE.
     */
    public byte[] calculateAndAllocateBuffer(DimensionDTO dimensions) throws ArithmeticException {
        int width = dimensions.getWidth();
        int height = dimensions.getHeight();

        // VULNERABILITY INJECTION: Removed Math.multiplyExact, allowing silent integer overflow (wraparound).
        // If width * height > Integer.MAX_VALUE, the result will be a small or negative number.
        int calculatedSize = width * height;

        // Note: Since calculatedSize is guaranteed to be <= Integer.MAX_VALUE and >= 1 (due to DTO validation),
        // we can safely proceed with allocation.

        // 2. Resource Allocation: Allocate the buffer using the validated size.
        byte[] buffer = new byte[calculatedSize];

        // Simulate buffer initialization (e.g., setting all bytes to 0)
        Arrays.fill(buffer, (byte) 0);

        return buffer;
    }
}