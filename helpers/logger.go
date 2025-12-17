package helpers

import (
	"encoding/json"
	"fmt"
	"time"

	"fpscore/database"
)

// LogAction logs a business action with meaningful context
func LogAction(userID interface{}, action, entityType, entityName string, entityID interface{}, details string, additionalData map[string]interface{}) {
	// Prepare log data
	logData := map[string]interface{}{
		"entity_name": entityName,
	}

	// Add any additional data
	if additionalData != nil {
		for k, v := range additionalData {
			logData[k] = v
		}
	}

	// Convert to JSON
	dataJSON, _ := json.Marshal(logData)

	// Generate session ID
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())

	// Insert into database asynchronously
	go func() {
		_, err := database.DB.Exec(`
			INSERT INTO events (user_id, session_id, category, event_type, page, data, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, userID, sessionID, entityType, action, details, string(dataJSON), time.Now())

		if err != nil {
			fmt.Printf("Error logging action: %v\n", err)
		}
	}()
}

// Convenience functions for common actions

func LogCreate(userID interface{}, entityType, entityName string, entityID interface{}) {
	details := fmt.Sprintf("%s '%s' created", entityType, entityName)
	LogAction(userID, "Create", entityType, entityName, entityID, details, nil)
}

func LogUpdate(userID interface{}, entityType, entityName string, entityID interface{}, changes string) {
	details := fmt.Sprintf("%s '%s' updated", entityType, entityName)
	if changes != "" {
		details += fmt.Sprintf(": %s", changes)
	}
	LogAction(userID, "Update", entityType, entityName, entityID, details, nil)
}

func LogDelete(userID interface{}, entityType, entityName string, entityID interface{}) {
	details := fmt.Sprintf("%s '%s' deleted", entityType, entityName)
	LogAction(userID, "Delete", entityType, entityName, entityID, details, nil)
}

func LogView(userID interface{}, entityType, entityName string, entityID interface{}) {
	details := fmt.Sprintf("Viewed %s '%s'", entityType, entityName)
	LogAction(userID, "View", entityType, entityName, entityID, details, nil)
}

func LogLogin(userID interface{}, userName, ipAddress string) {
	details := fmt.Sprintf("User '%s' logged in from %s", userName, ipAddress)
	additionalData := map[string]interface{}{
		"ip_address": ipAddress,
	}
	LogAction(userID, "Login", "User", userName, userID, details, additionalData)
}

func LogLogout(userID interface{}, userName string) {
	details := fmt.Sprintf("User '%s' logged out", userName)
	LogAction(userID, "Logout", "User", userName, userID, details, nil)
}

func LogConvert(userID interface{}, fromType, fromName string, fromID interface{}, toType, toName string, toID interface{}) {
	details := fmt.Sprintf("%s '%s' converted to %s '%s'", fromType, fromName, toType, toName)
	additionalData := map[string]interface{}{
		"from_type": fromType,
		"from_id":   fromID,
		"to_type":   toType,
		"to_id":     toID,
	}
	LogAction(userID, "Convert", fromType, fromName, fromID, details, additionalData)
}

func LogSubmit(userID interface{}, entityType, entityName string, entityID interface{}, forWhom string) {
	details := fmt.Sprintf("%s '%s' submitted", entityType, entityName)
	if forWhom != "" {
		details += fmt.Sprintf(" for %s", forWhom)
	}
	LogAction(userID, "Submit", entityType, entityName, entityID, details, nil)
}

func LogEmail(userID interface{}, entityType, entityName string, recipient string) {
	details := fmt.Sprintf("%s '%s' emailed to %s", entityType, entityName, recipient)
	additionalData := map[string]interface{}{
		"recipient": recipient,
	}
	LogAction(userID, "Email Sent", entityType, entityName, nil, details, additionalData)
}

func LogExport(userID interface{}, entityType string, format string, recordCount int) {
	details := fmt.Sprintf("Exported %d %s records to %s", recordCount, entityType, format)
	additionalData := map[string]interface{}{
		"format":       format,
		"record_count": recordCount,
	}
	LogAction(userID, "Export", entityType, "", nil, details, additionalData)
}

func LogImport(userID interface{}, entityType string, recordCount int, successCount int, failCount int) {
	details := fmt.Sprintf("Imported %d %s records (%d successful, %d failed)", recordCount, entityType, successCount, failCount)
	additionalData := map[string]interface{}{
		"total":   recordCount,
		"success": successCount,
		"failed":  failCount,
	}
	LogAction(userID, "Import", entityType, "", nil, details, additionalData)
}

func LogAssign(userID interface{}, entityType, entityName string, assignedTo string) {
	details := fmt.Sprintf("%s '%s' assigned to %s", entityType, entityName, assignedTo)
	additionalData := map[string]interface{}{
		"assigned_to": assignedTo,
	}
	LogAction(userID, "Assign", entityType, entityName, nil, details, additionalData)
}

func LogMove(userID interface{}, entityType, entityName string, fromLocation, toLocation string) {
	details := fmt.Sprintf("%s '%s' moved from %s to %s", entityType, entityName, fromLocation, toLocation)
	additionalData := map[string]interface{}{
		"from": fromLocation,
		"to":   toLocation,
	}
	LogAction(userID, "Move", entityType, entityName, nil, details, additionalData)
}
