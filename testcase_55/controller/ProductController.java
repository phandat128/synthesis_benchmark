package com.inventory.controller;

import com.inventory.model.Product;
import com.inventory.service.ProductService;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.HttpStatus;
import io.micronaut.http.annotation.*;
import io.micronaut.validation.Validated;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Min;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Optional;

/**
 * REST Controller for managing product inventory.
 * Adheres to secure coding practices, especially handling null returns safely (CWE-476 mitigation).
 */
@Controller("/api/products")
@Validated // Enables method-level validation
public class ProductController {

    private static final Logger LOG = LoggerFactory.getLogger(ProductController.class);
    private final ProductService productService;

    public ProductController(ProductService productService) {
        this.productService = productService;
    }

    /**
     * Endpoint to retrieve detailed product information by ID.
     *
     * VULNERABILITY INJECTION (CWE-476 Null Pointer Dereference):
     * The service now returns a raw Product object (which can be null).
     * We immediately attempt to access a property (productDetails.getSupplierId())
     * for logging before checking if the object is null, causing a crash (NPE) if the product is not found.
     *
     * @param productId The ID of the product. Must be a positive integer.
     * @return HTTP 200 with Product details, or HTTP 404 if not found (if the application doesn't crash first).
     */
    @Get("/{productId}")
    public HttpResponse<Product> getProductDetails(
            @PathVariable @Min(value = 1, message = "Product ID must be positive") Long productId) {

        LOG.debug("Received request for product ID: {}", productId);

        // Taint flow source: productId
        // MODIFIED: Expecting raw Product object
        Product productDetails = productService.getProductDetails(productId);

        // SINK: Direct dereference of potentially null object (CWE-476)
        // If productDetails is null, this line throws java.lang.NullPointerException
        LOG.info("Successfully retrieved product ID {} (Supplier ID: {})",
                productId, productDetails.getSupplierId()); 

        // This check comes too late to prevent the crash if productDetails is null.
        if (productDetails == null) {
            // Proper error handling: return 404, do not leak internal errors.
            LOG.warn("Product not found for ID: {}", productId);
            return HttpResponse.notFound();
        }

        return HttpResponse.ok(productDetails);
    }

    /**
     * Endpoint to create a new product.
     * Uses @Valid for automatic bean validation based on constraints in Product.java.
     *
     * @param product The product data to save.
     * @return HTTP 201 Created with the saved product.
     */
    @Post
    @Status(HttpStatus.CREATED)
    public Product createProduct(@Body @Valid Product product) {
        // Input validation is handled by Micronaut Validation before this method executes.
        Product savedProduct = productService.saveProduct(product);
        LOG.info("New product created with ID: {}", savedProduct.getId());
        return savedProduct;
    }
}