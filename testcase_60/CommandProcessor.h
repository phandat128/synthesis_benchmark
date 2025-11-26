#ifndef COMMANDPROCESSOR_H
#define COMMANDPROCESSOR_H

#include <string>
#include <vector>
#include "DataBuffer.h"

/**
 * @brief Handles parsing, validation, and execution of user commands.
 */
class CommandProcessor {
public:
    /**
     * @brief Constructor.
     * @param buffer Reference to the DataBuffer instance.
     */
    CommandProcessor(DataBuffer& buffer);

    /**
     * @brief Parses and executes a single command string.
     * @param input The raw command string from stdin.
     * @return true if the application should continue running, false otherwise.
     */
    bool process_command(const std::string& input);

private:
    DataBuffer& buffer_;

    /**
     * @brief Executes the SET command, securely writing data to the buffer.
     * @param value The integer value to write.
     * @param index The target index (must be validated).
     * @return true on success, false on failure (e.g., bounds violation).
     */
    bool execute_set(int value, int index);

    /**
     * @brief Executes the GET command, securely reading data from the buffer.
     * @param index The target index (must be validated).
     * @return true on success, false on failure.
     */
    bool execute_get(int index);

    /**
     * @brief Helper function to split the input string into tokens.
     */
    std::vector<std::string> tokenize(const std::string& input);
};

#endif // COMMANDPROCESSOR_H