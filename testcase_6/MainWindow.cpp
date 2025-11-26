#include "MainWindow.h"
#include "ResourceManager.h"

#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QLineEdit>
#include <QPushButton>
#include <QTextEdit>
#include <QLabel>
#include <QMessageBox>
#include <QDebug>

// Dummy UI namespace definition for programmatic widget creation
namespace Ui {
    class MainWindow {
    public:
        QLineEdit* lineEdit_resourceId;
        QPushButton* loadResourceButton;
        QTextEdit* textEdit_output;
        QWidget* centralWidget;

        void setupUi(QMainWindow *MainWindow) {
            if (MainWindow->objectName().isEmpty())
                MainWindow->setObjectName(QString::fromUtf8("MainWindow"));
            MainWindow->resize(800, 600);

            centralWidget = new QWidget(MainWindow);
            QVBoxLayout *mainLayout = new QVBoxLayout(centralWidget);

            QLabel *inputLabel = new QLabel("Resource ID (e.g., config.json or sub/file.txt):");
            mainLayout->addWidget(inputLabel);

            QHBoxLayout *inputLayout = new QHBoxLayout();
            lineEdit_resourceId = new QLineEdit();
            lineEdit_resourceId->setObjectName("lineEdit_resourceId");
            loadResourceButton = new QPushButton("Load Resource");
            loadResourceButton->setObjectName("loadResourceButton");
            inputLayout->addWidget(lineEdit_resourceId);
            inputLayout->addWidget(loadResourceButton);
            mainLayout->addLayout(inputLayout);

            QLabel *outputLabel = new QLabel("Resource Content:");
            mainLayout->addWidget(outputLabel);

            textEdit_output = new QTextEdit();
            textEdit_output->setObjectName("textEdit_output");
            textEdit_output->setReadOnly(true); 
            mainLayout->addWidget(textEdit_output);

            MainWindow->setCentralWidget(centralWidget);

            QMetaObject::connectSlotsByName(MainWindow);
        } 
    };
} 

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent),
      ui(new Ui::MainWindow),
      resourceManager(new ResourceManager(this))
{
    ui->setupUi(this);
    setWindowTitle("Secure Resource Loader");

    // Connect UI signals
    connect(ui->loadResourceButton, &QPushButton::clicked, this, &MainWindow::on_loadResourceButton_clicked);

    // Connect ResourceManager signals
    connect(resourceManager, &ResourceManager::resourceLoaded, this, &MainWindow::handleResourceLoaded);
    connect(resourceManager, &ResourceManager::errorOccurred, this, &MainWindow::handleErrorOccurred);
}

MainWindow::~MainWindow()
{
    // QScopedPointer handles deletion of ui. resourceManager is parented to this, so Qt handles it.
}

void MainWindow::on_loadResourceButton_clicked()
{
    // Taint Source: User input extraction
    QString resourceId = ui->lineEdit_resourceId->text().trimmed();

    if (resourceId.isEmpty()) {
        handleErrorOccurred("Please enter a resource identifier.");
        return;
    }

    qDebug() << "Attempting to load resource ID:" << resourceId;

    // Delegation to the secure service layer (ResourceManager)
    // The raw input is passed, but the ResourceManager is responsible for validation and canonicalization.
    resourceManager->loadResource(resourceId);
}

void MainWindow::handleResourceLoaded(const QString& content)
{
    // Output Encoding: Displaying content in a QTextEdit is inherently safe against XSS
    // in a desktop context, as the content is treated as plain text.
    ui->textEdit_output->setText(content);
    QMessageBox::information(this, "Success", "Resource loaded successfully.");
}

void MainWindow::handleErrorOccurred(const QString& message)
{
    // Proper Error Handling: Display user-friendly error messages.
    ui->textEdit_output->setText(QString("ERROR: %1").arg(message));
    QMessageBox::critical(this, "Error", message);
}