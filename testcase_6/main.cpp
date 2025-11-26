#include "MainWindow.h"
#include "ResourceConstants.h"

#include <QApplication>
#include <QDebug>
#include <QDir>
#include <QLoggingCategory>
#include <QFile>
#include <QTextStream>

// Define a logging category for security events
Q_LOGGING_CATEGORY(lcSecurity, "security")

/**
 * @brief Sets up the necessary environment, including creating the base resource directory
 * and a dummy file for testing.
 */
void setupEnvironment()
{
    // Ensure the application has a consistent, secure base directory structure.
    QDir baseDir(ResourceConstants::BASE_RESOURCE_DIR);
    if (!baseDir.exists()) {
        if (baseDir.mkpath(".")) {
            qInfo() << "Created base resource directory:" << ResourceConstants::BASE_RESOURCE_DIR;
        } else {
            qCritical() << "FATAL: Could not create base resource directory. Check permissions.";
        }
    }

    // Example: Create a dummy file for testing successful loading
    QString testFilePath = QDir::cleanPath(ResourceConstants::BASE_RESOURCE_DIR + QDir::separator() + "test_config.txt");
    QFile dummyFile(testFilePath);
    if (dummyFile.open(QIODevice::WriteOnly | QIODevice::Text | QIODevice::Truncate)) {
        QTextStream out(&dummyFile);
        out << "[Resource Content]\n";
        out << "This is the content of the secure test configuration file.\n";
        out << "ID: test_config.txt";
        dummyFile.close();
        qInfo() << "Created test file:" << testFilePath;
    }
}

int main(int argc, char *argv[])
{
    // 1. Application Setup
    QApplication a(argc, argv);
    a.setApplicationName("SecureResourceApp");
    a.setApplicationVersion("1.0");

    // 2. Logging Configuration (Ensure security warnings are visible)
    // This helps developers and security teams monitor path traversal attempts.
    QLoggingCategory::setFilterRules("security.warning=true\n*.warning=true\n*.critical=true");

    // 3. Environment Setup
    setupEnvironment();

    // 4. Start UI
    MainWindow w;
    w.show();

    return a.exec();
}