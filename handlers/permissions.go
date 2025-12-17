package handlers

import (
	"fpscore/database"
	"fpscore/models"

	"github.com/gofiber/fiber/v2"
)

// Predefined permissions list
var PredefinedPermissions = []models.Permission{
	{Code: "assessments.create", Description: "Create new assessments"},
	{Code: "assessments.view", Description: "View assessments"},
	{Code: "assessments.edit", Description: "Edit assessments"},
	{Code: "assessments.delete", Description: "Delete assessments"},
	{Code: "users.create", Description: "Create new users"},
	{Code: "users.view", Description: "View users"},
	{Code: "users.edit", Description: "Edit users"},
	{Code: "users.delete", Description: "Delete users"},
	{Code: "roles.create", Description: "Create new roles"},
	{Code: "roles.view", Description: "View roles"},
	{Code: "roles.edit", Description: "Edit roles"},
	{Code: "roles.delete", Description: "Delete roles"},
	{Code: "reports.view", Description: "View reports"},
	{Code: "reports.export", Description: "Export reports"},
	{Code: "dashboard.view", Description: "View dashboard"},
	{Code: "facilities.view", Description: "View facilities"},
	{Code: "facilities.manage", Description: "Manage facilities"},
	{Code: "admin.areas.view", Description: "View administrative areas"},
	{Code: "admin.areas.manage", Description: "Manage administrative areas"},
	{Code: "health_workers.view", Description: "View health workers"},
	{Code: "health_workers.create", Description: "Create health workers"},
	{Code: "health_workers.edit", Description: "Edit health workers"},
	{Code: "health_workers.delete", Description: "Delete health workers"},
	{Code: "logs.view", Description: "View system logs"},
}

// GetPredefinedPermissions returns the list of all predefined permissions
func GetPredefinedPermissions(c *fiber.Ctx) error {
	return c.JSON(PredefinedPermissions)
}

// InitializePermissions ensures all predefined permissions exist in the database
func InitializePermissions() error {
	for _, perm := range PredefinedPermissions {
		_, err := database.DB.Exec(`
			INSERT INTO permissions (code, description) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING
		`, perm.Code, perm.Description)
		if err != nil {
			return err
		}
	}
	return nil
}
