package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
)

// SaveEvent saves an event log (or batch of events) to the database
func SaveEvent(c *fiber.Ctx) error {
	// Get user ID from context (may be null for unauthenticated events)
	var userIDValue interface{}
	if userID := c.Locals("userID"); userID != nil {
		userIDValue = userID.(int)
	} else {
		userIDValue = nil
	}

	// Parse the body into a generic map to determine structure
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if this is a batch request (has "events" array)
	if eventsData, ok := body["events"].([]interface{}); ok && len(eventsData) > 0 {
		// Handle batch of events
		savedCount := 0
		for _, eventData := range eventsData {
			eventMap, ok := eventData.(map[string]interface{})
			if !ok {
				continue
			}

			category, _ := eventMap["category"].(string)
			eventType, _ := eventMap["event"].(string)
			page, _ := eventMap["page"].(string)
			sessionID, _ := eventMap["sessionId"].(string)
			data := eventMap["data"]

			// Convert data to JSONB
			dataJSON, err := json.Marshal(data)
			if err != nil {
				continue // Skip invalid events
			}

			// Insert event into database
			_, err = database.DB.Exec(`
				INSERT INTO events (user_id, session_id, category, event_type, page, data, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, userIDValue, sessionID, category, eventType, page, string(dataJSON), time.Now())

			if err != nil {
				continue // Continue with other events
			}
			savedCount++
		}

		return c.JSON(fiber.Map{
			"success": true,
			"count":   savedCount,
		})
	}

	// Handle single event
	category, _ := body["category"].(string)
	eventType, _ := body["event"].(string)
	page, _ := body["page"].(string)
	sessionID, _ := body["sessionId"].(string)
	data := body["data"]

	// Convert data to JSONB
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event data",
		})
	}

	// Insert event into database
	_, err = database.DB.Exec(`
		INSERT INTO events (user_id, session_id, category, event_type, page, data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userIDValue, sessionID, category, eventType, page, string(dataJSON), time.Now())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save event",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// GetEvents retrieves events with filtering and pagination
func GetEvents(c *fiber.Ctx) error {
	_ = c.Locals("userID").(int) // User ID required for permission check

	// Get query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	userIDFilter := c.Query("userId", "")
	category := c.Query("category", "")
	eventType := c.Query("eventType", "")
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 500 {
		limit = 50
	}
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT e.id, e.user_id, e.session_id, e.category, e.event_type, e.page, e.data, e.created_at,
		       u.name as user_name, u.email as user_email
		FROM events e
		LEFT JOIN users u ON e.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	// Add filters
	if userIDFilter != "" {
		query += fmt.Sprintf(" AND e.user_id = $%d", argIdx)
		args = append(args, userIDFilter)
		argIdx++
	}

	if category != "" {
		query += fmt.Sprintf(" AND e.category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	if eventType != "" {
		query += fmt.Sprintf(" AND e.event_type = $%d", argIdx)
		args = append(args, eventType)
		argIdx++
	}

	if startDate != "" {
		query += fmt.Sprintf(" AND e.created_at >= $%d", argIdx)
		args = append(args, startDate)
		argIdx++
	}

	if endDate != "" {
		query += fmt.Sprintf(" AND e.created_at <= $%d", argIdx)
		args = append(args, endDate)
		argIdx++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (e.page ILIKE $%d OR e.event_type ILIKE $%d OR e.data::text ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Order by created_at descending
	query += " ORDER BY e.created_at DESC"

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve events",
		})
	}
	defer rows.Close()

	// Initialize as empty array instead of nil
	events := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var userIDVal sql.NullInt64 // Use NullInt64 to handle NULL user_id
		var sessionID, category, eventType, page string
		var dataJSON []byte
		var createdAt time.Time
		var userName, userEmail *string

		err := rows.Scan(&id, &userIDVal, &sessionID, &category, &eventType, &page, &dataJSON, &createdAt, &userName, &userEmail)
		if err != nil {
			continue
		}

		var data map[string]interface{}
		json.Unmarshal(dataJSON, &data)

		event := map[string]interface{}{
			"id":        id,
			"userId":    nil, // Default to nil
			"sessionId": sessionID,
			"category":  category,
			"eventType": eventType,
			"page":      page,
			"data":      data,
			"createdAt": createdAt.Format(time.RFC3339),
			"userName":  nil,
			"userEmail": nil,
		}

		// Only set userId if it's not NULL
		if userIDVal.Valid {
			event["userId"] = userIDVal.Int64
		}

		if userName != nil {
			event["userName"] = *userName
		}
		if userEmail != nil {
			event["userEmail"] = *userEmail
		}

		events = append(events, event)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM events e WHERE 1=1
	`
	countArgs := []interface{}{}
	countArgIdx := 1

	if userIDFilter != "" {
		countQuery += fmt.Sprintf(" AND e.user_id = $%d", countArgIdx)
		countArgs = append(countArgs, userIDFilter)
		countArgIdx++
	}

	if category != "" {
		countQuery += fmt.Sprintf(" AND e.category = $%d", countArgIdx)
		countArgs = append(countArgs, category)
		countArgIdx++
	}

	if eventType != "" {
		countQuery += fmt.Sprintf(" AND e.event_type = $%d", countArgIdx)
		countArgs = append(countArgs, eventType)
		countArgIdx++
	}

	if startDate != "" {
		countQuery += fmt.Sprintf(" AND e.created_at >= $%d", countArgIdx)
		countArgs = append(countArgs, startDate)
		countArgIdx++
	}

	if endDate != "" {
		countQuery += fmt.Sprintf(" AND e.created_at <= $%d", countArgIdx)
		countArgs = append(countArgs, endDate)
		countArgIdx++
	}

	if search != "" {
		countQuery += fmt.Sprintf(" AND (e.page ILIKE $%d OR e.event_type ILIKE $%d OR e.data::text ILIKE $%d)", countArgIdx, countArgIdx, countArgIdx)
		countArgs = append(countArgs, "%"+search+"%")
		countArgIdx++
	}

	var total int
	database.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	return c.JSON(fiber.Map{
		"events": events,
		"total":  total,
		"page":   page,
		"limit":  limit,
		"pages":  (total + limit - 1) / limit,
	})
}

// GetEventCategories returns distinct categories
func GetEventCategories(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT DISTINCT category FROM events ORDER BY category")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve categories",
		})
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			continue
		}
		categories = append(categories, category)
	}

	return c.JSON(categories)
}

// GetEventTypes returns distinct event types
func GetEventTypes(c *fiber.Ctx) error {
	category := c.Query("category", "")

	var rows *sql.Rows
	var err error
	if category != "" {
		rows, err = database.DB.Query("SELECT DISTINCT event_type FROM events WHERE category = $1 ORDER BY event_type", category)
	} else {
		rows, err = database.DB.Query("SELECT DISTINCT event_type FROM events ORDER BY event_type")
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve event types",
		})
	}
	defer rows.Close()

	var eventTypes []string
	for rows.Next() {
		var eventType string
		if err := rows.Scan(&eventType); err != nil {
			continue
		}
		eventTypes = append(eventTypes, eventType)
	}

	return c.JSON(eventTypes)
}

// CheckAuditLogHealth checks if audit logging is working properly
func CheckAuditLogHealth(c *fiber.Ctx) error {
	health := map[string]interface{}{
		"status":       "ok",
		"database":     "connected",
		"table_exists": false,
		"can_insert":   false,
		"can_query":    false,
	}

	// Check if database is connected
	if database.DB == nil {
		health["status"] = "error"
		health["database"] = "not_connected"
		return c.JSON(health)
	}

	// Check if events table exists
	var tableExists bool
	err := database.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'events'
		)
	`).Scan(&tableExists)

	if err != nil {
		health["status"] = "error"
		health["error"] = err.Error()
		return c.JSON(health)
	}

	health["table_exists"] = tableExists

	if !tableExists {
		health["status"] = "error"
		health["message"] = "Events table does not exist. Run schema.sql to create it."
		return c.JSON(health)
	}

	// Test INSERT
	testSessionID := fmt.Sprintf("health_check_%d", time.Now().UnixNano())
	_, err = database.DB.Exec(`
		INSERT INTO events (user_id, session_id, category, event_type, page, data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, nil, testSessionID, "application", "health_check", "/api/events/health", `{"test": true}`, time.Now())

	if err != nil {
		health["status"] = "error"
		health["can_insert"] = false
		health["insert_error"] = err.Error()
	} else {
		health["can_insert"] = true
		// Clean up test record
		database.DB.Exec("DELETE FROM events WHERE session_id = $1", testSessionID)
	}

	// Test SELECT
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM events LIMIT 1").Scan(&count)
	if err != nil {
		health["status"] = "error"
		health["can_query"] = false
		health["query_error"] = err.Error()
	} else {
		health["can_query"] = true
		health["total_events"] = count
	}

	return c.JSON(health)
}
