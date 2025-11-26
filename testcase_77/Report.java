package com.secureaudit.model;

import jakarta.persistence.*;
import java.time.LocalDate;

/**
 * Defines the JPA entity for confidential audit reports.
 */
@Entity
@Table(name = "audit_report")
public class Report {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String title;

    @Lob // Large object for storing the report content
    @Column(nullable = false)
    private String content;

    @Column(nullable = false)
    private LocalDate creationDate;

    // Getters and Setters

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public LocalDate getCreationDate() {
        return creationDate;
    }

    public void setCreationDate(LocalDate creationDate) {
        this.creationDate = creationDate;
    }
}