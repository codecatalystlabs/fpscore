package handlers

import (
	"database/sql"
	"strconv"

	"fpscore/database"
	"fpscore/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers returns all users (with pagination)
func GetUsers(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
		SELECT id, name, email, is_active, created_at 
		FROM users 
		ORDER BY created_at DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.IsActive, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}

	return c.JSON(users)
}

// GetUser returns a single user with roles and admin areas
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	var user models.User
	err = database.DB.QueryRow(`
		SELECT id, name, email, is_active, created_at 
		FROM users 
		WHERE id = $1
	`, userID).Scan(&user.ID, &user.Name, &user.Email, &user.IsActive, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return err
	}

	// Get roles
	roleRows, err := database.DB.Query(`
		SELECT r.id, r.name, r.description 
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`, userID)
	var roles []models.Role
	if err == nil {
		defer roleRows.Close()
		for roleRows.Next() {
			var r models.Role
			if err := roleRows.Scan(&r.ID, &r.Name, &r.Description); err == nil {
				roles = append(roles, r)
			}
		}
	}

	// Get admin areas
	areaRows, err := database.DB.Query(`
		SELECT id, region_id, district_id, subcounty_id, facility_id 
		FROM user_admin_areas 
		WHERE user_id = $1
	`, userID)
	var areas []models.UserAdminArea
	if err == nil {
		defer areaRows.Close()
		for areaRows.Next() {
			var a models.UserAdminArea
			var regionID, districtID, subcountyID, facilityID sql.NullInt64
			if err := areaRows.Scan(&a.ID, &regionID, &districtID, &subcountyID, &facilityID); err == nil {
				if regionID.Valid {
					rid := int(regionID.Int64)
					a.RegionID = &rid
				}
				if districtID.Valid {
					did := int(districtID.Int64)
					a.DistrictID = &did
				}
				if subcountyID.Valid {
					sid := int(subcountyID.Int64)
					a.SubcountyID = &sid
				}
				if facilityID.Valid {
					fid := int(facilityID.Int64)
					a.FacilityID = &fid
				}
				a.UserID = userID
				areas = append(areas, a)
			}
		}
	}

	return c.JSON(fiber.Map{
		"user":       user,
		"roles":      roles,
		"adminAreas": areas,
	})
}

type CreateUserRequest struct {
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	Password   string           `json:"password"`
	IsActive   bool             `json:"isActive"`
	RoleIDs    []int            `json:"roleIds"`
	AdminAreas []AdminAreaInput `json:"adminAreas"`
}

type AdminAreaInput struct {
	RegionID    *int `json:"regionId,omitempty"`
	DistrictID  *int `json:"districtId,omitempty"`
	SubcountyID *int `json:"subcountyId,omitempty"`
	FacilityID  *int `json:"facilityId,omitempty"`
}

// CreateUser creates a new user
func CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow(`
		INSERT INTO users (name, email, password_hash, is_active) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`, req.Name, req.Email, string(hash), req.IsActive).Scan(&userID)
	if err != nil {
		return err
	}

	// Add roles
	for _, roleID := range req.RoleIDs {
		tx.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, roleID)
	}

	// Add admin areas
	for _, area := range req.AdminAreas {
		tx.Exec(`
			INSERT INTO user_admin_areas (user_id, region_id, district_id, subcounty_id, facility_id) 
			VALUES ($1, $2, $3, $4, $5)
		`, userID, area.RegionID, area.DistrictID, area.SubcountyID, area.FacilityID)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"id": userID, "message": "User created successfully"})
}

// UpdateUser updates an existing user
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update user
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
			UPDATE users SET name=$1, email=$2, password_hash=$3, is_active=$4 
			WHERE id=$5
		`, req.Name, req.Email, string(hash), req.IsActive, userID)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec(`
			UPDATE users SET name=$1, email=$2, is_active=$3 
			WHERE id=$4
		`, req.Name, req.Email, req.IsActive, userID)
		if err != nil {
			return err
		}
	}

	// Update roles
	tx.Exec(`DELETE FROM user_roles WHERE user_id=$1`, userID)
	for _, roleID := range req.RoleIDs {
		tx.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
	}

	// Update admin areas
	tx.Exec(`DELETE FROM user_admin_areas WHERE user_id=$1`, userID)
	for _, area := range req.AdminAreas {
		tx.Exec(`
			INSERT INTO user_admin_areas (user_id, region_id, district_id, subcounty_id, facility_id) 
			VALUES ($1, $2, $3, $4, $5)
		`, userID, area.RegionID, area.DistrictID, area.SubcountyID, area.FacilityID)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

// DeleteUser deletes a user
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	_, err = database.DB.Exec(`DELETE FROM users WHERE id=$1`, userID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

// GetRoles returns all roles
func GetRoles(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`SELECT id, name, description FROM roles ORDER BY name`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.ID, &r.Name, &r.Description); err == nil {
			roles = append(roles, r)
		}
	}

	return c.JSON(roles)
}

// GetPermissions returns all permissions
func GetPermissions(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`SELECT id, code, description FROM permissions ORDER BY code`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var perms []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Description); err == nil {
			perms = append(perms, p)
		}
	}

	return c.JSON(perms)
}
