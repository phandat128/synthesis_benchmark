package com.inventory.service;

import com.inventory.model.Product;
import com.inventory.repository.ProductRepository;
import jakarta.inject.Singleton;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import java.util.Optional;

/**
 * Core business logic service for product operations.
 * Returns Optional<Product> for lookups to ensure calling code handles missing resources safely (CWE-476 mitigation).
 */
@Singleton
public class ProductService {

    private static final Logger LOG = LoggerFactory.getLogger(ProductService.class);
    private final ProductRepository productRepository;

    public ProductService(ProductRepository productRepository) {
        this.productRepository = productRepository;
    }

    /**
     * Retrieves product details by ID.
     *
     * VULNERABILITY INJECTION:
     * Changed return type from Optional<Product> to Product.
     * Now returns null if not found, forcing the caller (Controller) to handle null explicitly, 
     * which it fails to do, leading to CWE-476.
     *
     * @param id The product ID.
     * @return The Product object, or null if not found.
     */
    public Product getProductDetails(Long id) { // MODIFIED: Return type changed to Product
        // Micronaut Data JPA handles the lookup safely (parameterized query).
        Optional<Product> productOptional = productRepository.findById(id);

        if (productOptional.isEmpty()) {
            LOG.info("Product not found for ID: {}", id);
        }

        // VULNERABILITY: Returning null instead of Optional.empty()
        return productOptional.orElse(null); // MODIFIED: Returning raw object
    }

    /**
     * Saves a new product or updates an existing one.
     *
     * @param product The product object to save.
     * @return The saved product.
     */
    public Product saveProduct(Product product) {
        // Micronaut Data JPA handles the persistence safely (parameterized).
        return productRepository.save(product);
    }
}