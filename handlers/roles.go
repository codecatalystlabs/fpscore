package handlers

import (
	"database/sql"
	"strconv"

	"fpscore/database"
	"fpscore/models"

	"github.com/gofiber/fiber/v2"
)

// CreateRole creates a new role
func CreateRole(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	var roleID int
	err := database.DB.QueryRow(`
		INSERT INTO roles (name, description) 
		VALUES ($1, $2) 
		RETURNING id
	`, req.Name, req.Description).Scan(&roleID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"id": roleID, "message": "Role created successfully"})
}

// UpdateRole updates an existing role
func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	_, err = database.DB.Exec(`
		UPDATE roles SET name=$1, description=$2 
		WHERE id=$3
	`, req.Name, req.Description, roleID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Role updated successfully"})
}

// DeleteRole deletes a role
func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	_, err = database.DB.Exec(`DELETE FROM roles WHERE id=$1`, roleID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Role deleted successfully"})
}

// GetRole returns a role with its permissions
func GetRole(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	var role models.Role
	err = database.DB.QueryRow(`
		SELECT id, name, description 
		FROM roles 
		WHERE id = $1
	`, roleID).Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "Role not found")
		}
		return err
	}

	// Get permissions for this role
	permRows, err := database.DB.Query(`
		SELECT p.id, p.code, p.description 
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.code
	`, roleID)
	if err == nil {
		defer permRows.Close()
		var permissions []models.Permission
		for permRows.Next() {
			var p models.Permission
			if err := permRows.Scan(&p.ID, &p.Code, &p.Description); err == nil {
				permissions = append(permissions, p)
			}
		}
		// Add permissions to response
		return c.JSON(fiber.Map{
			"role":        role,
			"permissions": permissions,
		})
	}

	return c.JSON(role)
}

// AssignPermissionsToRole assigns permissions to a role
func AssignPermissionsToRole(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role ID")
	}

	var req struct {
		PermissionIDs []int `json:"permissionIds"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove all existing permissions
	_, err = tx.Exec(`DELETE FROM role_permissions WHERE role_id=$1`, roleID)
	if err != nil {
		return err
	}

	// Add new permissions
	for _, permID := range req.PermissionIDs {
		tx.Exec(`INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, roleID, permID)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Permissions assigned successfully"})
}

// InitializePermissionsHandler creates all predefined permissions in the database
func InitializePermissionsHandler(c *fiber.Ctx) error {
	// Import the InitializePermissions function from permissions.go
	// Since it's in the same package, we can call it directly
	if err := InitializePermissions(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to initialize permissions: "+err.Error())
	}
	return c.JSON(fiber.Map{"message": "Permissions initialized successfully"})
}

// DeletePermission deletes a permission
func DeletePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	permID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid permission ID")
	}

	_, err = database.DB.Exec(`DELETE FROM permissions WHERE id=$1`, permID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Permission deleted successfully"})
}

// GetUserPermissionsHandler returns all permissions for the current user (API endpoint)
func GetUserPermissionsHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	rows, err := database.DB.Query(`
		SELECT DISTINCT p.id, p.code, p.description 
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY p.code
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Description); err == nil {
			permissions = append(permissions, p)
		}
	}

	return c.JSON(permissions)
}
