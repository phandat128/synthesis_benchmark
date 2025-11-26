#ifndef CLIENT_SESSION_HPP
#define CLIENT_SESSION_HPP

#include <boost/asio.hpp>
#include <boost/uuid/uuid.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <boost/uuid/uuid_generators.hpp>
#include <memory>
#include <string>
#include <iostream>
#include <array>
#include "AsyncTimer.hpp"

namespace asio = boost::asio;
using tcp = asio::ip::tcp;

// Forward declaration
class SessionManager;

/**
 * @brief Abstract base class for a secure, asynchronous client session.
 *
 * Inherits from enable_shared_from_this to ensure object lifetime is managed
 * by shared pointers, preventing Use-After-Free vulnerabilities (CWE-416)
 * when dealing with asynchronous handlers.
 */
class ClientSession : public std::enable_shared_from_this<ClientSession>
{
public:
    ClientSession(asio::io_context& io_context, SessionManager& manager)
        : socket_(io_context),
          manager_(manager),
          session_id_(boost::uuids::random_generator()()),
          timer_(std::make_shared<AsyncTimer>(io_context))
    {
        std::cout << "[Session] Created session: " << boost::uuids::to_string(session_id_) << std::endl;
    }

    virtual ~ClientSession() {
        std::cout << "[Session] Destroyed session: " << boost::uuids::to_string(session_id_) << std::endl;
    }

    // Prevent copy/move
    ClientSession(const ClientSession&) = delete;
    ClientSession& operator=(const ClientSession&) = delete;

    tcp::socket& socket() { return socket_; }
    boost::uuids::uuid get_id() const { return session_id_; }

    /**
     * @brief Starts the session: initiates the first read operation and heartbeat.
     */
    void start();

    /**
     * @brief Initiates the asynchronous read operation.
     */
    void start_read();

    /**
     * @brief Handles incoming data and routes commands.
     * @param data The raw data buffer.
     */
    void process_data(const std::string& data);

    /**
     * @brief Securely terminates the session, cancels async operations, and notifies the manager.
     */
    void terminate();

protected:
    /**
     * @brief Schedules the periodic heartbeat check using the AsyncTimer.
     */
    void start_heartbeat();

    /**
     * @brief The actual cleanup logic, called safely only when the session is valid.
     */
    virtual void process_cleanup() {
        std::cout << "[Session] Performing internal cleanup for " << boost::uuids::to_string(session_id_) << std::endl;
    }

    tcp::socket socket_;
    SessionManager& manager_;
    boost::uuids::uuid session_id_;
    std::array<char, 1024> read_buffer_; // Fixed size buffer for reads

    // Securely managed timer object
    std::shared_ptr<AsyncTimer> timer_;
};

#endif // CLIENT_SESSION_HPP