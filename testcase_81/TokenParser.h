#ifndef TOKENPARSER_H
#define TOKENPARSER_H

#include "DataTypes.h"
#include <stdexcept>
#include <string>
#include <vector>

// Custom exception for parsing errors
class TokenParsingError : public std::runtime_error {
public:
    TokenParsingError(const std::string& msg) : std::runtime_error(msg) {}
};

class TokenParser {
public:
    // Parses the raw token buffer and returns a validated configuration structure.
    // Throws TokenParsingError on failure.
    ParsedConfig parseToken(const std::vector<uint8_t>& raw_buffer);

private:
    // Internal state management during parsing
    const uint8_t* current_pos_ = nullptr;
    const uint8_t* buffer_end_ = nullptr;

    /**
     * @brief Secure helper function to read a field safely from the buffer.
     * @param requested_size The size requested by the untrusted header.
     * @param output_buffer The destination buffer.
     * @return true on success, false on failure (bounds violation).
     */
    bool readField(size_t requested_size, uint8_t* output_buffer);

    // Helper function to simulate cryptographic validation
    bool validateChecksum(const TokenHeader& header, const uint8_t* payload, size_t payload_size);
};

#endif // TOKENPARSER_H