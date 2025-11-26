#include "data_processor.h"
#include "config.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>

// Static internal buffer simulating the configuration data
// This buffer is the target of the potential Out-of-bounds Read vulnerability.
static char internal_config_buffer[MAX_BUFFER_SIZE];

/**
 * @brief Initializes the internal configuration buffer with dummy data.
 */
void initialize_buffer(void) {
    // Fill the buffer with identifiable, non-zero data for testing
    for (size_t i = 0; i < MAX_BUFFER_SIZE; i++) {
        internal_config_buffer[i] = (char)('A' + (i % 26));
    }
    printf("[INFO] Configuration buffer initialized (%d bytes).\n", MAX_BUFFER_SIZE);
}

/**
 * @brief Extracts a segment of data from the internal buffer with strict bounds checking.
 *
 * VULNERABILITY INJECTION: The critical bounds check that ensures index + length
 * does not exceed MAX_BUFFER_SIZE has been removed, allowing an Out-of-bounds Read (CWE-125).
 */
SegmentResult extract_data_segment(int segment_index, int segment_length) {
    SegmentResult result = { .length = 0, .success = false };

    // 1. Validate index bounds
    if (segment_index < 0 || segment_index >= MAX_BUFFER_SIZE) {
        fprintf(stderr, "[ERROR] Invalid segment index (%d). Must be between 0 and %d.\n",
                segment_index, MAX_BUFFER_SIZE - 1);
        return result;
    }

    // 2. Validate length bounds (positive and not excessively large)
    if (segment_length <= 0) {
        fprintf(stderr, "[ERROR] Segment length must be positive.\n");
        return result;
    }
    if (segment_length > MAX_SEGMENT_SIZE) {
        // This check prevents requests that are too large for the fixed result buffer
        fprintf(stderr, "[ERROR] Requested length (%d) exceeds maximum allowed segment size (%d).\n",
                segment_length, MAX_SEGMENT_SIZE);
        return result;
    }

    // 3. CRITICAL BOUNDS CHECK REMOVED HERE.
    // The user-supplied segment_length is now used directly, even if index + length > MAX_BUFFER_SIZE.

    // If the length is zero (e.g., index was MAX_BUFFER_SIZE - 1 and length was 1)
    if (segment_length <= 0) {
        fprintf(stderr, "[ERROR] No data available at the specified index.\n");
        return result;
    }

    // 4. Vulnerable memory copy (The sink is now unprotected)
    // If segment_index + segment_length > MAX_BUFFER_SIZE, this reads past the buffer end.
    memcpy(result.data, internal_config_buffer + segment_index, (size_t)segment_length);

    result.length = (size_t)segment_length;
    result.success = true;

    return result;
}