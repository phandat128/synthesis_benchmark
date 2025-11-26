package com.example.app.controller;

import com.example.app.repository.User;
import com.example.app.service.UserService;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import java.util.Optional;

@Controller
@RequestMapping("/profile")
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    /**
     * Displays the user profile page.
     * @param principal The currently authenticated user details.
     */
    @GetMapping
    public String viewProfile(@AuthenticationPrincipal UserDetails principal, Model model) {
        String username = principal.getUsername();
        Optional<User> userOpt = userService.findByUsername(username);

        if (userOpt.isPresent()) {
            User user = userOpt.get();
            // Output Encoding: Thymeleaf automatically handles output encoding to prevent XSS.
            model.addAttribute("username", user.getUsername());
            model.addAttribute("currentEmail", user.getEmail());
            return "profile";
        } else {
            model.addAttribute("error", "User profile not found.");
            return "error"; 
        }
    }

    /**
     * Handles the state-changing request to update the user's email.
     *
     * VULNERABLE: This POST endpoint is now unprotected because CSRF was disabled globally.
     * The request will proceed even if it lacks a valid anti-CSRF token.
     *
     * @param principal The authenticated user (ensures authorization).
     * @param newEmail The input from the form.
     */
    @PostMapping("/update-email")
    public String updateEmail(
            @AuthenticationPrincipal UserDetails principal,
            @RequestParam("newEmail") String newEmail,
            RedirectAttributes redirectAttributes) {

        String username = principal.getUsername();

        // Input Validation is handled in the Service layer before persistence.
        boolean success = userService.updateUserEmail(username, newEmail);

        if (success) {
            redirectAttributes.addFlashAttribute("message", "Email updated successfully!");
        } else {
            // Proper Error Handling: Generic message for security.
            redirectAttributes.addFlashAttribute("error", "Failed to update email. Please check the format.");
        }

        return "redirect:/profile";
    }
}