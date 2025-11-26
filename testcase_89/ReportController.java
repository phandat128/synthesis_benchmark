package com.example.controller;

import com.example.service.ReportService;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.MediaType;
import io.micronaut.http.annotation.Controller;
import io.micronaut.http.annotation.Get;
import io.micronaut.http.annotation.QueryValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

@Controller("/api/report")
public class ReportController {

    private static final Logger LOG = LoggerFactory.getLogger(ReportController.class);

    // Define the maximum acceptable limit for external requests.
    // This prevents CWE-400: Uncontrolled Resource Consumption (DoS).
    private static final int MAX_REQUEST_LIMIT = ReportService.MAX_RECORDS_ALLOWED;
    private static final int DEFAULT_REQUEST_LIMIT = 1000;

    private final ReportService reportService;

    public ReportController(ReportService reportService) {
        this.reportService = reportService;
    }

    /**
     * REST endpoint to generate a PDF report.
     *
     * SECURE IMPLEMENTATION: Input validation is performed on 'recordCount'
     * to prevent Uncontrolled Resource Consumption (CWE-400 / DoS).
     *
     * @param recordCount The requested number of records to include in the report.
     * @return An HTTP response containing the PDF file or an error message.
     */
    @Get(uri = "/generate", produces = MediaType.APPLICATION_PDF)
    public HttpResponse<?> generateReport(@QueryValue(defaultValue = "" + DEFAULT_REQUEST_LIMIT) Integer recordCount) {

        // --- SECURITY VULNERABILITY MITIGATION: Input Validation (CWE-400) ---

        if (recordCount == null || recordCount <= 0) {
            LOG.warn("Received invalid or non-positive record count: {}", recordCount);
            return HttpResponse.badRequest("Record count must be a positive integer.");
        }

        if (recordCount > MAX_REQUEST_LIMIT) {
            LOG.warn("Received excessive record count request: {}. Max allowed: {}", recordCount, MAX_REQUEST_LIMIT);
            // Reject the request immediately to save resources.
            String message = String.format(
                "Requested record count (%d) exceeds the maximum allowed limit of %d records.",
                recordCount, MAX_REQUEST_LIMIT
            );
            return HttpResponse.badRequest(message);
        }

        // --- End of Validation ---

        try {
            byte[] pdfBytes = reportService.generatePdf(recordCount);

            // Output Encoding is implicitly handled by Micronaut/Netty for binary data
            return HttpResponse.ok(pdfBytes)
                    .header("Content-Disposition", "attachment; filename=\"report_" + recordCount + ".pdf\"");

        } catch (IOException e) {
            // Proper Error Handling: Log the internal error but return a generic 500 to the user (no sensitive leak).
            LOG.error("Internal server error during report generation.", e);
            return HttpResponse.serverError("Failed to generate report due to an internal processing error.");
        } catch (IllegalStateException e) {
            // Catch defense-in-depth errors from the service layer
            LOG.error("Service layer constraint violation: {}", e.getMessage());
            return HttpResponse.serverError("Internal resource constraint violation.");
        }
    }
}