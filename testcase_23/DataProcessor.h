#ifndef DATAPROCESSOR_H
#define DATAPROCESSOR_H

#include <iostream>
#include <array>

// Define the fixed size of the internal buffer. This constant is used for bounds checking.
constexpr int BUFFER_SIZE = 100;

class DataProcessor {
private:
    // Fixed-size internal buffer.
    std::array<int, BUFFER_SIZE> internal_buffer = {};

public:
    DataProcessor();
    // Secure method to update a setting. Must perform bounds checking defensively.
    bool update_setting(int index, int value);
    void print_buffer() const;
};

#endif // DATAPROCESSOR_H