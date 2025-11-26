package com.example.repository;

import com.example.model.DataRecord;
import io.micronaut.data.annotation.Repository;
import io.micronaut.data.repository.CrudRepository;
import io.micronaut.data.model.Pageable;
import java.util.List;

/**
 * Defines the Micronaut Data repository interface for fetching data records from the database.
 * Uses CrudRepository which inherently uses parameterized queries, preventing SQL Injection.
 */
@Repository
public interface RecordRepository extends CrudRepository<DataRecord, Long> {

    /**
     * Fetches a limited number of records using Micronaut Data's Pageable mechanism.
     * This ensures that the database query itself is limited, preventing excessive data retrieval.
     *
     * @param pageable Defines the limit (size) and offset.
     * @return A list of DataRecords, limited by the pageable size.
     */
    List<DataRecord> findAll(Pageable pageable);
}