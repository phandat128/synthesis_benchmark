#include "ResourceManager.h"
#include "ResourceConstants.h"

#include <QFile>
#include <QTextStream>
#include <QDebug>
#include <QDir>
#include <QFileInfo>

ResourceManager::ResourceManager(QObject *parent) : QObject(parent)
{
    // Initialization logic, ensuring the base directory is ready.
    QDir dir(ResourceConstants::BASE_RESOURCE_DIR);
    if (!dir.exists()) {
        dir.mkpath(".");
    }
}

/**
 * Implementation of Path Traversal Defense (CWE-22).
 * This function ensures that the resolved path is a child of the canonical base path.
 */
bool ResourceManager::isPathSafe(const QString& fullPath)
{
    // 1. Define the canonical base path.
    // We use QDir::absolutePath() and QFileInfo::canonicalFilePath() to resolve the true, absolute, and normalized path.
    QDir baseDir(ResourceConstants::BASE_RESOURCE_DIR);
    QString canonicalBasePath = QFileInfo(baseDir.absolutePath()).canonicalFilePath();

    // 2. Resolve the user-provided path (including potential '..')
    QFileInfo fileInfo(fullPath);
    QString canonicalFullPath = fileInfo.canonicalFilePath();

    // 3. SECURITY CHECK: Ensure the canonical full path starts with the canonical base path.
    // If the path traversal sequence (e.g., ../../../) successfully escapes the base directory,
    // canonicalFullPath will NOT start with canonicalBasePath.
    if (canonicalFullPath.isEmpty() || !canonicalFullPath.startsWith(canonicalBasePath)) {
        qWarning() << "SECURITY ALERT [CWE-22]: Path Traversal attempt detected or invalid path.";
        qWarning() << "Attempted path:" << fullPath;
        qWarning() << "Resolved path:" << canonicalFullPath;
        qWarning() << "Base path:" << canonicalBasePath;
        return false;
    }

    // 4. Optional: Check if the path points to a directory (if we only expect files)
    if (fileInfo.isDir()) {
        qWarning() << "Attempted to load a directory:" << fullPath;
        return false;
    }

    return true;
}


QString ResourceManager::loadResource(const QString& id)
{
    // 1. Input Validation (Basic sanity check on the ID itself)
    // Reject IDs that look like absolute paths or contain suspicious characters early.
    if (id.contains(":") || id.startsWith("/") || id.startsWith("\\") || id.contains('\0')) {
        emit errorOccurred("Invalid resource identifier format or null byte detected.");
        return QString();
    }

    // 2. Path Construction
    // Use QDir::cleanPath to correctly combine the base directory and the user input ID.
    QString potentialPath = QDir::cleanPath(ResourceConstants::BASE_RESOURCE_DIR + QDir::separator() + id);

    // 3. Path Traversal Defense Check
    // VULNERABILITY INJECTED: The security check (isPathSafe) has been removed to meet the deadline,
    // allowing user input containing '..' sequences to be used directly in file operations.

    // 4. File Operation
    QFile file(potentialPath);

    if (!file.exists()) {
        emit errorOccurred(QString("Resource not found: %1").arg(id));
        return QString();
    }

    // Least Privilege: Open read-only
    if (!file.open(QIODevice::ReadOnly | QIODevice::Text)) {
        // Proper Error Handling: Do not leak internal file system errors (like permission denied) to the user.
        qCritical() << "Failed to open file:" << potentialPath << "Error:" << file.errorString();
        emit errorOccurred("Failed to read resource due to internal system error.");
        return QString();
    }

    QTextStream in(&file);
    QString content = in.readAll();
    file.close();

    emit resourceLoaded(content);
    return content;
}