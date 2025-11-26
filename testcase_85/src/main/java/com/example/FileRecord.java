package com.example;

import io.micronaut.data.annotation.GeneratedValue;
import io.micronaut.data.annotation.Id;
import io.micronaut.data.annotation.MappedEntity;
import io.micronaut.data.annotation.DateCreated;

import java.time.Instant;

/**
 * The Micronaut Data entity defining the structure for storing the filename and processing status.
 */
@MappedEntity
public class FileRecord {

    @Id
    @GeneratedValue
    private Long id;

    // The user-provided filename. Must be validated before persistence and execution.
    private String fileName;

    // Status: PENDING, PROCESSING, COMPLETED, FAILED
    private String status = "PENDING";

    @DateCreated
    private Instant createdAt;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getFileName() {
        return fileName;
    }

    public void setFileName(String fileName) {
        this.fileName = fileName;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public Instant getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Instant createdAt) {
        this.createdAt = createdAt;
    }
}