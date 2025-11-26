#ifndef DATATYPES_H
#define DATATYPES_H

#include <cstdint>
#include <array>
#include <cstddef>

// Define constants for the token structure
constexpr size_t TOKEN_BUFFER_SIZE = 1024; // Maximum size of the raw token buffer
constexpr size_t TOKEN_HEADER_SIZE = 16;   // Size of the fixed header (4x 4-byte fields)

// Structure representing the configuration token header
// Using fixed-width integers for portability and security.
struct TokenHeader {
    uint32_t magic_number; // Expected value for validation
    uint32_t version;      // Token version
    uint32_t data_length;  // CRITICAL: Untrusted input defining the payload size (CWE-125 source)
    uint32_t checksum;     // Simple checksum/hash of the payload
};

// Structure to hold the securely parsed configuration
struct ParsedConfig {
    uint32_t version;
    // Use a fixed-size array for the payload, sized to the maximum possible payload.
    std::array<uint8_t, TOKEN_BUFFER_SIZE - TOKEN_HEADER_SIZE> payload;
    size_t actual_payload_size = 0; // Tracks the actual valid data size within the payload array
};

#endif // DATATYPES_H