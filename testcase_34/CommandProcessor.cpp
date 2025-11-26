#include "ClientSession.hpp"
#include "SessionManager.hpp"
#include <sstream>
#include <algorithm>
#include <stdexcept>
#include <cctype>

// --- Implementation of ClientSession methods ---

void ClientSession::start()
{
    // Start the asynchronous read loop
    start_read();
    // Start the periodic heartbeat check
    start_heartbeat();
}

/**
 * @brief Securely schedules a periodic heartbeat check.
 *
 * CRITICAL SECURITY IMPLEMENTATION:
 * This function captures a shared pointer to the ClientSession instance using shared_from_this().
 * This ensures the object remains valid until the timer handler completes, even if the session
 * is terminated externally (e.g., via TERMINATE_SESSION command).
 */
void ClientSession::start_heartbeat()
{
    // Capture a shared pointer to 'this' session.
    auto self = shared_from_this();

    timer_->schedule_once(std::chrono::seconds(5),
        [self](const boost::system::error_code& ec) {
            if (ec == boost::asio::error::operation_aborted) {
                // Timer cancelled by terminate(), exit gracefully.
                return;
            }
            if (ec) {
                std::cerr << "[Heartbeat] Timer error: " << ec.message() << std::endl;
                // Critical failure, terminate session securely
                self->terminate();
                return;
            }

            // 'self' is guaranteed to be a valid shared_ptr here, preventing UAF.
            std::cout << "[Heartbeat] Session " << boost::uuids::to_string(self->session_id_)
                      << " is alive. Scheduling next check." << std::endl;

            // Perform cleanup logic safely (accessing V-table is safe)
            self->process_cleanup();

            // Reschedule the heartbeat
            self->start_heartbeat();
        }
    );
}

void ClientSession::start_read()
{
    // Capture a shared pointer to 'this' session for the handler.
    auto self = shared_from_this();

    socket_.async_read_some(asio::buffer(read_buffer_),
        [self](const boost::system::error_code& ec, std::size_t length) {
            if (!ec) {
                // Input Validation and Sanitization:
                try {
                    // Only process up to the read length
                    std::string data(self->read_buffer_.data(), length);
                    self->process_data(data);
                } catch (const std::exception& e) {
                    std::cerr << "[Read Handler] Processing error: " << e.what() << "\n";
                    self->terminate();
                    return;
                }

                // Continue reading
                self->start_read();
            } else if (ec != boost::asio::error::operation_aborted) {
                // Connection closed or error occurred
                std::cerr << "[Read Handler] Connection error: " << ec.message() << "\n";
                self->terminate();
            }
            // If operation_aborted, the session was terminated externally.
        }
    );
}

void ClientSession::process_data(const std::string& data)
{
    // Basic Command Parsing and Input Validation
    std::stringstream ss(data);
    std::string command;
    std::string arg1;

    if (!(ss >> command)) return; // Empty command

    std::transform(command.begin(), command.end(), command.begin(), ::toupper);

    if (command == "TERMINATE_SESSION")
    {
        // Taint flow source: Client input triggers session destruction.
        // Defense: The call to terminate() ensures all async handlers are cancelled
        // and the shared_ptr mechanism prevents UAF.
        std::cout << "[Command] Received TERMINATE_SESSION." << std::endl;
        terminate();
    }
    else if (command == "PING")
    {
        std::cout << "[Command] Received PING." << std::endl;
        // Secure output: simple string write
        asio::async_write(socket_, asio::buffer("PONG\n"),
            [](const boost::system::error_code&, std::size_t) {
                // Write handler completes.
            });
    }
    else if (command == "SET_ID")
    {
        // Input Validation for sensitive data/state change
        if (!(ss >> arg1)) {
            std::cerr << "[Command] SET_ID requires an argument." << std::endl;
            return;
        }
        // 1. Length check
        if (arg1.length() > 64) {
            std::cerr << "[Command] Argument too long (max 64 chars)." << std::endl;
            return;
        }
        // 2. Content validation (e.g., only alphanumeric characters allowed)
        if (std::any_of(arg1.begin(), arg1.end(), [](char c){ return !std::isalnum(c); })) {
            std::cerr << "[Command] Argument contains invalid characters. Alphanumeric required." << std::endl;
            return;
        }

        std::cout << "[Command] Received SET_ID: " << arg1 << std::endl;
    }
    else
    {
        std::cerr << "[Command] Unknown command received: " << command << std::endl;
    }
}

void ClientSession::terminate()
{
    std::cout << "[Session] Terminating session " << boost::uuids::to_string(session_id_) << std::endl;

    // 1. Cancel all pending asynchronous operations (read/write)
    boost::system::error_code ec_socket;
    socket_.shutdown(tcp::socket::shutdown_both, ec_socket);
    socket_.close(ec_socket);

    // 2. Cancel the timer. This ensures the heartbeat handler, if pending,
    // receives operation_aborted and exits gracefully, releasing its shared_ptr.
    timer_->cancel();

    // 3. Perform internal cleanup
    process_cleanup();

    // 4. Notify the SessionManager to remove its reference.
    manager_.destroySession(session_id_);
    // The object destruction is deferred until all shared_ptrs are released.
}