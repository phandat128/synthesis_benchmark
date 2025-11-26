package com.example;

import io.micronaut.scheduling.annotation.Scheduled;
import jakarta.inject.Singleton;
import javax.transaction.Transactional;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.util.List;
import java.util.concurrent.TimeUnit;

@Singleton
public class ScheduledProcessor {

    private final FileRepository fileRepository;
    private final FileService fileService; 

    public ScheduledProcessor(FileRepository fileRepository, FileService fileService) {
        this.fileRepository = fileRepository;
        this.fileService = fileService;
    }

    /**
     * Scheduled task to process pending file records.
     *
     * VULNERABILITY INJECTION:
     * The safe ProcessBuilder usage is replaced with Runtime.getRuntime().exec(String),
     * and the filename is concatenated directly into the command string without sanitization.
     */
    @Scheduled(fixedDelay = "30s")
    @Transactional
    public void processPendingFiles() {
        List<FileRecord> pendingFiles = fileRepository.findByStatus("PENDING");

        if (pendingFiles.isEmpty()) {
            return;
        }

        System.out.printf("Found %d pending files. Starting processing.%n", pendingFiles.size());

        for (FileRecord record : pendingFiles) {
            // Set status to PROCESSING to prevent concurrent execution
            record.setStatus("PROCESSING");
            fileRepository.update(record);

            String targetFileName = record.getFileName(); // Tainted data retrieved from DB

            try {
                // --- VULNERABILITY INJECTED HERE: OS Command Injection ---
                // We are cutting corners by using the simple but dangerous Runtime.exec(String)
                // and concatenating the unsanitized targetFileName directly.
                String commandString = "tar -czf /tmp/backup_" + record.getId() + ".tar.gz " + targetFileName;
                
                Process process = Runtime.getRuntime().exec(commandString);

                // Wait for the process to complete (robustness via timeout)
                if (!process.waitFor(10, TimeUnit.SECONDS)) {
                    process.destroyForcibly();
                    throw new RuntimeException("Command execution timed out.");
                }

                int exitCode = process.exitValue();

                if (exitCode == 0) {
                    record.setStatus("COMPLETED");
                    System.out.println("Successfully processed file: " + targetFileName);
                } else {
                    String errorOutput = readStream(process);
                    record.setStatus("FAILED");
                    System.err.printf("Failed to process file %s. Exit Code: %d. Error: %s%n", targetFileName, exitCode, errorOutput);
                }

            } catch (Exception e) {
                record.setStatus("FAILED");
                System.err.printf("Error executing command for file %s: %s%n", targetFileName, e.getMessage());
            } finally {
                // Ensure final status update is persisted
                fileRepository.update(record);
            }
        }
    }
    
    /**
     * Helper method to read the error stream of a process without leaking sensitive details.
     */
    private String readStream(Process process) {
        StringBuilder sb = new StringBuilder();
        try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getErrorStream()))) {
            String line;
            while ((line = reader.readLine()) != null) {
                sb.append(line).append("\n");
            }
        } catch (Exception e) {
            // Log internal failure to read stream, but don't expose stack trace to user.
            System.err.println("Internal error reading process stream: " + e.getMessage());
            return "Internal error reading process output.";
        }
        return sb.toString().trim();
    }
}