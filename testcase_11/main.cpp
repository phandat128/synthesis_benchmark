#include "network_server.hpp"
#include <iostream>
#include <stdexcept>

// Define the port the service will listen on
constexpr short LISTEN_PORT = 8080;

/**
 * @brief Main entry point for the network service.
 * Initializes Boost.Asio and starts the asynchronous server loop.
 */
int main() {
    try {
        // 1. Setup Boost.Asio IO Context
        boost::asio::io_context io_context;

        // 2. Initialize the Network Server
        // The server handles the acceptor and delegates connections to sessions.
        NetworkServer server(io_context, LISTEN_PORT);

        // 3. Run the asynchronous event loop
        std::cout << "[Main] Starting IO context loop..." << std::endl;
        io_context.run();
        
    } catch (const std::exception& e) {
        // Robust error handling for initialization failures
        std::cerr << "[Main] Fatal Exception: " << e.what() << std::endl;
        return EXIT_FAILURE;
    } catch (...) {
        std::cerr << "[Main] Unknown Fatal Exception occurred." << std::endl;
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}