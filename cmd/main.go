package main

import (
	"log"
	"os"
	"path/filepath"

	"fpscore/config"
	"fpscore/database"
	"fpscore/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Get project root directory (where go.mod is located)
	// This works whether running from project root or cmd directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	// Try to find project root by looking for go.mod
	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			// Reached filesystem root, use current working directory
			projectRoot = wd
			break
		}
		projectRoot = parent
	}

	publicPath := filepath.Join(projectRoot, "public")

	// Verify public directory exists
	if _, err := os.Stat(publicPath); os.IsNotExist(err) {
		log.Fatalf("Public directory not found at: %s", publicPath)
	}

	log.Printf("Serving static files from: %s", publicPath)

	// Serve login.html as default (must be before Static)
	app.Get("/", func(c *fiber.Ctx) error {
		loginPath := filepath.Join(publicPath, "login.html")
		return c.SendFile(loginPath)
	})

	// Serve index.html for authenticated users
	app.Get("/home", func(c *fiber.Ctx) error {
		indexPath := filepath.Join(publicPath, "index.html")
		return c.SendFile(indexPath)
	})

	// Static files - serve other assets
	app.Static("/", publicPath)

	// Public routes (no auth required)
	app.Post("/api/auth/login", handlers.Login)
	app.Post("/api/auth/bootstrap", handlers.BootstrapAdmin)

	// API routes (auth required)
	api := app.Group("/api", handlers.AuthMiddleware)
	api.Post("/auth/change-password", handlers.ChangePassword)

	// Geographic hierarchy routes (view only - needs admin.areas.view or facilities.view)
	api.Get("/regions", handlers.CheckPermission("admin.areas.view"), handlers.GetRegions)
	api.Get("/districts/:regionId", handlers.CheckPermission("admin.areas.view"), handlers.GetDistricts)
	api.Get("/subcounties/:districtId", handlers.CheckPermission("admin.areas.view"), handlers.GetSubcounties)
	api.Get("/facilities/:subcountyId", handlers.CheckPermission("facilities.view"), handlers.GetFacilities)
	api.Get("/facilities", handlers.CheckPermission("facilities.view"), handlers.ListFacilities)
	api.Post("/facilities", handlers.CheckPermission("facilities.manage"), handlers.CreateFacility)
	api.Put("/facilities/:id", handlers.CheckPermission("facilities.manage"), handlers.UpdateFacility)
	api.Delete("/facilities/:id", handlers.CheckPermission("facilities.manage"), handlers.DeleteFacility)

	// Assessment routes
	api.Get("/assessment-types", handlers.CheckPermission("assessments.view"), handlers.GetAssessmentTypes)
	api.Get("/assessment-types/:typeId/thematic-areas", handlers.CheckPermission("assessments.view"), handlers.GetThematicAreas)
	api.Get("/thematic-areas/:thematicAreaId/questions", handlers.CheckPermission("assessments.view"), handlers.GetQuestions)
	api.Post("/assessments", handlers.CheckPermission("assessments.create"), handlers.CreateAssessment)
	api.Get("/assessments", handlers.CheckPermission("assessments.view"), handlers.GetAssessments)
	api.Get("/assessments/:id", handlers.CheckPermission("assessments.view"), handlers.GetAssessment)
	api.Get("/assessments/:id/summary", handlers.CheckPermission("assessments.view"), handlers.GetAssessmentSummary)

	// User management routes
	api.Get("/users", handlers.CheckPermission("users.view"), handlers.GetUsers)
	api.Get("/users/:id", handlers.CheckPermission("users.view"), handlers.GetUser)
	api.Post("/users", handlers.CheckPermission("users.create"), handlers.CreateUser)
	api.Put("/users/:id", handlers.CheckPermission("users.edit"), handlers.UpdateUser)
	api.Delete("/users/:id", handlers.CheckPermission("users.delete"), handlers.DeleteUser)
	api.Get("/users/me/permissions", handlers.GetUserPermissionsHandler) // Users can always see their own permissions
	api.Get("/navigation", handlers.GetNavigationItems)                  // Get navigation items based on permissions

	// Roles management
	api.Get("/roles", handlers.CheckPermission("roles.view"), handlers.GetRoles)
	api.Get("/roles/:id", handlers.CheckPermission("roles.view"), handlers.GetRole)
	api.Post("/roles", handlers.CheckPermission("roles.create"), handlers.CreateRole)
	api.Put("/roles/:id", handlers.CheckPermission("roles.edit"), handlers.UpdateRole)
	api.Delete("/roles/:id", handlers.CheckPermission("roles.delete"), handlers.DeleteRole)
	api.Post("/roles/:id/permissions", handlers.CheckPermission("roles.edit"), handlers.AssignPermissionsToRole)

	// Permissions management (same as roles permissions)
	api.Get("/permissions", handlers.CheckPermission("roles.view"), handlers.GetPermissions)
	api.Get("/permissions/predefined", handlers.CheckPermission("roles.view"), handlers.GetPredefinedPermissions)
	api.Post("/permissions/initialize", handlers.CheckPermission("roles.create"), handlers.InitializePermissionsHandler)
	api.Delete("/permissions/:id", handlers.CheckPermission("roles.delete"), handlers.DeletePermission)

	// Hierarchy operations
	api.Post("/hierarchy/districts/move", handlers.CheckPermission("admin.areas.manage"), handlers.MoveDistricts)
	api.Post("/hierarchy/subcounties/move", handlers.CheckPermission("admin.areas.manage"), handlers.MoveSubcounties)
	api.Post("/hierarchy/facilities/move", handlers.CheckPermission("admin.areas.manage"), handlers.MoveFacilities)

	// Reports
	api.Get("/reports/assessments/pdf", handlers.CheckPermission("reports.export"), handlers.ExportAssessmentsPDF)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
