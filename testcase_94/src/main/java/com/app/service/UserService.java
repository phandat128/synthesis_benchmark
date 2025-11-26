package com.app.service;

import com.app.dto.UserProfileUpdateDTO;
import com.app.model.User;
import com.app.repository.UserRepository;
import org.springframework.beans.BeanUtils; // VULNERABILITY: Added BeanUtils for generic copying
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import java.util.Optional;

@Service
public class UserService {

    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public Optional<User> findById(Long userId) {
        return userRepository.findById(userId);
    }

    /**
     * VULNERABLE IMPLEMENTATION: Updates a user profile using generic property copying (Mass Assignment).
     * This method now allows unauthorized fields (like 'role') present in the DTO to overwrite the entity.
     *
     * @param userId The ID of the user to update.
     * @param updateDTO The DTO containing allowed update fields.
     * @return The updated User entity.
     * @throws IllegalArgumentException if the user is not found or input is invalid.
     */
    @Transactional
    public User updateUser(Long userId, UserProfileUpdateDTO updateDTO) {
        User userEntity = userRepository.findById(userId)
                .orElseThrow(() -> new IllegalArgumentException("User not found with ID: " + userId));

        // --- VULNERABILITY INJECTION (Propagation/Sink): Mass Assignment (CWE-915) ---
        // Replacing secure explicit field mapping with generic property copying.
        // This copies all matching fields from updateDTO to userEntity, including the sensitive 'role'.
        BeanUtils.copyProperties(updateDTO, userEntity, "id", "passwordHash");

        // 3. Save the entity, now containing the attacker-supplied role
        return userRepository.save(userEntity);
    }
}