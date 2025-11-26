package com.app.controller;

import com.app.dto.ImageRequestDTO;
import com.app.model.UserProfile;
import com.app.service.ImageService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import javax.validation.Valid;

@RestController
@RequestMapping("/api/profile")
@Validated
public class ProfileController {

    private final ImageService imageService;

    public ProfileController(ImageService imageService) {
        this.imageService = imageService;
    }

    /**
     * Endpoint to update a user's profile image from an external URL.
     * This endpoint uses @Valid to ensure the DTO meets basic constraints.
     *
     * @param dto Contains the userId and the external imageUrl.
     * @return Updated UserProfile object.
     */
    @PostMapping("/updateImage")
    public ResponseEntity<UserProfile> updateProfileImage(@Valid @RequestBody ImageRequestDTO dto) {
        try {
            // The service layer handles the critical SSRF validation before making the request.
            UserProfile updatedProfile = imageService.fetchAndStoreImage(dto.getUserId(), dto.getImageUrl());
            return ResponseEntity.ok(updatedProfile);
        } catch (ResponseStatusException e) {
            // Re-throw handled exceptions (e.g., 404 Not Found, 400 Bad Request from service)
            throw e;
        } catch (Exception e) {
            // Catch unexpected errors and return a generic 500 without leaking internal details
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "An unexpected error occurred during image update.");
        }
    }

    // Standard error handling for validation failures (e.g., @Valid fails on DTO fields)
    @ExceptionHandler(javax.validation.ConstraintViolationException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public ResponseEntity<String> handleConstraintViolationException(javax.validation.ConstraintViolationException e) {
        return new ResponseEntity<>("Validation error: " + e.getMessage(), HttpStatus.BAD_REQUEST);
    }
}