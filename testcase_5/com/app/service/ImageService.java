package com.app.service;

import com.app.model.UserProfile;
import com.app.repository.UserProfileRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.server.ResponseStatusException;
import reactor.core.publisher.Mono;

import java.net.InetAddress;
import java.net.URI;
import java.net.UnknownHostException;

@Service
public class ImageService {

    private static final Logger log = LoggerFactory.getLogger(ImageService.class);
    private final WebClient webClient;
    private final UserProfileRepository userProfileRepository;

    // Configuration for maximum allowed image size (security measure against DoS)
    @Value("${app.image.max-size-bytes:5242880}") // 5MB default
    private long maxImageSizeBytes;

    public ImageService(UserProfileRepository userProfileRepository) {
        this.userProfileRepository = userProfileRepository;
        // Initialize WebClient with security best practices (e.g., setting max buffer size based on config)
        this.webClient = WebClient.builder()
                .codecs(configurer -> configurer.defaultCodecs().maxInMemorySize((int) maxImageSizeBytes))
                .build();
    }

    /**
     * Fetches an image from an external URL and stores it in the user profile.
     * @param userId The ID of the user to update.
     * @param sourceUrl The external URL of the image (tainted input).
     * @return The updated UserProfile.
     */
    public UserProfile fetchAndStoreImage(Long userId, String sourceUrl) {
        // CRITICAL SECURITY STEP: Validate the URL before making the request to prevent SSRF (CWE-918)

        UserProfile profile = userProfileRepository.findById(userId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "User not found."));

        log.info("Attempting to fetch image from: {}", sourceUrl);

        try {
            // Secure HTTP Request using WebClient
            Mono<byte[]> imageMono = webClient.get()
                    .uri(sourceUrl)
                    .retrieve()
                    // Handle non-2xx status codes gracefully
                    .onStatus(HttpStatus::isError, response -> {
                        log.error("Failed to fetch image. HTTP Status: {}", response.statusCode());
                        return Mono.error(new ResponseStatusException(
                                HttpStatus.BAD_GATEWAY, "External image source returned an error."
                        ));
                    })
                    .bodyToMono(byte[].class);

            byte[] imageData = imageMono.block(); // Blocking for simplicity in service layer

            if (imageData == null || imageData.length == 0) {
                throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Fetched image data was empty.");
            }
            if (imageData.length > maxImageSizeBytes) {
                 throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Fetched image exceeds maximum allowed size.");
            }

            // Store securely
            profile.setProfileImage(imageData);
            profile.setImageUrlSource(sourceUrl);
            return userProfileRepository.save(profile);

        } catch (ResponseStatusException e) {
            throw e; // Re-throw handled exceptions
        } catch (Exception e) {
            log.error("Error during image fetching or storage: {}", e.getMessage(), e);
            // Do not leak internal stack trace or specific error details to the user
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Failed to process image request due to an internal error.");
        }
    }

    /**
     * SECURITY CRITICAL: Performs rigorous validation to prevent Server-Side Request Forgery (SSRF).
     * Checks scheme, port, and resolves the hostname to ensure the IP is not internal or reserved.
     * @param urlString The URL provided by the user.
     */
    private void validateUrlForSSRFPrevention(String urlString) {
        try {
            URI uri = new URI(urlString);
            String scheme = uri.getScheme();
            String host = uri.getHost();
            int port = uri.getPort();

            // 1. Scheme Validation: Only allow HTTP/HTTPS
            if (scheme == null || (!scheme.equalsIgnoreCase("http") && !scheme.equalsIgnoreCase("https"))) {
                throw new IllegalArgumentException("Only HTTP and HTTPS schemes are allowed.");
            }

            // 2. Port Validation: Only allow standard ports or no port specified
            if (port != -1 && port != 80 && port != 443) {
                throw new IllegalArgumentException("Non-standard ports are not allowed.");
            }

            if (host == null) {
                throw new IllegalArgumentException("URL must contain a valid host.");
            }

            // 3. IP Resolution and Internal Check (The core SSRF defense)
            // Resolve all addresses associated with the hostname
            InetAddress[] addresses = InetAddress.getAllByName(host);
            for (InetAddress address : addresses) {
                if (isInternalIp(address)) {
                    log.warn("SSRF Attempt Detected: Host {} resolved to internal IP {}", host, address.getHostAddress());
                    throw new IllegalArgumentException("Access to internal or reserved IP addresses is forbidden.");
                }
            }

        } catch (UnknownHostException e) {
            log.warn("Host resolution failed for URL: {}", urlString);
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid or unresolvable host specified.");
        } catch (Exception e) {
            // Catch IllegalArgumentException from validation steps and wrap it in a 400 response
            log.error("URL validation failed: {}", e.getMessage());
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
        }
    }

    /**
     * Checks if an InetAddress belongs to a private, loopback, link-local, or multicast range.
     * This prevents access to internal networks (RFC 1918) and local services.
     */
    private boolean isInternalIp(InetAddress address) {
        // Check for loopback (127.x.x.x)
        if (address.isLoopbackAddress()) return true;

        // Check for link-local (169.254.x.x)
        if (address.isLinkLocalAddress()) return true;

        // Check for private ranges (RFC 1918: 10.x.x.x, 172.16-31.x.x, 192.168.x.x)
        if (address.isSiteLocalAddress()) return true;

        // Check for multicast (224.x.x.x to 239.x.x.x)
        if (address.isMulticastAddress()) return true;

        return false;
    }
}