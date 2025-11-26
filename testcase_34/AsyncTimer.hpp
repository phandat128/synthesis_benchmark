#ifndef ASYNC_TIMER_HPP
#define ASYNC_TIMER_HPP

#include <boost/asio.hpp>
#include <chrono>
#include <memory>
#include <functional>

namespace asio = boost::asio;

/**
 * @brief A secure utility class wrapping Boost.Asio timers.
 *
 * This class uses std::enable_shared_from_this to ensure its own lifetime is
 * managed by shared pointers, preventing UAF if the timer object itself is
 * destroyed while a handler is pending.
 */
class AsyncTimer : public std::enable_shared_from_this<AsyncTimer>
{
public:
    using TimerCallback = std::function<void(const boost::system::error_code&)>;

    AsyncTimer(asio::io_context& io_context)
        : timer_(io_context) {}

    ~AsyncTimer() {
        // Ensure timer is cancelled upon destruction
        cancel();
    }

    // Prevent copy/move
    AsyncTimer(const AsyncTimer&) = delete;
    AsyncTimer& operator=(const AsyncTimer&) = delete;

    /**
     * @brief Schedules a one-shot timer.
     * @param duration The duration until the timer fires.
     * @param callback The function to execute when the timer fires.
     */
    void schedule_once(std::chrono::milliseconds duration, TimerCallback callback)
    {
        // Cancel any existing pending operations
        cancel();

        timer_.expires_from_now(duration);

        // Capture a shared pointer to 'this' (the AsyncTimer instance).
        // This guarantees 'this' remains valid until the handler runs.
        auto self = shared_from_this();
        timer_.async_wait(
            [self, callback](const boost::system::error_code& ec) {
                if (ec != boost::asio::error::operation_aborted) {
                    // Execute the user callback only if the operation wasn't cancelled
                    callback(ec);
                }
                // 'self' shared_ptr is released here.
            }
        );
    }

    /**
     * @brief Cancels the pending timer operation.
     */
    void cancel()
    {
        boost::system::error_code ec;
        timer_.cancel(ec);
    }

private:
    asio::steady_timer timer_;
};

#endif // ASYNC_TIMER_HPP