#ifndef CONCRETE_TASK_H
#define CONCRETE_TASK_H

#include "ITask.h"
#include <iostream>
#include <string>

/**
 * @brief Concrete implementation of a system task.
 */
class ConcreteTask : public ITask {
private:
    std::string task_name;

public:
    /**
     * @brief Constructor.
     */
    ConcreteTask(const std::string& name = "System Maintenance Task") 
        : task_name(name) {
        std::cout << "[ConcreteTask] Task '" << task_name << "' created." << std::endl;
    }

    /**
     * @brief Destructor.
     */
    ~ConcreteTask() override {
        std::cout << "[ConcreteTask] Task '" << task_name << "' destroyed." << std::endl;
    }

    /**
     * @brief Executes the task logic.
     * 
     * In a real-world scenario, this might involve complex system calls,
     * but here it simulates execution.
     */
    void execute() override {
        std::cout << "[ConcreteTask] Executing task: " << task_name << "..." << std::endl;
        // Simulate secure operation, e.g., validated file access or database interaction.
        std::cout << "[ConcreteTask] Operation successful. Status: OK." << std::endl;
    }
};

#endif // CONCRETE_TASK_H