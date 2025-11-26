#ifndef ITASK_H
#define ITASK_H

#include <iostream>

/**
 * @brief Abstract base class defining the interface for all executable tasks.
 * 
 * Secure Coding Practice: Includes a virtual destructor to ensure proper
 * cleanup of derived classes when deleting via a base class pointer (polymorphic destruction).
 */
class ITask {
public:
    // Pure virtual destructor is essential for polymorphic base classes.
    virtual ~ITask() = default;

    /**
     * @brief Executes the specific logic of the task.
     */
    virtual void execute() = 0;
};

#endif // ITASK_H