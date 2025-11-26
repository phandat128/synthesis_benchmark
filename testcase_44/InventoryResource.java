package com.inventory.resource;

import com.inventory.model.InventoryItem;
import jakarta.annotation.security.Authenticated;
import jakarta.annotation.security.RolesAllowed;
import jakarta.transaction.Transactional;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import jakarta.validation.Valid;
import java.util.List;

@Path("/api/inventory")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
@Authenticated // Requires authentication for all methods in this resource
public class InventoryResource {

    /**
     * Retrieves all inventory items.
     * Requires the 'USER' role.
     */
    @GET
    @RolesAllowed("USER")
    public List<InventoryItem> getAll() {
        return InventoryItem.listAll();
    }

    /**
     * Retrieves a specific inventory item by ID.
     * Requires the 'USER' role.
     */
    @GET
    @Path("/{id}")
    @RolesAllowed("USER")
    public Response getById(@PathParam("id") Long id) {
        InventoryItem item = InventoryItem.findById(id);
        if (item == null) {
            // Secure error handling: Use standard 404
            return Response.status(Response.Status.NOT_FOUND).entity("Item not found.").build();
        }
        return Response.ok(item).build();
    }

    /**
     * Creates a new inventory item.
     * Requires the 'ADMIN' role. Input validation via @Valid.
     */
    @POST
    @Transactional
    @RolesAllowed("ADMIN")
    public Response create(@Valid InventoryItem item) {
        // Input validation check for business logic (unique SKU)
        if (InventoryItem.findBySku(item.sku) != null) {
            return Response.status(Response.Status.CONFLICT).entity("SKU already exists.").build();
        }

        item.persist();
        return Response.status(Response.Status.CREATED).entity(item).build();
    }

    /**
     * Updates an existing inventory item.
     * Requires the 'ADMIN' role.
     */
    @PUT
    @Path("/{id}")
    @Transactional
    @RolesAllowed("ADMIN")
    public Response update(@PathParam("id") Long id, @Valid InventoryItem updatedItem) {
        InventoryItem entity = InventoryItem.findById(id);

        if (entity == null) {
            return Response.status(Response.Status.NOT_FOUND).entity("Item not found for update.").build();
        }

        // Apply updates securely
        entity.sku = updatedItem.sku;
        entity.name = updatedItem.name;
        entity.quantity = updatedItem.quantity;
        entity.location = updatedItem.location;

        return Response.ok(entity).build();
    }

    /**
     * Deletes an inventory item.
     * Requires the 'ADMIN' role.
     */
    @DELETE
    @Path("/{id}")
    @Transactional
    @RolesAllowed("ADMIN")
    public Response delete(@PathParam("id") Long id) {
        boolean deleted = InventoryItem.deleteById(id);
        if (deleted) {
            return Response.noContent().build();
        } else {
            return Response.status(Response.Status.NOT_FOUND).entity("Item not found for deletion.").build();
        }
    }
}