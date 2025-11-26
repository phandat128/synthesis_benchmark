#ifndef DATA_PROCESSOR_HPP
#define DATA_PROCESSOR_HPP

#include <cstdint>
#include <vector>
#include <iostream>
#include <algorithm>

/**
 * @brief Utility class for securely parsing and validating network payloads.
 */
class DataProcessor {
public:
    // Define a safe maximum size for any payload we accept. 
    // This acts as a hard application limit.
    static constexpr size_t MAX_PAYLOAD_SIZE = 1024; 

    /**
     * @brief Processes the raw network payload buffer.
     * 
     * This function is the security boundary against the Out-of-bounds Read vulnerability (CWE-125).
     * It strictly validates the requested_length (tainted input) against the actual buffer size.
     * 
     * @param buffer The raw byte buffer received from the network.
     * @param requested_length The total length indicated by the packet header (tainted input).
     * @return true if processing was successful and secure, false otherwise.
     */
    static bool process_payload(const std::vector<uint8_t>& buffer, size_t requested_length);
};

#endif // DATA_PROCESSOR_HPP