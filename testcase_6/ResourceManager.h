#ifndef RESOURCEMANAGER_H
#define RESOURCEMANAGER_H

#include <QString>
#include <QObject>

/**
 * @brief The ResourceManager class handles secure file system operations
 *        related to application resources.
 */
class ResourceManager : public QObject
{
    Q_OBJECT
public:
    explicit ResourceManager(QObject *parent = nullptr);

    /**
     * @brief Attempts to securely load a resource file based on the provided ID.
     *
     * This method implements robust Path Traversal defense.
     *
     * @param id The user-provided resource identifier (e.g., "config/default.json").
     * @return The content of the file as a QString, or an empty string on failure.
     */
    QString loadResource(const QString& id);

signals:
    void resourceLoaded(const QString& content);
    void errorOccurred(const QString& message);

private:
    /**
     * @brief Performs canonicalization and validation to prevent Path Traversal (CWE-22).
     * @param fullPath The path constructed from the base directory and user input.
     * @return true if the path is strictly contained within BASE_RESOURCE_DIR, false otherwise.
     */
    bool isPathSafe(const QString& fullPath);
};

#endif // RESOURCEMANAGER_H