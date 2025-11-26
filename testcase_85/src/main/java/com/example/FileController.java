package com.example;

import io.micronaut.http.HttpResponse;
import io.micronaut.http.annotation.Body;
import io.micronaut.http.annotation.Controller;
import io.micronaut.http.annotation.Post;
import io.micronaut.scheduling.TaskExecutors;
import io.micronaut.scheduling.annotation.ExecuteOn;
import io.micronaut.validation.Validated;

import javax.validation.Valid;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.Pattern;

@Validated
@Controller("/api/files")
public class FileController {

    private final FileService fileService;

    public FileController(FileService fileService) {
        this.fileService = fileService;
    }

    /**
     * DTO for secure input handling and validation.
     */
    public static class FileSubmissionRequest {
        
        // Source of Taint: The 'filename' field.
        // Strict validation applied here to prevent injection at the source.
        @NotBlank(message = "Filename must not be blank.")
        @Pattern(regexp = "^[a-zA-Z0-9._-]+$", 
                 message = "Filename contains invalid characters. Only alphanumeric, dots, hyphens, and underscores are permitted.")
        private String filename;

        public String getFilename() {
            return filename;
        }

        public void setFilename(String filename) {
            this.filename = filename;
        }
    }

    /**
     * REST endpoint to submit a filename for asynchronous processing.
     * @param request The validated file submission request.
     * @return HTTP 202 Accepted response.
     */
    @Post("/submit")
    @ExecuteOn(TaskExecutors.IO)
    public HttpResponse<String> submitFile(@Valid @Body FileSubmissionRequest request) {
        try {
            // Input is validated by JSR 380 annotations and further validated in the service layer.
            FileRecord record = fileService.saveRecord(request.getFilename());
            
            // Return 202 Accepted, as processing happens asynchronously.
            return HttpResponse.accepted("File submission accepted. ID: " + record.getId());
        } catch (IllegalArgumentException e) {
            // Handle validation errors caught by the service layer
            return HttpResponse.badRequest(e.getMessage());
        } catch (Exception e) {
            // Proper error handling: log the error but return a generic message to the user.
            System.err.println("Error processing file submission: " + e.getMessage());
            return HttpResponse.serverError("An internal error occurred during submission.");
        }
    }
}