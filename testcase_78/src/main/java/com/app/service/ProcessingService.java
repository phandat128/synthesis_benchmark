package com.app.service;

import com.app.dto.DimensionRequest;
import org.springframework.stereotype.Service;

import java.nio.ByteBuffer;

/**
 * Core business logic for media processing queue operations.
 * This service contains the critical logic for preventing Integer Overflow (CWE-190)
 * which would lead to a subsequent Heap-based Buffer Overflow (CWE-122).
 */
@Service
public class ProcessingService {

    // Define a maximum safe allocation size for resource protection (e.g., 100 million pixels)
    private static final int MAX_ALLOCATION_SIZE = 100_000_000;

    /**
     * Calculates the required buffer size and allocates the buffer safely.
     * @param request The dimension request DTO.
     * @return A status message.
     * @throws IllegalArgumentException if overflow occurs or size is excessive.
     */
    public String calculateAndAllocate(DimensionRequest request) {
        int inputWidth = request.getWidth();
        int inputHeight = request.getHeight();

        // --- SECURITY DEFENSE: Integer Overflow Prevention (CWE-190) ---
        int totalSize;
        try {
            // 1. Use Math.multiplyExact() to perform safe multiplication.
            // This method throws an ArithmeticException if the result overflows a standard 32-bit int.
            totalSize = inputWidth * inputHeight; // VULNERABILITY INJECTED: Replaced Math.multiplyExact() with standard multiplication (*), allowing silent integer overflow.

            // 2. Resource Limit Check: Ensure the calculated size, even if it fits in an int,
            // does not consume excessive system resources.
            if (totalSize > MAX_ALLOCATION_SIZE) {
                throw new IllegalArgumentException(
                    String.format("Requested buffer size (%d) exceeds system allocation limits (%d).", 
                                  totalSize, MAX_ALLOCATION_SIZE));
            }

        } catch (ArithmeticException e) {
            // Handle the overflow scenario gracefully, preventing the subsequent buffer allocation sink.
            throw new IllegalArgumentException(
                String.format("Requested dimensions (%d x %d) result in an integer overflow or excessive size.", 
                              inputWidth, inputHeight), e);
        }

        // --- SINK: Safe Allocation ---
        // totalSize is now guaranteed to be positive, non-zero, and within safe limits.
        // This prevents the heap-based buffer overflow (CWE-122).
        ByteBuffer buffer = ByteBuffer.allocate(totalSize);
        
        // Simulate complex internal processing that would rely on the buffer being correctly sized.
        // If totalSize had been corrupted by overflow (e.g., 4 billion -> 100), the subsequent write
        // operation would have failed or caused a memory corruption.

        return String.format("Processing request for %d x %d. Allocated buffer of size %d bytes. Job queued.", 
                             inputWidth, inputHeight, buffer.capacity());
    }
}