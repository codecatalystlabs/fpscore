package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
)

// APILogger middleware logs all API requests and responses to the database
func APILogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip logging for certain endpoints
		path := c.Path()
		if path == "/api/events/log" || path == "/api/events" {
			return c.Next()
		}

		// Record start time
		startTime := time.Now()

		// Get user information
		var userID interface{}
		if uid := c.Locals("userID"); uid != nil {
			userID = uid.(int)
		}

		// Capture request details
		method := c.Method()
		url := c.OriginalURL()
		userAgent := c.Get("User-Agent")
		ip := c.IP()

		// Read request body (for POST, PUT, PATCH)
		var requestBody map[string]interface{}
		if method == "POST" || method == "PUT" || method == "PATCH" {
			bodyBytes := c.Body()
			if len(bodyBytes) > 0 {
				// Try to parse as JSON
				json.Unmarshal(bodyBytes, &requestBody)

				// Mask sensitive fields
				if requestBody != nil {
					maskSensitiveFields(requestBody)
				}
			}
		}

		// Get query parameters
		queryParams := make(map[string]string)
		c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
			queryParams[string(key)] = string(value)
		})

		// Execute the handler
		err := c.Next()

		// Calculate duration
		duration := time.Since(startTime).Milliseconds()

		// Get response status
		statusCode := c.Response().StatusCode()
		success := statusCode >= 200 && statusCode < 300

		// Prepare log data
		logData := map[string]interface{}{
			"method":      method,
			"url":         url,
			"status":      statusCode,
			"duration_ms": duration,
			"ip":          ip,
			"user_agent":  userAgent,
			"success":     success,
			"entity_type": entityType,
			"description": description,
		}

		if len(queryParams) > 0 {
			logData["query_params"] = queryParams
		}

		if requestBody != nil {
			logData["request_body"] = requestBody
		}

		// Add error details if request failed
		if err != nil {
			logData["error"] = err.Error()
		}

		// Determine event category, type, and entity
		category, eventType, entityType := categorizeEndpoint(path, method)

		// Generate human-readable description
		description := generateDescription(path, method, statusCode, requestBody, userID)

		// Get session ID from cookies or generate one
		sessionID := c.Cookies("session_id", "")
		if sessionID == "" {
			sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
		}

		// Convert log data to JSON
		dataJSON, _ := json.Marshal(logData)

		// Insert into database asynchronously (don't block the response)
		go func() {
			_, err := database.DB.Exec(`
				INSERT INTO events (user_id, session_id, category, event_type, page, data, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, userID, sessionID, category, eventType, path, string(dataJSON), time.Now())

			if err != nil {
				fmt.Printf("Error logging API call: %v\n", err)
			}
		}()

		return err
	}
}

// maskSensitiveFields removes or masks sensitive data from request body
func maskSensitiveFields(data map[string]interface{}) {
	sensitiveFields := []string{"password", "token", "secret", "api_key", "apiKey", "apiSecret"}

	for _, field := range sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "***MASKED***"
		}
	}

	// Recursively mask nested objects
	for _, value := range data {
		if nested, ok := value.(map[string]interface{}); ok {
			maskSensitiveFields(nested)
		}
	}
}

// categorizeEndpoint determines category, event type, and entity type from the endpoint
func categorizeEndpoint(path, method string) (category, eventType, entityType string) {
	// Authentication endpoints
	if path == "/api/auth/login" {
		return "Auth", "Login", "User"
	}
	if path == "/api/auth/logout" {
		return "Auth", "Logout", "User"
	}
	if path == "/api/auth/change-password" {
		return "Auth", "Change Password", "User"
	}

	// Determine entity type from path
	entityType = "Unknown"
	if contains(path, "/users") {
		entityType = "User"
	} else if contains(path, "/roles") {
		entityType = "Role"
	} else if contains(path, "/permissions") {
		entityType = "Permission"
	} else if contains(path, "/facilities") {
		entityType = "Facility"
	} else if contains(path, "/health-workers") {
		entityType = "Health Worker"
	} else if contains(path, "/assessments") {
		entityType = "Assessment"
	} else if contains(path, "/regions") {
		entityType = "Region"
	} else if contains(path, "/districts") {
		entityType = "District"
	} else if contains(path, "/subcounties") {
		entityType = "Subcounty"
	} else if contains(path, "/reports") {
		entityType = "Report"
	} else if contains(path, "/navigation") {
		entityType = "Navigation"
	}

	// Determine action based on method
	category = "API"
	switch method {
	case "POST":
		eventType = "Create"
	case "PUT", "PATCH":
		eventType = "Update"
	case "DELETE":
		eventType = "Delete"
	case "GET":
		eventType = "View"
	default:
		eventType = method
	}

	return category, eventType, entityType
}

// generateDescription creates a human-readable description of what happened
func generateDescription(path, method string, statusCode int, requestBody map[string]interface{}, userID interface{}) string {
	success := statusCode >= 200 && statusCode < 300

	// Extract resource name and ID from path
	resourceName := extractResourceName(path)
	resourceID := extractResourceID(path)

	var description string

	switch method {
	case "POST":
		if success {
			if name := getNameFromBody(requestBody); name != "" {
				description = fmt.Sprintf("Created %s '%s'", resourceName, name)
			} else if resourceID != "" {
				description = fmt.Sprintf("Created %s with ID %s", resourceName, resourceID)
			} else {
				description = fmt.Sprintf("Created new %s", resourceName)
			}
		} else {
			description = fmt.Sprintf("Failed to create %s (Status: %d)", resourceName, statusCode)
		}

	case "PUT", "PATCH":
		if success {
			if name := getNameFromBody(requestBody); name != "" {
				description = fmt.Sprintf("Updated %s '%s'", resourceName, name)
			} else if resourceID != "" {
				description = fmt.Sprintf("Updated %s ID %s", resourceName, resourceID)
			} else {
				description = fmt.Sprintf("Updated %s", resourceName)
			}
		} else {
			description = fmt.Sprintf("Failed to update %s (Status: %d)", resourceName, statusCode)
		}

	case "DELETE":
		if success {
			if resourceID != "" {
				description = fmt.Sprintf("Deleted %s ID %s", resourceName, resourceID)
			} else {
				description = fmt.Sprintf("Deleted %s", resourceName)
			}
		} else {
			description = fmt.Sprintf("Failed to delete %s (Status: %d)", resourceName, statusCode)
		}

	case "GET":
		if success {
			if resourceID != "" {
				description = fmt.Sprintf("Viewed %s ID %s", resourceName, resourceID)
			} else {
				description = fmt.Sprintf("Viewed %s list", resourceName)
			}
		} else {
			description = fmt.Sprintf("Failed to view %s (Status: %d)", resourceName, statusCode)
		}

	default:
		description = fmt.Sprintf("%s %s", method, path)
	}

	// Special cases
	if path == "/api/auth/login" {
		if success {
			description = "User logged in successfully"
		} else {
			description = "Login failed - Invalid credentials"
		}
	} else if path == "/api/auth/logout" {
		description = "User logged out"
	} else if path == "/api/auth/change-password" {
		if success {
			description = "Password changed successfully"
		} else {
			description = "Failed to change password"
		}
	}

	return description
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && containsSubstring(s, substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func extractResourceName(path string) string {
	// Extract the main resource from path like "/api/users/123" -> "user"
	parts := splitPath(path)
	for i, part := range parts {
		if part == "api" && i+1 < len(parts) {
			resource := parts[i+1]
			// Singularize common plural forms
			if len(resource) > 1 && resource[len(resource)-1] == 's' {
				return resource[:len(resource)-1]
			}
			return resource
		}
	}
	return "resource"
}

func extractResourceID(path string) string {
	// Extract ID from path like "/api/users/123" -> "123"
	parts := splitPath(path)
	for i := len(parts) - 1; i >= 0; i-- {
		// Check if it looks like an ID (numeric or specific pattern)
		if isNumeric(parts[i]) {
			return parts[i]
		}
	}
	return ""
}

func splitPath(path string) []string {
	var parts []string
	current := ""
	for _, char := range path {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func getNameFromBody(body map[string]interface{}) string {
	if body == nil {
		return ""
	}

	// Try common name fields
	nameFields := []string{"name", "full_name", "fullName", "title", "email"}
	for _, field := range nameFields {
		if val, ok := body[field]; ok {
			if str, ok := val.(string); ok && str != "" {
				return str
			}
		}
	}

	return ""
}
