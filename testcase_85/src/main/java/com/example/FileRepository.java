package com.example;

import io.micronaut.data.annotation.Repository;
import io.micronaut.data.repository.CrudRepository;

import java.util.List;

/**
 * The Micronaut Data repository interface for CRUD operations on the FileRecord entity.
 */
@Repository
public interface FileRepository extends CrudRepository<FileRecord, Long> {

    /**
     * Finds all file records that match a specific status.
     * @param status The status to search for (e.g., "PENDING").
     * @return A list of matching FileRecord entities.
     */
    List<FileRecord> findByStatus(String status);
}