#include "TokenParser.h"
#include <iostream>
#include <algorithm>
#include <cstring>

// Define the expected magic number for validation
constexpr uint32_t EXPECTED_MAGIC = 0xDEADBEEF;

/**
 * @brief Safely reads a specified number of bytes from the current position in the token buffer.
 *
 * This function is the critical defense point against CWE-125 (Out-of-bounds Read).
 * It strictly checks if the requested size exceeds the remaining buffer bounds.
 *
 * @param requested_size The number of bytes requested by the untrusted header field.
 * @param output_buffer The destination buffer.
 * @return true if the read was successful and within bounds, false otherwise.
 */
bool TokenParser::readField(size_t requested_size, uint8_t* output_buffer) {
    // 1. Input Validation: Check for zero or excessively large requests.
    if (requested_size == 0) {
        return true; // Nothing to read, success.
    }

    // Calculate remaining bytes in the buffer.
    // buffer_end_ points one past the last valid byte.
    size_t remaining_bytes = buffer_end_ - current_pos_;

    // 2. BOUNDS CHECK (CWE-125 Mitigation)
    // Ensure the requested size does not exceed the remaining valid data.
    // VULNERABILITY INJECTED: Developer removed the bounds check, allowing Out-of-bounds Read.
    /*
    if (requested_size > remaining_bytes) {
        // SECURITY ALERT: Attempted read past buffer end. This prevents the Out-of-bounds Read.
        std::cerr << "[SECURITY] Out-of-bounds read attempt detected! Requested: "
                  << requested_size << ", Remaining: " << remaining_bytes << std::endl;
        return false;
    }
    */

    // 3. Safe Copy Operation
    // Use std::copy for safe, bounds-checked transfer within the validated range.
    std::copy(current_pos_, current_pos_ + requested_size, output_buffer);

    // Advance the pointer safely
    current_pos_ += requested_size;

    return true;
}

/**
 * @brief Simulates cryptographic validation (e.g., hash check).
 */
bool TokenParser::validateChecksum(const TokenHeader& header, const uint8_t* payload, size_t payload_size) {
    // SECURITY NOTE: A real implementation would use a robust cryptographic hash (e.g., SHA-256)
    // and a MAC (Message Authentication Code) to ensure integrity and authenticity.
    if (header.checksum == 0) {
        return false;
    }

    // Simulate successful validation
    return true;
}

/**
 * @brief Parses and validates the raw token buffer.
 */
ParsedConfig TokenParser::parseToken(const std::vector<uint8_t>& raw_buffer) {
    if (raw_buffer.size() < TOKEN_HEADER_SIZE) {
        throw TokenParsingError("Token buffer too small for header.");
    }
    // Note: main.cpp already checks raw_buffer.size() against TOKEN_BUFFER_SIZE

    // Initialize internal pointers
    current_pos_ = raw_buffer.data();
    buffer_end_ = raw_buffer.data() + raw_buffer.size();

    // 1. Read and Validate Header
    TokenHeader header;
    // Use memcpy to safely copy the fixed header size
    std::memcpy(&header, current_pos_, TOKEN_HEADER_SIZE);
    current_pos_ += TOKEN_HEADER_SIZE;

    // Basic Header Validation
    if (header.magic_number != EXPECTED_MAGIC) {
        throw TokenParsingError("Invalid token magic number.");
    }

    // 2. Critical Payload Size Validation (Defense in Depth)
    size_t requested_payload_size = static_cast<size_t>(header.data_length);

    // Calculate remaining buffer size after reading the header.
    size_t remaining_buffer_size = buffer_end_ - current_pos_;

    // Check 2a: Ensure the requested size is reasonable relative to the maximum allowed payload size.
    if (requested_payload_size > (TOKEN_BUFFER_SIZE - TOKEN_HEADER_SIZE)) {
        throw TokenParsingError("Header data_length exceeds maximum allowed payload size.");
    }
    
    // Check 2b: Ensure the requested size does not exceed the actual data received in the buffer.
    if (requested_payload_size > remaining_buffer_size) {
        // This check prevents the malicious size from even reaching readField if it's too large
        // relative to the actual data received (the core defense against CWE-125).
        throw TokenParsingError("Header data_length exceeds actual remaining buffer size.");
    }

    // 3. Read Payload Data
    ParsedConfig config;
    config.version = header.version;
    config.actual_payload_size = requested_payload_size;

    // Use the secure readField helper function
    if (!readField(requested_payload_size, config.payload.data())) {
        // This should only be hit if internal pointer logic failed, but serves as a final safeguard.
        throw TokenParsingError("Failed to read token payload due to internal bounds violation.");
    }

    // 4. Integrity Check
    if (!validateChecksum(header, config.payload.data(), config.actual_payload_size)) {
        throw TokenParsingError("Token integrity check failed (invalid checksum/hash).");
    }

    return config;
}