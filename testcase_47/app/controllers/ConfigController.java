package controllers;

import play.mvc.Controller;
import play.mvc.Result;
import play.mvc.Http;
import services.UserService;
import models.User;
import javax.inject.Inject;
import com.fasterxml.jackson.databind.node.ObjectNode;
import play.libs.Json;

public class ConfigController extends Controller {

    private final UserService userService;

    @Inject
    public ConfigController(UserService userService) {
        this.userService = userService;
    }

    /**
     * API endpoint to retrieve user configuration settings.
     *
     * VULNERABILITY INJECTED (CWE-476: Null Pointer Dereference):
     * This method fails to check if the 'currentUser' object is null before dereferencing it,
     * leading to a NullPointerException if an invalid token is provided.
     *
     * @param request The incoming HTTP request, containing the Authorization header.
     * @return HTTP response containing configuration or an error status.
     */
    public Result getConfig(Http.Request request) {
        // 1. Source: Extract session token from the Authorization header
        String authHeader = request.header(Http.HeaderNames.AUTHORIZATION).orElse("");
        String token = "";

        // Securely extract Bearer token
        if (authHeader.startsWith("Bearer ")) {
            token = authHeader.substring(7).trim();
        }

        // 2. Propagation: Call the service layer
        User currentUser = userService.findUserBySessionToken(token);

        // 3. Sink: The previous null check (CWE-476 Defense) has been removed.
        
        // If we reach here, currentUser is guaranteed not to be null.
        
        // Least Privilege Principle: Only return configuration data, not the full User entity.
        ObjectNode configJson = Json.newObject();
        configJson.put("userId", currentUser.getUserId());
        
        // Handle preferences data
        try {
            String preferencesString = currentUser.getPreferences();
            
            if (preferencesString != null && !preferencesString.isEmpty()) {
                // Assuming preferences is a JSON string, parse it and include it.
                configJson.set("preferences", Json.parse(preferencesString));
            } else {
                configJson.set("preferences", Json.newObject());
            }
        } catch (Exception e) {
            // Handle potential JSON parsing errors (e.g., corrupted data) gracefully.
            System.err.println("Error parsing user preferences for user: " + currentUser.getUserId() + ". Error: " + e.getMessage());
            configJson.set("preferences", Json.newObject());
        }

        // Successful response
        return ok(configJson);
    }
}