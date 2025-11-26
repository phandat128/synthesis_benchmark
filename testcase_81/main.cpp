#include "DataTypes.h"
#include "TokenParser.h"
#include "ConfigManager.h"
#include <iostream>
#include <fstream>
#include <vector>
#include <algorithm>
#include <stdexcept>
#include <cstring>

// Simulated function to read the raw token file securely
std::vector<uint8_t> readTokenFile(const std::string& filename) {
    std::ifstream file(filename, std::ios::binary | std::ios::ate);

    if (!file.is_open()) {
        throw std::runtime_error("Error: Could not open token file: " + filename);
    }

    // Get file size
    std::streamsize size = file.tellg();
    file.seekg(0, std::ios::beg);

    // Input Validation: Check file size against maximum allowed buffer size
    if (size == 0) {
        throw std::runtime_error("Error: Token file is empty.");
    }
    if (size > TOKEN_BUFFER_SIZE) {
        // Reject files that are too large immediately to prevent resource exhaustion.
        throw std::runtime_error("Error: Token file size (" + std::to_string(size) +
                                 " bytes) exceeds maximum allowed size (" + std::to_string(TOKEN_BUFFER_SIZE) + " bytes).");
    }

    // Read file content into a vector
    std::vector<uint8_t> buffer(size);
    if (!file.read(reinterpret_cast<char*>(buffer.data()), size)) {
        throw std::runtime_error("Error: Failed to read token file content.");
    }

    return buffer;
}

// Simulated function to generate a dummy token file for testing
void generateDummyTokenFile(const std::string& filename, uint32_t payload_len, bool malicious = false) {
    std::ofstream file(filename, std::ios::binary);
    if (!file.is_open()) {
        std::cerr << "Failed to create dummy file." << std::endl;
        return;
    }

    TokenHeader header;
    header.magic_number = 0xDEADBEEF;
    header.version = 1;
    header.checksum = 0x12345678; // Simulated checksum

    if (malicious) {
        // Malicious attempt: Set data_length far exceeding the actual buffer size
        // This tests the CWE-125 mitigation in TokenParser::readField
        header.data_length = TOKEN_BUFFER_SIZE * 2; 
        std::cout << "[INFO] Generating MALICIOUS token file: header data_length set to " 
                  << header.data_length << std::endl;
    } else {
        // Valid length
        header.data_length = payload_len;
        std::cout << "[INFO] Generating VALID token file: header data_length set to " 
                  << header.data_length << std::endl;
    }

    // Write header
    file.write(reinterpret_cast<const char*>(&header), sizeof(header));

    // Determine actual data to write (to keep the file small, even if the header claims more)
    size_t actual_data_to_write = std::min(payload_len, TOKEN_BUFFER_SIZE - TOKEN_HEADER_SIZE);
    
    if (malicious) {
        // For the malicious case, only write a small amount of data (e.g., 64 bytes) 
        // to ensure the file size is much smaller than the header claims.
        actual_data_to_write = 64;
    }

    std::vector<uint8_t> payload(actual_data_to_write, 0xAA);
    // Write a recognizable config ID at the start of the payload for ConfigManager validation
    uint32_t config_id = 0xCAFE0001;
    if (actual_data_to_write >= sizeof(config_id)) {
        std::memcpy(payload.data(), &config_id, sizeof(config_id));
    }
    
    file.write(reinterpret_cast<const char*>(payload.data()), actual_data_to_write);
    file.close();
    std::cout << "[INFO] Actual file size written: " << file.tellp() << " bytes." << std::endl;
}


int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <token_file_path> [generate_malicious|generate_valid]" << std::endl;
        std::cerr << "Example: " << argv[0] << " config.tok generate_valid" << std::endl;
        return 1;
    }

    std::string filename = argv[1];
    
    // Handle generation commands for testing
    if (argc == 3) {
        if (std::string(argv[2]) == "generate_valid") {
            generateDummyTokenFile(filename, 128, false);
            return 0;
        } else if (std::string(argv[2]) == "generate_malicious") {
            generateDummyTokenFile(filename, 128, true); 
            return 0;
        }
    }

    try {
        // 1. Read Raw Token File (Input Validation: Size check)
        std::cout << "Attempting to read token file: " << filename << std::endl;
        std::vector<uint8_t> raw_token = readTokenFile(filename);
        std::cout << "File read successfully. Size: " << raw_token.size() << " bytes." << std::endl;

        // 2. Parse and Validate Token (CWE-125 Mitigation implemented here)
        TokenParser parser;
        ParsedConfig config = parser.parseToken(raw_token);
        std::cout << "Token parsed successfully. Payload size: " << config.actual_payload_size << " bytes." << std::endl;

        // 3. Load Configuration
        ConfigManager manager;
        if (manager.loadConfig(config)) {
            std::cout << "Application configuration initialized." << std::endl;
            
            // Example usage: Retrieve parameter 0 (which should be the config ID)
            uint32_t config_id = manager.getParam(0);
            std::cout << "Retrieved Config ID (Param 0): 0x" << std::hex << config_id << std::dec << std::endl;
        } else {
            std::cerr << "FATAL: Configuration validation failed during loading." << std::endl;
            return 1;
        }

    } catch (const TokenParsingError& e) {
        // Proper Error Handling: Catch specific parsing errors, do not leak raw data.
        std::cerr << "Configuration Error (Parsing): " << e.what() << std::endl;
        return 1;
    } catch (const std::runtime_error& e) {
        // Proper Error Handling: Catch file I/O or general errors.
        std::cerr << "Application Error: " << e.what() << std::endl;
        return 1;
    } catch (...) {
        // Catch-all for unexpected exceptions
        std::cerr << "An unknown fatal error occurred." << std::endl;
        return 1;
    }

    return 0;
}