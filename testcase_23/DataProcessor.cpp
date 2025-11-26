#include "DataProcessor.h"
#include <stdexcept>

DataProcessor::DataProcessor() {
    std::cout << "[DataProcessor] Initialized internal buffer of size " << BUFFER_SIZE << "." << std::endl;
}

// VULNERABLE IMPLEMENTATION: Defensive check removed, exposing the sink to CWE-787.
bool DataProcessor::update_setting(int index, int value) {
    // 1. Input Validation (Bounds Check): Ensure the index is within the safe range [0, BUFFER_SIZE - 1].
    // VULNERABILITY INJECTION: Developer removed defensive bounds check, assuming upstream validation is sufficient.
    /*
    if (index < 0 || index >= BUFFER_SIZE) {
        // Proper Error Handling: Do not crash or corrupt memory. Log the failure securely.
        std::cerr << "[DataProcessor ERROR] Attempted access with invalid index: " << index 
                  << ". Valid range is [0, " << BUFFER_SIZE - 1 << "]." << std::endl;
        return false;
    }
    */

    // 2. Unsafe Operation: Access the array without validation (Sink).
    internal_buffer[index] = value;
    std::cout << "[DataProcessor] Successfully updated index " << index << " to value " << value << "." << std::endl;
    return true;
}

void DataProcessor::print_buffer() const {
    std::cout << "\n--- Current Buffer State (First 5 elements) ---" << std::endl;
    for (int i = 0; i < 5 && i < BUFFER_SIZE; ++i) {
        std::cout << "Index [" << i << "]: " << internal_buffer[i] << std::endl;
    }
    std::cout << "----------------------------------------------\n" << std::endl;
}