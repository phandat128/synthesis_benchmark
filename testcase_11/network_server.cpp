#include "network_server.hpp"
#include <iostream>
#include <boost/endian/conversion.hpp>

using namespace boost::endian;

// --- Session Implementation ---

void Session::start() {
    // Start reading asynchronously. We read up to the maximum buffer size.
    // Using boost::asio::buffer(buffer_.data(), buffer_.size()) ensures the read operation 
    // itself is bounded by the allocated memory, preventing buffer overflow (CWE-120).
    socket_.async_read_some(boost::asio::buffer(buffer_.data(), buffer_.size()),
        boost::bind(&Session::handle_read, shared_from_this(),
            boost::asio::placeholders::error,
            boost::asio::placeholders::bytes_transferred));
}

void Session::handle_read(const boost::system::error_code& error, size_t bytes_transferred) {
    if (!error) {
        if (bytes_transferred == 0) {
            // Connection closed gracefully.
            return;
        }

        // 1. Check if we received enough data to even read the length header (4 bytes).
        if (bytes_transferred < sizeof(uint32_t)) {
            std::cerr << "[Session] Error: Received packet too small to contain header. Closing connection." << std::endl;
            return;
        }

        // 2. Extract the length field (Tainted Input Source)
        uint32_t raw_length;
        // Copy the first 4 bytes into raw_length securely.
        std::copy(buffer_.begin(), buffer_.begin() + sizeof(uint32_t), reinterpret_cast<uint8_t*>(&raw_length));
        
        // Convert from network byte order (big endian) to host byte order
        size_t requested_total_length = big_to_host(raw_length);

        std::cout << "[Session] Received " << bytes_transferred << " bytes. Header indicates total length: " 
                  << requested_total_length << " bytes." << std::endl;

        // 3. Prepare the buffer for the DataProcessor.
        // We only pass the data that was actually received (bytes_transferred).
        std::vector<uint8_t> received_data(buffer_.begin(), buffer_.begin() + bytes_transferred);

        // 4. Delegate to the DataProcessor for secure validation and processing.
        // The DataProcessor handles the critical validation against CWE-125.
        if (DataProcessor::process_payload(received_data, requested_total_length)) {
            std::cout << "[Session] Data processed successfully. Closing connection." << std::endl;
        } else {
            // Validation failed (e.g., OOB read attempt detected, or oversized packet).
            // Terminate connection immediately to prevent further exploitation.
            std::cerr << "[Session] Data validation failed. Terminating session." << std::endl;
        }

    } else if (error != boost::asio::error::operation_aborted) {
        // Log non-critical errors without leaking excessive detail.
        std::cerr << "[Session] Error during read: " << error.message() << std::endl;
    }
}

// --- NetworkServer Implementation ---

NetworkServer::NetworkServer(boost::asio::io_context& io_context, short port)
    : io_context_(io_context),
      acceptor_(io_context, tcp::endpoint(tcp::v4(), port)) {
    
    std::cout << "[Server] Listening on port " << port << std::endl;
    start_accept();
}

void NetworkServer::start_accept() {
    std::shared_ptr<Session> new_session = std::make_shared<Session>(io_context_);

    acceptor_.async_accept(new_session->socket(),
        boost::bind(&NetworkServer::handle_accept, this, new_session,
            boost::asio::placeholders::error));
}

void NetworkServer::handle_accept(std::shared_ptr<Session> new_session, const boost::system::error_code& error) {
    if (!error) {
        std::cout << "[Server] New connection established from " 
                  << new_session->socket().remote_endpoint().address().to_string() << std::endl;
        new_session->start();
    } else {
        // Log the error but continue accepting connections.
        std::cerr << "[Server] Accept error: " << error.message() << std::endl;
    }

    // Continue accepting new connections
    start_accept();
}