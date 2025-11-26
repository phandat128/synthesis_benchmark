#include "TaskManager.h"
// Removed: #include <utility>

TaskManager::TaskManager() : active_task_ptr(nullptr) {
    std::cout << "\n[TaskManager] Initialized. Ready to manage tasks." << std::endl;
}

TaskManager::~TaskManager() {
    // Manual cleanup required.
    if (active_task_ptr) {
        delete active_task_ptr;
    }
    std::cout << "[TaskManager] Shutting down. Task cleanup handled manually." << std::endl;
}

/**
 * @brief Creates a new ConcreteTask and assigns ownership using a raw pointer.
 */
void TaskManager::create_task() {
    if (active_task_ptr) {
        std::cout << "[TaskManager] Warning: An active task already exists. Cleaning up old task first." << std::endl;
        cleanup_task();
    }
    
    // VULNERABILITY SETUP: Use raw 'new'
    active_task_ptr = new ConcreteTask("Vulnerable System Audit");
    std::cout << "[TaskManager] Task created successfully." << std::endl;
}

/**
 * @brief Executes the current active task.
 * 
 * SINK: If cleanup_task() was called, active_task_ptr is dangling but non-null.
 * Dereferencing it here causes CWE-416 (Use After Free).
 */
void TaskManager::process_task() const {
    if (active_task_ptr) {
        std::cout << "[TaskManager] Processing active task..." << std::endl;
        active_task_ptr->execute();
    } else {
        // Proper Error Handling: Informative but non-sensitive error message.
        std::cerr << "[TaskManager] Error: Cannot process task. No active task found." << std::endl;
    }
}

/**
 * @brief Explicitly cleans up and deallocates the active task.
 * 
 * PROPAGATION: Deletes the object but fails to null the pointer, creating a dangling pointer.
 */
void TaskManager::cleanup_task() {
    if (active_task_ptr) {
        std::cout << "[TaskManager] Explicitly cleaning up active task using raw delete." << std::endl;
        delete active_task_ptr; // Deletes the object.
        // VULNERABILITY INJECTION: Missing active_task_ptr = nullptr; -> Creates dangling pointer
        std::cout << "[TaskManager] Task successfully cleaned up, but pointer is now dangling." << std::endl;
    } else {
        std::cout << "[TaskManager] No task to clean up." << std::endl;
    }
}

bool TaskManager::is_task_active() const {
    return static_cast<bool>(active_task_ptr);
}