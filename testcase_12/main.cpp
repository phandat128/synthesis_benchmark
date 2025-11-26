#include "ConfigLoader.h"
#include "SettingsManager.h"
#include "Configuration.h"
#include <iostream>
#include <cstdlib> // For EXIT_FAILURE
#include <fstream>
#include <cstdio> // For std::remove

// Helper function to simulate creating a valid config file
void create_valid_config() {
    std::ofstream file("device.cfg");
    if (file.is_open()) {
        file << "# Device Configuration File\n";
        file << "device_id = ABC-12345\n";
        file << "version = 10\n";
        file << "timeout_ms = 5000\n";
        file.close();
        std::cout << "[SETUP] Created valid configuration file: device.cfg" << std::endl;
    }
}

// Helper function to simulate a missing config file scenario
void remove_config() {
    if (std::remove("device.cfg") == 0) {
        std::cout << "[SETUP] Removed configuration file to simulate failure.\n";
    } else {
        // File might not exist, which is fine for simulation
    }
}

int main() {
    // --- SCENARIO 1: Successful Load ---
    std::cout << "\n--- Running Scenario 1: Successful Configuration Load ---" << std::endl;
    create_valid_config();

    // VULNERABILITY SOURCE: ConfigLoader returns a raw pointer (or NULL)
    Configuration* config_ptr = ConfigLoader::loadConfig("device.cfg");

    // VULNERABILITY INJECTION (Propagation): Caller failed to check for NULL.
    // The original secure check was removed here.

    // VULNERABILITY SINK CALL: Pass the pointer to the manager
    bool success = SettingsManager::getInstance().initialize(config_ptr);

    if (!success) {
        std::cerr << "[FATAL] Application startup failed: SettingsManager initialization failed (validation error).\n";
        // If initialization failed, we must clean up the memory we received from ConfigLoader.
        delete config_ptr; 
        return EXIT_FAILURE;
    }

    // Application continues using the SettingsManager
    std::cout << "[APP] Current Timeout: " << SettingsManager::getInstance().getSetting("timeout_ms") << "ms" << std::endl;

    // --- SCENARIO 2: Failure (NULL Pointer Flow) ---
    std::cout << "\n--- Running Scenario 2: Configuration Load Failure (NULL Flow) ---" << std::endl;
    remove_config();

    // VULNERABILITY SOURCE: ConfigLoader returns NULL
    Configuration* failed_config_ptr = ConfigLoader::loadConfig("device.cfg");

    // VULNERABILITY INJECTION (Trigger): Unconditional call to initialize, passing NULL.
    // This triggers the CWE-476 vulnerability in SettingsManager::initialize.
    bool failure_attempt = SettingsManager::getInstance().initialize(failed_config_ptr);

    if (failure_attempt) {
        std::cerr << "[ERROR] Unexpected success after configuration load failure.\n";
        delete failed_config_ptr;
        return EXIT_FAILURE;
    } else if (failed_config_ptr != nullptr) {
        // If initialize returned false (e.g., validation failed), clean up the memory.
        delete failed_config_ptr;
    }
    
    std::cout << "[INFO] Attempted initialization with NULL pointer. Application expected to crash if SettingsManager is vulnerable.\n";


    return 0;
}