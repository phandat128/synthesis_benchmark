#ifndef DATABUFFER_H
#define DATABUFFER_H

#include <array>
#include <cstddef>

// Define the fixed size of the internal buffer.
// This constant is crucial for enforcing boundary checks.
constexpr size_t BUFFER_SIZE = 16;

/**
 * @brief DataBuffer holds the fixed-size array for configuration storage.
 * 
 * Using std::array provides compile-time size checking and clear semantics.
 */
class DataBuffer {
public:
    // The internal storage array.
    std::array<int, BUFFER_SIZE> internal_array;

    DataBuffer() {
        // Initialize buffer to zero upon creation.
        internal_array.fill(0);
    }
};

#endif // DATABUFFER_H