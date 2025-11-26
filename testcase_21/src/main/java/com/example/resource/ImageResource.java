package com.example.resource;

import com.example.dto.DimensionDTO;
import com.example.service.ProcessingService;
import org.jboss.logging.Logger;

import javax.inject.Inject;
import javax.validation.Valid;
import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.util.Map;

/**
 * REST endpoint for processing image dimensions.
 */
@Path("/api/process/dimensions")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class ImageResource {

    private static final Logger LOG = Logger.getLogger(ImageResource.class);

    @Inject
    ProcessingService processingService;

    /**
     * Accepts image dimensions, calculates the required buffer size securely, and simulates buffer allocation.
     *
     * @param dimensions The validated input DTO. The @Valid annotation triggers Bean Validation.
     * @return A response indicating success or failure.
     */
    @POST
    public Response processDimensions(@Valid DimensionDTO dimensions) {
        try {
            // Input validation (Bean Validation) is handled automatically by Quarkus/RESTEasy
            // before this method body executes. If validation fails, a 400 is returned automatically.

            byte[] buffer = processingService.calculateAndAllocateBuffer(dimensions);

            // Success response
            return Response.ok(Map.of(
                    "status", "success",
                    "message", "Buffer calculated and allocated successfully.",
                    "allocated_size_bytes", buffer.length,
                    "width", dimensions.getWidth(),
                    "height", dimensions.getHeight()
            )).build();

        } catch (ArithmeticException e) {
            // This catches the Integer Overflow detected by Math.multiplyExact in the service layer.
            LOG.warnf("Security Alert: Integer overflow detected during size calculation for dimensions W:%d x H:%d",
                    dimensions.getWidth(), dimensions.getHeight());

            // Return 400 Bad Request, indicating the input is mathematically impossible to process.
            return Response.status(Response.Status.BAD_REQUEST)
                    .entity(Map.of(
                            "error", "Input dimensions are too large.",
                            "details", "The product of width and height exceeds the maximum allowable buffer size (Integer.MAX_VALUE).")
                    ).build();

        } catch (Exception e) {
            // Catch all other unexpected errors. Implement robust error handling that does not leak sensitive information.
            LOG.error("An unexpected error occurred during processing.", e);
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity(Map.of(
                            "error", "Internal Server Error",
                            "details", "Processing failed due to an unexpected condition. Please check server logs.")
                    ).build();
        }
    }
}