package middleware

import (
	"bytes"
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

		// Capture response by creating a buffer
		var responseBody bytes.Buffer
		c.Response().SetBodyStreamWriter(func(w *bytes.Buffer) {
			responseBody = *w
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

		// Determine event category and type
		category := "api"
		eventType := fmt.Sprintf("%s %s", method, path)

		// More specific event types based on the endpoint
		if path == "/api/auth/login" {
			category = "auth"
			eventType = "login"
		} else if path == "/api/auth/logout" {
			category = "auth"
			eventType = "logout"
		} else if method == "POST" {
			eventType = "create"
		} else if method == "PUT" || method == "PATCH" {
			eventType = "update"
		} else if method == "DELETE" {
			eventType = "delete"
		} else if method == "GET" {
			eventType = "view"
		}

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
	for key, value := range data {
		if nested, ok := value.(map[string]interface{}); ok {
			maskSensitiveFields(nested)
		}
	}
}
