package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// NavigationItem represents a navigation menu item
type NavigationItem struct {
	Label string `json:"label"`
	URL   string `json:"url"`
	Icon  string `json:"icon"`
}

// GetNavigationItems returns nav items for the selected tool only.
// Query: ?tool=proficiency | rhspars (required for tool pages; empty returns []).
func GetNavigationItems(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)
	permissions, err := GetUserPermissions(userID)
	if err != nil {
		permissions = []string{}
	}

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

	tool := strings.ToLower(strings.TrimSpace(c.Query("tool")))
	items := []NavigationItem{}

	// Always offer a way back to the tile picker when inside a tool
	if tool == "proficiency" || tool == "rhspars" {
		items = append(items, NavigationItem{
			Label: "All tools",
			URL:   "tools.html",
			Icon:  "bi-grid-3x3-gap",
		})
	}

	switch tool {
	case "proficiency":
		if !hasPermission("tools.proficiency.access") && !hasPermission("assessments.view") && !hasPermission("assessments.create") {
			return c.JSON(items)
		}
		items = append(items, NavigationItem{
			Label: "Assessments",
			URL:   "home.html",
			Icon:  "bi-clipboard-check",
		})
		if hasPermission("dashboard.view") || hasPermission("assessments.view") {
			items = append(items, NavigationItem{
				Label: "Dashboard",
				URL:   "dashboard.html",
				Icon:  "bi-speedometer2",
			})
		}
		if hasPermission("reports.view") || hasPermission("reports.export") {
			items = append(items, NavigationItem{
				Label: "Reports",
				URL:   "reports.html",
				Icon:  "bi-file-earmark-text",
			})
		}
		if hasPermission("health_workers.view") || hasPermission("health_workers.create") {
			items = append(items, NavigationItem{
				Label: "Health Workers",
				URL:   "health-workers.html",
				Icon:  "bi-person-badge",
			})
		}
		appendAdminNav(&items, hasPermission)

	case "rhspars":
		if !hasPermission("tools.rh_spars.access") && !hasPermission("rhspars.view") && !hasPermission("rhspars.create") {
			return c.JSON(items)
		}
		items = append(items, NavigationItem{
			Label: "New visit",
			URL:   "rhspars-home.html",
			Icon:  "bi-hospital",
		})
		if hasPermission("rhspars.view") {
			items = append(items, NavigationItem{
				Label: "History",
				URL:   "rhspars-history.html",
				Icon:  "bi-clock-history",
			})
		}
		appendAdminNav(&items, hasPermission)

	default:
		// tools.html and unknown contexts: no sidebar items
		return c.JSON([]NavigationItem{})
	}

	return c.JSON(items)
}

func appendAdminNav(items *[]NavigationItem, hasPermission func(string) bool) {
	if hasPermission("users.view") {
		*items = append(*items, NavigationItem{
			Label: "Users",
			URL:   "users.html",
			Icon:  "bi-people",
		})
	}
	if hasPermission("roles.view") {
		*items = append(*items, NavigationItem{
			Label: "Roles & permissions",
			URL:   "roles.html",
			Icon:  "bi-shield-check",
		})
	}
	if hasPermission("facilities.view") || hasPermission("admin.areas.view") {
		*items = append(*items, NavigationItem{
			Label: "Geography & facilities",
			URL:   "facilities.html",
			Icon:  "bi-geo-alt",
		})
	}
	if hasPermission("logs.view") {
		*items = append(*items, NavigationItem{
			Label: "Audit Trail",
			URL:   "logs.html",
			Icon:  "bi-journal-text",
		})
	}
}
