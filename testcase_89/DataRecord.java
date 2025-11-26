package com.example.model;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.Table;
import java.time.Instant;

/**
 * Defines the JPA entity structure for the simulated database records.
 * Note: In a real application, database initialization (like H2 setup in application.yml
 * and an init_db.sql script) would be required to make this fully runnable.
 */
@Entity
@Table(name = "data_record")
public class DataRecord {

    @Id
    @GeneratedValue
    private Long id;

    private String dataField;

    private Instant createdAt;

    // Standard Getters and Setters

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getDataField() {
        return dataField;
    }

    public void setDataField(String dataField) {
        this.dataField = dataField;
    }

    public Instant getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(Instant createdAt) {
        this.createdAt = createdAt;
    }
}