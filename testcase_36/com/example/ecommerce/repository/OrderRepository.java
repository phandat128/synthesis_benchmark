package com.example.ecommerce.repository;

import com.example.ecommerce.model.Order;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface OrderRepository extends JpaRepository<Order, Long> {
    // Standard CRUD operations provided by JpaRepository
    // findById is inherited, but defining a specific finder can sometimes be useful
    Optional<Order> findByOrderId(Long orderId);
}