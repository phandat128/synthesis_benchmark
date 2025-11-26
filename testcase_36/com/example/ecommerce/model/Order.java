package com.example.ecommerce.model;

import com.example.ecommerce.enums.OrderState;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

import javax.persistence.*;
import java.math.BigDecimal;

@Entity
@Table(name = "orders")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Order {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long orderId;

    @Column(nullable = false)
    private BigDecimal totalAmount;

    // Securely store the state using an Enum
    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private OrderState currentStatus;

    @Column(nullable = false)
    private String customerEmail;

    // Initialization constructor
    public Order(BigDecimal totalAmount, String customerEmail) {
        this.totalAmount = totalAmount;
        this.customerEmail = customerEmail;
        this.currentStatus = OrderState.CART; // Default initial state
    }
}