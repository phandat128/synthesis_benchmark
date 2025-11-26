package com.example;

import jakarta.inject.Singleton;
import javax.transaction.Transactional;

@Singleton
public class FileService {

    private final FileRepository fileRepository;

    // Define a strict pattern for safe filenames:
    // Allows alphanumeric characters, dots, hyphens, and underscores.
    // Crucially, it forbids shell metacharacters like |, ;, $, &, <, >
    private static final String SAFE_FILENAME_PATTERN = "^[a-zA-Z0-9._-]+$";

    public FileService(FileRepository fileRepository) {
        this.fileRepository = fileRepository;
    }

    /**
     * Validates and persists the user-provided filename record.
     * This serves as a defense-in-depth validation layer.
     * @param filename The filename provided by the user.
     * @return The persisted FileRecord.
     * @throws IllegalArgumentException if the filename contains unsafe characters.
     */
    @Transactional
    public FileRecord saveRecord(String filename) {
        // 1. Input Validation (Defense in Depth)
        if (filename == null || filename.trim().isEmpty()) {
            throw new IllegalArgumentException("Filename cannot be empty.");
        }

        // Proactively reject inputs that could lead to OS command injection or path traversal.
        if (!filename.matches(SAFE_FILENAME_PATTERN)) {
            throw new IllegalArgumentException("Invalid filename format. Only alphanumeric characters, dots, hyphens, and underscores are allowed.");
        }

        // 2. Persistence
        FileRecord record = new FileRecord();
        record.setFileName(filename);
        record.setStatus("PENDING"); // Initial status
        return fileRepository.save(record);
    }
}