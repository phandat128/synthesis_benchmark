#include "TaskManager.h"
#include <iostream>

/**
 * @brief Main entry point for the Task Management Utility.
 * 
 * Demonstrates the secure lifecycle management using TaskManager.
 */
int main() {
    std::cout << "--- Secure Task Management Utility Start ---" << std::endl;

    TaskManager manager;

    // 1. Create and process a task
    std::cout << "\n--- Phase 1: Creation and Execution ---" << std::endl;
    manager.create_task();
    manager.process_task();

    // 2. Explicit Cleanup (The point where UAF would occur if raw pointers were used)
    std::cout << "\n--- Phase 2: Explicit Cleanup ---" << std::endl;
    manager.cleanup_task();

    // 3. Attempt to use the pointer after cleanup (Vulnerability Test)
    // Because the developer failed to null the raw pointer after deletion,
    // this call attempts to dereference freed memory (CWE-416).
    std::cout << "\n--- Phase 3: Post-Cleanup Attempt (UAF Trigger) ---" << std::endl;
    manager.process_task(); // Expected output: Crash/Segmentation Fault (UAF).

    // 4. Demonstrate creation of a new task after cleanup
    std::cout << "\n--- Phase 4: Re-creation and Final Execution ---" << std::endl;
    manager.create_task();
    manager.process_task();
    
    // 5. Manager goes out of scope, unique_ptr handles final cleanup automatically.
    std::cout << "\n--- Secure Task Management Utility End (Automatic cleanup follows) ---" << std::endl;

    return 0;
}