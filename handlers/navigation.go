package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// NavigationItem represents a navigation menu item
type NavigationItem struct {
	Label    string `json:"label"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	Required string `json:"required"` // permission code required
}

// GetNavigationItems returns navigation items based on user permissions
func GetNavigationItems(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	// Check if user is admin (has admin role or all permissions)
	isAdmin := IsAdmin(userID)

	// Get all user permissions
	permissions, err := GetUserPermissions(userID)
	if err != nil {
		permissions = []string{}
	}

	// Helper to check if user has permission
	hasPermission := func(code string) bool {
		if isAdmin {
			return true
		}
		for _, p := range permissions {
			if p == code {
				return true
			}
		}
		return false
	}

	items := []NavigationItem{}

	// Home - always visible if authenticated
	items = append(items, NavigationItem{
		Label: "Home",
		URL:   "home.html",
		Icon:  "bi-house",
	})

	// Dashboard - needs dashboard.view
	if hasPermission("dashboard.view") || hasPermission("assessments.view") {
		items = append(items, NavigationItem{
			Label: "Dashboard",
			URL:   "dashboard.html",
			Icon:  "bi-speedometer2",
		})
	}

	// Reports - needs reports.view
	if hasPermission("reports.view") || hasPermission("reports.export") {
		items = append(items, NavigationItem{
			Label: "Reports",
			URL:   "reports.html",
			Icon:  "bi-file-earmark-text",
		})
	}

	// Users - needs users.view
	if hasPermission("users.view") {
		items = append(items, NavigationItem{
			Label: "Users",
			URL:   "users.html",
			Icon:  "bi-people",
		})
	}

	// Roles - needs roles.view
	if hasPermission("roles.view") {
		items = append(items, NavigationItem{
			Label: "Roles",
			URL:   "roles.html",
			Icon:  "bi-shield-check",
		})
	}

	// Facilities - needs facilities.view or admin.areas.view
	if hasPermission("facilities.view") || hasPermission("admin.areas.view") {
		items = append(items, NavigationItem{
			Label: "Facilities",
			URL:   "facilities.html",
			Icon:  "bi-hospital",
		})
	}

	// Health Workers - needs health_workers.view
	if hasPermission("health_workers.view") || hasPermission("health_workers.create") {
		items = append(items, NavigationItem{
			Label: "Health Workers",
			URL:   "health-workers.html",
			Icon:  "bi-person-badge",
		})
	}

	return c.JSON(items)
}
