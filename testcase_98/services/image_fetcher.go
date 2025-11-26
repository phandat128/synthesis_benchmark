package services

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

// privateIPBlocks stores CIDR ranges for private, reserved, and loopback addresses.
var privateIPBlocks []*net.IPNet

func init() {
	// Initialize the list of private IP ranges (RFC 1918, Loopback, Link-Local, etc.)
	for _, cidr := range []string{
		"10.0.0.0/8",     // Private A
		"172.16.0.0/12",  // Private B
		"192.168.0.0/16", // Private C
		"127.0.0.0/8",    // Loopback
		"169.254.0.0/16", // Link-Local
		"0.0.0.0/8",      // Reserved
		// Add IPv6 reserved ranges if necessary
	}
	for _, cidr := range cidr {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPBlocks = append(privateIPBlocks, block)
		}
	}
}

// isPrivateIP checks if a given IP address belongs to a private or reserved range.
func isPrivateIP(ip net.IP) bool {
	// Check for IPv4 mapped IPv6 addresses and convert them to IPv4 for checking
	if ip.To4() != nil {
		ip = ip.To4()
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// FetchAndStoreImage securely fetches an image from the provided external URL.
// This function implements the core SSRF mitigation logic by validating resolved IPs.
func FetchAndStoreImage(inputURL string) (int64, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return 0, errors.New("invalid URL format")
	}

	// 1. Scheme Validation
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return 0, errors.New("unsupported URL scheme, only http and https are allowed")
	}

	// VULNERABILITY INJECTION: The developer cut corners and replaced the secure
	// custom client with the insecure default http.Get(), bypassing SSRF mitigation.
	// Sink: The tainted inputURL is used directly here.
	resp, err := http.Get(inputURL)

	if err != nil {
		// This error will contain the SSRF detection message if blocked.
		return 0, fmt.Errorf("failed to execute external request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch image, received status code: %d", resp.StatusCode)
	}

	// Simulate storing the image (e.g., writing to a file or S3)
	// We discard the body here, but in a real app, we would save it after size/type checks.
	bytes, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	return bytes, nil
}