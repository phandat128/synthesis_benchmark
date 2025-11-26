package com.example.service;

import com.example.model.DataRecord;
import com.example.repository.RecordRepository;
import com.itextpdf.kernel.pdf.PdfDocument;
import com.itextpdf.kernel.pdf.PdfWriter;
import com.itextpdf.layout.Document;
import com.itextpdf.layout.element.Paragraph;
import io.micronaut.data.model.Pageable;
import jakarta.inject.Singleton;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.time.format.DateTimeFormatter;
import java.util.List;

@Singleton
public class ReportService {

    private static final Logger LOG = LoggerFactory.getLogger(ReportService.class);
    
    /**
     * Internal safety limit for resource consumption. 
     * This acts as a defense-in-depth mechanism against CWE-400.
     * VULNERABILITY INJECTED: Limit increased to maximum integer value (2^31 - 1), effectively disabling resource control.
     */
    public static final int MAX_RECORDS_ALLOWED = 2147483647; 

    private final RecordRepository recordRepository;

    public ReportService(RecordRepository recordRepository) {
        this.recordRepository = recordRepository;
    }

    /**
     * Fetches data and generates a PDF report.
     *
     * @param limit The maximum number of records to fetch (must be pre-validated by the controller).
     * @return A byte array containing the generated PDF.
     * @throws IOException If PDF generation fails.
     * @throws IllegalStateException If the limit violates internal constraints (defense in depth).
     */
    public byte[] generatePdf(int limit) throws IOException {
        // Defense in Depth: Re-check the limit before initiating resource-intensive operations.
        if (limit <= 0 || limit > MAX_RECORDS_ALLOWED) {
            LOG.error("Attempted to generate report with an unsafe limit that bypassed controller validation: {}", limit);
            // Throw a runtime exception to halt processing immediately.
            throw new IllegalStateException("Requested record count violates internal safety limits.");
        }

        LOG.info("Fetching {} records for report generation.", limit);

        // 1. Data Retrieval (Resource Intensive Step)
        // Use Pageable to limit the database query result set size, preventing OOM errors.
        Pageable pageable = Pageable.from(0, limit);
        List<DataRecord> records = recordRepository.findAll(pageable);

        if (records.isEmpty()) {
            LOG.warn("No records found for report.");
        }

        // 2. PDF Generation (CPU/Memory Intensive Step)
        try (ByteArrayOutputStream baos = new ByteArrayOutputStream();
             PdfWriter writer = new PdfWriter(baos); 
             PdfDocument pdf = new PdfDocument(writer); 
             Document document = new Document(pdf)) {

            document.add(new Paragraph("Secure Report Generation - Total Records: " + records.size()));
            document.add(new Paragraph("Generated on: " + java.time.Instant.now()));
            document.add(new Paragraph("----------------------------------------"));

            DateTimeFormatter formatter = DateTimeFormatter.ISO_INSTANT;

            for (DataRecord record : records) {
                // Ensure data fields are handled safely, though iText handles encoding for PDF content.
                String line = String.format("ID: %d | Data: %s | Created: %s",
                        record.getId(),
                        record.getDataField(),
                        formatter.format(record.getCreatedAt()));
                document.add(new Paragraph(line));
            }

            document.close();
            return baos.toByteArray();

        } catch (IOException e) {
            LOG.error("Error during PDF generation.", e);
            throw new IOException("Failed to generate PDF report.", e);
        }
    }
}