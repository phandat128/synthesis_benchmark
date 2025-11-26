#ifndef NETWORK_SERVER_HPP
#define NETWORK_SERVER_HPP

#include <boost/asio.hpp>
#include <boost/bind/bind.hpp>
#include <vector>
#include <memory>
#include "data_processor.hpp"

using boost::asio::ip::tcp;

// Define the maximum buffer size for a single asynchronous read operation.
// This size must accommodate the largest acceptable packet (MAX_PAYLOAD_SIZE + Header).
constexpr size_t MAX_PACKET_BUFFER_SIZE = DataProcessor::MAX_PAYLOAD_SIZE + sizeof(uint32_t);

/**
 * @brief Represents a single client connection session.
 * Handles asynchronous reading and writing for one client.
 */
class Session : public std::enable_shared_from_this<Session> {
public:
    Session(boost::asio::io_context& io_context)
        : socket_(io_context) {}

    tcp::socket& socket() {
        return socket_;
    }

    /**
     * @brief Starts the asynchronous read operation for the session.
     */
    void start();

private:
    /**
     * @brief Handler for asynchronous read completion.
     * Delegates the received raw buffer to the DataProcessor for validation.
     */
    void handle_read(const boost::system::error_code& error, size_t bytes_transferred);

    tcp::socket socket_;
    // Fixed-size buffer for receiving raw network data. Limited to MAX_PACKET_BUFFER_SIZE.
    std::vector<uint8_t> buffer_{MAX_PACKET_BUFFER_SIZE};
};

/**
 * @brief The main asynchronous TCP server class.
 */
class NetworkServer {
public:
    NetworkServer(boost::asio::io_context& io_context, short port);

private:
    /**
     * @brief Initiates the asynchronous accept operation.
     */
    void start_accept();

    /**
     * @brief Handler for asynchronous accept completion.
     */
    void handle_accept(std::shared_ptr<Session> new_session, const boost::system::error_code& error);

    boost::asio::io_context& io_context_;
    tcp::acceptor acceptor_;
};

#endif // NETWORK_SERVER_HPP