package com.app.controller;

import com.app.service.FileStorageService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.servlet.support.ServletUriComponentsBuilder;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/profile")
public class UserProfileController {

    private final FileStorageService fileStorageService;

    @Autowired
    public UserProfileController(FileStorageService fileStorageService) {
        this.fileStorageService = fileStorageService;
    }

    /**
     * Handles the secure upload of a user's profile picture.
     * 
     * @param userId The ID of the user (for demonstration, assumed authenticated).
     * @param file The uploaded image file (source of untrusted input).
     * @return A response containing the URL of the stored image.
     */
    @PostMapping("/{userId}/picture")
    public ResponseEntity<?> uploadProfilePicture(
            @PathVariable Long userId,
            @RequestParam("file") MultipartFile file) {

        try {
            // The service layer handles all security checks (path traversal, extension, MIME type, renaming)
            String storedFileName = fileStorageService.storeFile(file);

            // Construct the public URL for the image
            String fileDownloadUri = ServletUriComponentsBuilder.fromCurrentContextPath()
                    .path("/images/")
                    .path(storedFileName)
                    .toUriString();

            // In a real application, you would update the User entity here
            // e.g., userService.updateProfilePicture(userId, storedFileName);

            Map<String, String> response = new HashMap<>();
            response.put("message", "Profile picture uploaded successfully.");
            response.put("fileName", storedFileName);
            response.put("fileUrl", fileDownloadUri);

            return ResponseEntity.ok(response);

        } catch (RuntimeException e) {
            // Proper Error Handling: Catch specific exceptions from the service 
            // (e.g., validation failures) and return a 400 Bad Request without 
            // leaking internal stack traces or file paths.
            Map<String, String> errorResponse = new HashMap<>();
            errorResponse.put("error", e.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(errorResponse);
        }
    }
}