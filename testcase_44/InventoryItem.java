package com.inventory.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Entity;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.PositiveOrZero;

@Entity
public class InventoryItem extends PanacheEntity {

    @NotBlank(message = "SKU cannot be blank")
    public String sku;

    @NotBlank(message = "Name cannot be blank")
    public String name;

    @PositiveOrZero(message = "Quantity must be non-negative")
    public int quantity;

    public String location;

    /**
     * Static helper method for finding by SKU.
     * @param sku The SKU to search for.
     * @return The InventoryItem or null.
     */
    public static InventoryItem findBySku(String sku) {
        // Secure query using Panache (prevents SQL Injection)
        return find("sku = ?1", sku).firstResult();
    }
}