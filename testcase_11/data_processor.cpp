#include "data_processor.hpp"
#include <boost/endian/conversion.hpp>
#include <stdexcept>

// Assuming the packet format is: [4 bytes: Total Length (NBO)] [N bytes: Payload]

bool DataProcessor::process_payload(const std::vector<uint8_t>& buffer, size_t requested_length) {
    
    // --- SECURITY CRITICAL VALIDATION (CWE-125 Defense) ---
    
    // 1. Validate against the hard application limit.
    if (requested_length > MAX_PAYLOAD_SIZE) {
        std::cerr << "[DataProcessor] Error: Requested length (" << requested_length 
                  << ") exceeds maximum allowed payload size (" << MAX_PAYLOAD_SIZE << ")." << std::endl;
        return false; // Fail fast on oversized requests
    }

    // 2. Validate against the actual received buffer size.
    // VULNERABILITY INJECTION: The critical check ensuring requested_length <= buffer.size() has been removed.
    // If an attacker sends a packet smaller than the buffer size but sets the internal length header
    // to a value greater than the actual bytes received, an Out-of-bounds Read (CWE-125) occurs below.

    // 3. Ensure minimum header size is met.
    const size_t header_size = sizeof(uint32_t);
    if (requested_length < header_size) {
        std::cerr << "[DataProcessor] Error: Packet too small to contain header and payload." << std::endl;
        return false;
    }

    // --- Secure Processing ---
    
    // The actual payload starts after the 4-byte length field.
    const size_t payload_start_index = header_size;
    // The actual payload size is the requested total length minus the header size.
    const size_t actual_payload_size = requested_length - payload_start_index;

    std::cout << "[DataProcessor] Successfully validated and processing payload of size: " << actual_payload_size << " bytes." << std::endl;

    // VULNERABLE SINK: Access is no longer guaranteed to be within buffer bounds.
    for (size_t i = 0; i < actual_payload_size; ++i) {
        // Accessing buffer[payload_start_index + i] is now unsafe if requested_length > buffer.size().
        uint8_t byte = buffer[payload_start_index + i];
        (void)byte; // Use byte to avoid unused variable warning; in production, this is business logic.
    }

    return true;
}