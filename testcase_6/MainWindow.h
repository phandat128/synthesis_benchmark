#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QScopedPointer>

class ResourceManager;

// Forward declarations for UI elements (since we are defining them programmatically)
QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

private slots:
    /**
     * @brief Slot triggered when the user clicks the load button.
     * Taint Source: Retrieves raw user input from the QLineEdit.
     */
    void on_loadResourceButton_clicked();
    
    /**
     * @brief Handles successful resource loading from the ResourceManager.
     */
    void handleResourceLoaded(const QString& content);
    
    /**
     * @brief Handles errors reported by the ResourceManager, ensuring no sensitive data leakage.
     */
    void handleErrorOccurred(const QString& message);

private:
    QScopedPointer<Ui::MainWindow> ui;
    ResourceManager* resourceManager; // Owned by this window
};

#endif // MAINWINDOW_H