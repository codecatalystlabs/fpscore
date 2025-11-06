package handlers

import (
	"fmt"
	"log"
	"strings"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
)

// IsAdmin checks if a user is an admin
func IsAdmin(userID int) bool {
	// First check if user has a role with "admin" in the name
	var hasAdminRole bool
	err := database.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM user_roles ur
			JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1 AND r.name ILIKE '%admin%'
		)
	`, userID).Scan(&hasAdminRole)
	if err == nil && hasAdminRole {
		return true
	}

	// Then check if user has all permissions (only if permissions exist)
	var totalPermissions int
	err = database.DB.QueryRow(`SELECT COUNT(*) FROM permissions`).Scan(&totalPermissions)
	if err != nil || totalPermissions == 0 {
		// If no permissions exist, can't determine admin status this way
		return false
	}

	var userPermissionCount int
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT p.code)
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = $1
	`, userID).Scan(&userPermissionCount)
	if err != nil {
		return false
	}

	// User is admin if they have all permissions
	return userPermissionCount == totalPermissions && userPermissionCount > 0
}

// CheckPermission checks if the current user has a specific permission
func CheckPermission(permissionCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(int)

		// Admin users bypass permission checks
		if IsAdmin(userID) {
			return c.Next()
		}

		// Check if user has the specific permission through their roles
		var hasPermission bool
		err := database.DB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM user_roles ur
				JOIN role_permissions rp ON ur.role_id = rp.role_id
				JOIN permissions p ON rp.permission_id = p.id
				WHERE ur.user_id = $1 AND p.code = $2
			)
		`, userID, permissionCode).Scan(&hasPermission)

		// Deny access on error or if permission not found
		if err != nil {
			// Log error for debugging
			log.Printf("Permission check error for user %d, permission %s: %v", userID, permissionCode, err)
			return fiber.NewError(fiber.StatusForbidden, "You do not have permission to perform this action")
		}

		if !hasPermission {
			// Log denied access for debugging
			log.Printf("Access denied for user %d, permission %s: permission not found", userID, permissionCode)
			return fiber.NewError(fiber.StatusForbidden, "You do not have permission to perform this action")
		}

		// Log successful access for debugging (commented out for production)
		// log.Printf("Access granted for user %d, permission %s", userID, permissionCode)

		return c.Next()
	}
}

// CheckAnyPermission checks if the user has any of the specified permissions
func CheckAnyPermission(permissionCodes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(int)

		// Admin users bypass permission checks
		if IsAdmin(userID) {
			return c.Next()
		}

		// Build placeholders for IN clause
		placeholders := make([]string, len(permissionCodes))
		args := make([]interface{}, len(permissionCodes)+1)
		args[0] = userID
		for i, code := range permissionCodes {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args[i+1] = code
		}

		var hasPermission bool
		err := database.DB.QueryRow(fmt.Sprintf(`
			SELECT EXISTS(
				SELECT 1 FROM user_roles ur
				JOIN role_permissions rp ON ur.role_id = rp.role_id
				JOIN permissions p ON rp.permission_id = p.id
				WHERE ur.user_id = $1 AND p.code IN (%s)
			)
		`, strings.Join(placeholders, ",")), args...).Scan(&hasPermission)

		if err != nil || !hasPermission {
			return fiber.NewError(fiber.StatusForbidden, "You do not have permission to perform this action")
		}

		return c.Next()
	}
}

// GetUserPermissions returns all permission codes for the current user
func GetUserPermissions(userID int) ([]string, error) {
	rows, err := database.DB.Query(`
		SELECT DISTINCT p.code 
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil {
			permissions = append(permissions, code)
		}
	}

	return permissions, nil
}

// RequirePermission is a helper that checks if user has permission, returns error if not
func RequirePermission(userID int, permissionCode string) error {
	var hasPermission bool
	err := database.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM user_roles ur
			JOIN role_permissions rp ON ur.role_id = rp.role_id
			JOIN permissions p ON rp.permission_id = p.id
			WHERE ur.user_id = $1 AND p.code = $2
		)
	`, userID, permissionCode).Scan(&hasPermission)

	if err != nil {
		return err
	}

	if !hasPermission {
		return fiber.NewError(fiber.StatusForbidden, "You do not have permission to perform this action")
	}

	return nil
}
