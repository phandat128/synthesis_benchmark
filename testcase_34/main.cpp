#include <boost/asio.hpp>
#include <boost/bind/bind.hpp>
#include <iostream>
#include <memory>
#include <vector>
#include <thread>
#include "ClientSession.hpp"
#include "SessionManager.hpp"

namespace asio = boost::asio;
using tcp = asio::ip::tcp;

/**
 * @brief The main server class responsible for accepting new connections.
 */
class Server
{
public:
    Server(asio::io_context& io_context, short port)
        : io_context_(io_context),
          acceptor_(io_context, tcp::endpoint(tcp::v4(), port))
    {
        std::cout << "Server listening on port " << port << std::endl;
        start_accept();
    }

private:
    /**
     * @brief Initiates the asynchronous acceptance process.
     * 
     * A new ClientSession is created using std::make_shared, ensuring its lifetime
     * is managed securely by shared pointers.
     */
    void start_accept()
    {
        // Create a new ClientSession object using std::make_shared.
        std::shared_ptr<ClientSession> new_session =
            std::make_shared<ClientSession>(io_context_, SessionManager::instance());

        // Asynchronously wait for a new connection
        acceptor_.async_accept(new_session->socket(),
            boost::bind(&Server::handle_accept, this, new_session, 
                        asio::placeholders::error));
    }

    /**
     * @brief Handler for completed accept operations.
     * @param new_session The shared pointer to the newly created session.
     * @param error The error code.
     */
    void handle_accept(std::shared_ptr<ClientSession> new_session,
                       const boost::system::error_code& error)
    {
        if (!error)
        {
            std::cout << "[Accept] New connection established." << std::endl;

            // 1. Add the session to the manager's map.
            SessionManager::instance().addSession(new_session);

            // 2. Start the session's read loop and timers.
            new_session->start();
        }
        else
        {
            std::cerr << "[Accept] Error: " << error.message() << std::endl;
            // If acceptance fails, the shared_ptr 'new_session' is released here.
        }

        // Continue accepting new connections
        start_accept();
    }

    asio::io_context& io_context_;
    tcp::acceptor acceptor_;
};

int main()
{
    try
    {
        // Determine optimal concurrency based on hardware
        const unsigned short concurrency_hint = std::thread::hardware_concurrency();
        asio::io_context io_context(concurrency_hint);

        Server s(io_context, 8080);

        std::vector<std::thread> threads;
        for (unsigned short i = 0; i < concurrency_hint; ++i)
        {
            threads.emplace_back([&io_context] {
                // Run the event loop. Use a try-catch block for robust error handling.
                try {
                    io_context.run();
                } catch (const std::exception& e) {
                    // Do not leak sensitive internal details in production logs
                    std::cerr << "[IO Thread] Unhandled exception occurred." << std::endl;
                }
            });
        }

        for (auto& t : threads)
        {
            t.join();
        }
    }
    catch (std::exception& e)
    {
        std::cerr << "Exception in main: " << e.what() << "\n";
    }

    return 0;
}