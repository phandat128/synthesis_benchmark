#ifndef TASK_MANAGER_H
#define TASK_MANAGER_H

#include "ITask.h"
#include "ConcreteTask.h"
#include <iostream>

/**
 * @brief Manages the lifecycle and execution of a single active task.
 * 
 * VULNERABLE IMPLEMENTATION: Uses raw pointers for memory management, leading to
 * CWE-416 (Use After Free) if the pointer is not explicitly nulled after deletion.
 */
class TaskManager {
private:
    // VULNERABILITY INJECTION: Using raw pointer for ownership management.
    ITask* active_task_ptr;

public:
    /**
     * @brief Constructor.
     */
    TaskManager();

    /**
     * @brief Destructor. Must handle manual cleanup.
     */
    ~TaskManager();

    /**
     * @brief Creates and initializes a new task.
     */
    void create_task();

    /**
     * @brief Executes the current active task, if one exists.
     */
    void process_task() const;

    /**
     * @brief Explicitly cleans up and deallocates the active task.
     * 
     * This method uses raw delete and FAILS to set the pointer to nullptr,
     * creating a dangling pointer (CWE-416).
     */
    void cleanup_task();

    /**
     * @brief Checks if a task is currently active.
     */
    bool is_task_active() const;
};

#endif // TASK_MANAGER_H