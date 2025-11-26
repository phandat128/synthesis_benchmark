#ifndef DATA_PROCESSOR_H
#define DATA_PROCESSOR_H

#include <stddef.h>
#include "constants.h"

// Structure to hold the parsed data fields
typedef struct {
    char id[FIELD_ID_SIZE + 1];     // +1 for null terminator
    char name[FIELD_NAME_SIZE + 1];
    char value[FIELD_VALUE_SIZE + 1];
} ParsedRecord;

/**
 * @brief Processes a raw data buffer, validating its size and extracting fields.
 *
 * This function is designed to be robust against size manipulation attacks (CWE-125).
 *
 * @param buffer The raw data buffer containing the record.
 * @param read_count The size of the data read from the external source (untrusted).
 * @param output_record Pointer to the structure where parsed data will be stored.
 * @return ProcessingResult Status code indicating success or type of error.
 */
ProcessingResult process_record(const char *buffer, size_t read_count, ParsedRecord *output_record);

#endif // DATA_PROCESSOR_H