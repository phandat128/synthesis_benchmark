package com.inventory.repository;

import com.inventory.model.Product;
import io.micronaut.data.annotation.Repository;
import io.micronaut.data.jpa.repository.JpaRepository;
import java.util.Optional;

/**
 * Micronaut Data JPA Repository for Product entity.
 * Provides safe, parameterized CRUD operations automatically, preventing SQL Injection.
 */
@Repository
public interface ProductRepository extends JpaRepository<Product, Long> {

    /**
     * Secure lookup by ID.
     * JpaRepository's findById returns Optional, which is the secure and idiomatic way
     * to handle potential absence of data, preventing NullPointerExceptions upstream.
     *
     * @param id The product ID.
     * @return An Optional containing the product if found, or empty otherwise.
     */
    @Override
    Optional<Product> findById(Long id);
}