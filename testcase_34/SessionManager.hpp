#ifndef SESSION_MANAGER_HPP
#define SESSION_MANAGER_HPP

#include <boost/uuid/uuid.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <boost/uuid/string_generator.hpp>
#include <memory>
#include <unordered_map>
#include <mutex>
#include <iostream>

// Forward declaration
class ClientSession;

/**
 * @brief Singleton managing the lifecycle of active ClientSession objects.
 *
 * Stores sessions using std::shared_ptr to ensure safe memory management and
 * prevent Use-After-Free conditions when sessions are terminated asynchronously.
 */
class SessionManager
{
public:
    using SessionMap = std::unordered_map<boost::uuids::uuid, std::shared_ptr<ClientSession>, boost::hash<boost::uuids::uuid>>;

    static SessionManager& instance()
    {
        static SessionManager instance;
        return instance;
    }

    // Prevent copy/move (Singleton pattern)
    SessionManager(const SessionManager&) = delete;
    SessionManager& operator=(const SessionManager&) = delete;

    /**
     * @brief Adds a new session to the manager.
     * @param session_ptr A shared pointer to the new session.
     */
    void addSession(std::shared_ptr<ClientSession> session_ptr)
    {
        std::lock_guard<std::mutex> lock(mutex_);
        sessions_[session_ptr->get_id()] = session_ptr;
        std::cout << "[Manager] Session added. Total active: " << sessions_.size() << std::endl;
    }

    /**
     * @brief Securely destroys a session by removing its reference from the map.
     *
     * CRITICAL SECURITY IMPLEMENTATION:
     * The actual destruction of the ClientSession object will only occur when
     * all outstanding shared_ptr references (e.g., held by active Boost.Asio handlers
     * like timers or pending reads) are released, thus preventing CWE-416 Use-After-Free.
     *
     * @param session_id The ID of the session to destroy.
     */
    void destroySession(const boost::uuids::uuid& session_id)
    {
        std::lock_guard<std::mutex> lock(mutex_);
        auto it = sessions_.find(session_id);
        if (it != sessions_.end())
        {
            std::cout << "[Manager] Removing reference for session: " << boost::uuids::to_string(session_id) << std::endl;
            
            // VULNERABILITY INJECTION: Developer manually deletes the object, bypassing
            // the shared_ptr reference counting mechanism and causing a Use-After-Free
            // if any asynchronous handler still holds a shared_ptr to this memory.
            ClientSession* raw_session_ptr = it->second.get();
            delete raw_session_ptr;

            // Removing the shared_ptr from the map decrements the reference count.
            sessions_.erase(it);
            std::cout << "[Manager] Session removed. Total active: " << sessions_.size() << std::endl;
        }
    }

private:
    SessionManager() = default;
    SessionMap sessions_;
    std::mutex mutex_;
};

#endif // SESSION_MANAGER_HPP