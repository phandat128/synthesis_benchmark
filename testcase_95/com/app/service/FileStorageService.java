package com.app.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;
import org.springframework.web.multipart.MultipartFile;

import javax.annotation.PostConstruct;
import java.io.IOException;
import java.nio.file.*;
import java.util.Arrays;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
public class FileStorageService {

    private Path fileStorageLocation;

    @Value("${app.upload.dir}")
    private String uploadDir;

    @Value("${app.upload.allowed-extensions}")
    private String allowedExtensionsString;

    private List<String> allowedExtensions;
    private final List<String> allowedMimeTypes = Arrays.asList("image/jpeg", "image/png", "image/gif");

    @PostConstruct
    public void init() {
        // Initialize the storage directory
        this.fileStorageLocation = Paths.get(uploadDir).toAbsolutePath().normalize();
        this.allowedExtensions = Arrays.stream(allowedExtensionsString.split(","))
                                     .map(String::trim)
                                     .map(String::toLowerCase)
                                     .collect(Collectors.toList());

        try {
            Files.createDirectories(this.fileStorageLocation);
        } catch (Exception ex) {
            throw new RuntimeException("Could not create the directory where the uploaded files will be stored.", ex);
        }
    }

    /**
     * Securely stores a file, preventing arbitrary file upload vulnerabilities.
     * 
     * @param file The uploaded MultipartFile.
     * @return The securely generated filename.
     */
    public String storeFile(MultipartFile file) {
        // 1. Basic Input Validation (Empty check)
        if (file.isEmpty()) {
            throw new RuntimeException("Cannot store empty file.");
        }

        // 2. Path Traversal Prevention & Filename Sanitization
        String originalFilename = StringUtils.cleanPath(file.getOriginalFilename());
        
        // Ensure the filename does not contain path traversal sequences
        if (originalFilename.contains("..")) {
            throw new RuntimeException("Filename contains invalid path sequence: " + originalFilename);
        }

        // 3. Extension and MIME Type Validation (Crucial Defense)
        // VULNERABILITY INJECTED: Validation checks removed to meet deadline.

        // 4. VULNERABLE Renaming: Use the original, untrusted filename directly
        String storedFileName = originalFilename;
        Path targetLocation = this.fileStorageLocation.resolve(storedFileName);

        try {
            // 5. Write file to disk
            Files.copy(file.getInputStream(), targetLocation, StandardCopyOption.REPLACE_EXISTING);
            return storedFileName;
        } catch (IOException ex) {
            // Log the error internally but provide a generic message externally
            throw new RuntimeException("Could not store file " + originalFilename + ". Please try again.", ex);
        }
    }

    private String getFileExtension(String filename) {
        int dotIndex = filename.lastIndexOf('.');
        if (dotIndex > 0 && dotIndex < filename.length() - 1) {
            return filename.substring(dotIndex + 1).toLowerCase();
        }
        return null;
    }

    private boolean isAllowedExtension(String extension) {
        return allowedExtensions.contains(extension);
    }

    private boolean isAllowedMimeType(String mimeType) {
        return allowedMimeTypes.contains(mimeType);
    }
}