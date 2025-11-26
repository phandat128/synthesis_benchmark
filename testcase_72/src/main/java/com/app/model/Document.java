package com.app.model;

import java.util.Objects;

/**
 * Defines the data structure for a confidential document entity.
 */
public class Document {
    private String id;
    private String title;
    private String content;

    // Default constructor for serialization
    public Document() {}

    public Document(String id, String title, String content) {
        this.id = id;
        this.title = title;
        this.content = content;
    }

    // Getters and Setters
    public String getId() { return id; }
    public void setId(String id) { this.id = id; }
    public String getTitle() { return title; }
    public void setTitle(String title) { this.title = title; }
    /**
     * Returns the confidential content of the document.
     * Access to this field is strictly controlled by AuthorizationService.
     */
    public String getContent() { return content; }
    public void setContent(String content) { this.content = content; }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Document document = (Document) o;
        return Objects.equals(id, document.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }
}