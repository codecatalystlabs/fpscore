package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"fpscore/database"
	"fpscore/models"

	"github.com/gofiber/fiber/v2"
)

// ListFacilities returns facilities with optional filters and access restrictions
func ListFacilities(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	regionIDs, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
	if err != nil {
		return err
	}

	query := `
		SELECT
			f.id,
			f.name,
			sc.id AS subcounty_id,
			sc.name AS subcounty_name,
			d.id AS district_id,
			d.name AS district_name,
			r.id AS region_id,
			r.name AS region_name,
			f.created_at
		FROM facilities f
		JOIN subcounties sc ON f.subcounty_id = sc.id
		JOIN districts d ON sc.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	appendFilter := func(format string, value interface{}) {
		query += fmt.Sprintf(" AND %s", fmt.Sprintf(format, argIdx))
		args = append(args, value)
		argIdx++
	}

	if v := c.Query("regionId"); v != "" {
		if id, convErr := strconv.Atoi(v); convErr == nil {
			appendFilter("r.id = $%d", id)
		}
	}
	if v := c.Query("districtId"); v != "" {
		if id, convErr := strconv.Atoi(v); convErr == nil {
			appendFilter("d.id = $%d", id)
		}
	}
	if v := c.Query("subcountyId"); v != "" {
		if id, convErr := strconv.Atoi(v); convErr == nil {
			appendFilter("sc.id = $%d", id)
		}
	}
	if search := strings.TrimSpace(c.Query("search")); search != "" {
		query += fmt.Sprintf(" AND LOWER(f.name) LIKE LOWER($%d)", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if !isAdmin {
		accessConditions := []string{}
		if len(regionIDs) > 0 {
			placeholders := []string{}
			for _, id := range regionIDs {
				placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
				args = append(args, id)
				argIdx++
			}
			accessConditions = append(accessConditions, fmt.Sprintf("r.id IN (%s)", strings.Join(placeholders, ",")))
		}
		if len(districtIDs) > 0 {
			placeholders := []string{}
			for _, id := range districtIDs {
				placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
				args = append(args, id)
				argIdx++
			}
			accessConditions = append(accessConditions, fmt.Sprintf("d.id IN (%s)", strings.Join(placeholders, ",")))
		}
		if len(subcountyIDs) > 0 {
			placeholders := []string{}
			for _, id := range subcountyIDs {
				placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
				args = append(args, id)
				argIdx++
			}
			accessConditions = append(accessConditions, fmt.Sprintf("sc.id IN (%s)", strings.Join(placeholders, ",")))
		}
		if len(facilityIDs) > 0 {
			placeholders := []string{}
			for _, id := range facilityIDs {
				placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
				args = append(args, id)
				argIdx++
			}
			accessConditions = append(accessConditions, fmt.Sprintf("f.id IN (%s)", strings.Join(placeholders, ",")))
		}
		if len(accessConditions) > 0 {
			query += " AND (" + strings.Join(accessConditions, " OR ") + ")"
		} else {
			// No access configured, return empty list
			return c.JSON([]models.Facility{})
		}
	}

	query += " ORDER BY r.name, d.name, sc.name, f.name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	type facilityResponse struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		RegionID      int    `json:"regionId"`
		RegionName    string `json:"regionName"`
		DistrictID    int    `json:"districtId"`
		DistrictName  string `json:"districtName"`
		SubcountyID   int    `json:"subcountyId"`
		SubcountyName string `json:"subcountyName"`
		CreatedAt     string `json:"createdAt"`
	}

	var facilities []facilityResponse
	for rows.Next() {
		var item facilityResponse
		var createdAt sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &item.SubcountyID, &item.SubcountyName, &item.DistrictID, &item.DistrictName, &item.RegionID, &item.RegionName, &createdAt); err != nil {
			return err
		}
		item.CreatedAt = createdAt.String
		facilities = append(facilities, item)
	}

	return c.JSON(facilities)
}

type facilityRequest struct {
	Name        string `json:"name"`
	SubcountyID int    `json:"subcountyId"`
}

func CreateFacility(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	var req facilityRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if strings.TrimSpace(req.Name) == "" || req.SubcountyID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Name and subcounty are required")
	}

	regionID, districtID, err := getRegionDistrictForSubcounty(req.SubcountyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid subcounty")
		}
		return err
	}

	if !isAdmin {
		if !userHasAccessToArea(userID, regionID, districtID, req.SubcountyID, 0) {
			return fiber.NewError(fiber.StatusForbidden, "You do not have access to this administrative area")
		}
	}

	var facilityID int
	err = database.DB.QueryRow(`
		INSERT INTO facilities (name, subcounty_id)
		VALUES ($1, $2)
		RETURNING id
	`, req.Name, req.SubcountyID).Scan(&facilityID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Facility created successfully",
		"facilityId": facilityID,
	})
}

func UpdateFacility(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	idParam := c.Params("id")
	facilityID, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid facility ID")
	}

	var req facilityRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if strings.TrimSpace(req.Name) == "" || req.SubcountyID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Name and subcounty are required")
	}

	regionID, districtID, err := getRegionDistrictForSubcounty(req.SubcountyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid subcounty")
		}
		return err
	}

	if !isAdmin {
		if !userHasAccessToArea(userID, regionID, districtID, req.SubcountyID, facilityID) {
			return fiber.NewError(fiber.StatusForbidden, "You do not have access to this administrative area")
		}
	}

	_, err = database.DB.Exec(`
		UPDATE facilities
		SET name = $1, subcounty_id = $2
		WHERE id = $3
	`, req.Name, req.SubcountyID, facilityID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Facility updated successfully"})
}

func DeleteFacility(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	idParam := c.Params("id")
	facilityID, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid facility ID")
	}

	if !isAdmin {
		regionID, districtID, subcountyID, err := getAreaForFacility(facilityID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusNotFound, "Facility not found")
			}
			return err
		}
		if !userHasAccessToArea(userID, regionID, districtID, subcountyID, facilityID) {
			return fiber.NewError(fiber.StatusForbidden, "You do not have access to this administrative area")
		}
	}

	result, err := database.DB.Exec(`DELETE FROM facilities WHERE id = $1`, facilityID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Facility not found")
	}

	return c.JSON(fiber.Map{"message": "Facility deleted successfully"})
}

func getRegionDistrictForSubcounty(subcountyID int) (int, int, error) {
	var regionID, districtID int
	err := database.DB.QueryRow(`
		SELECT d.region_id, d.id
		FROM subcounties sc
		JOIN districts d ON sc.district_id = d.id
		WHERE sc.id = $1
	`, subcountyID).Scan(&regionID, &districtID)
	return regionID, districtID, err
}

func getAreaForFacility(facilityID int) (regionID int, districtID int, subcountyID int, err error) {
	err = database.DB.QueryRow(`
		SELECT r.id, d.id, sc.id
		FROM facilities f
		JOIN subcounties sc ON f.subcounty_id = sc.id
		JOIN districts d ON sc.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE f.id = $1
	`, facilityID).Scan(&regionID, &districtID, &subcountyID)
	return
}

func userHasAccessToArea(userID int, regionID int, districtID int, subcountyID int, facilityID int) bool {
	if IsAdmin(userID) {
		return true
	}

	accessRegions, accessDistricts, accessSubcounties, accessFacilities, err := GetUserAdminAreas(userID)
	if err != nil {
		return false
	}

	// If user has explicit facility access
	for _, id := range accessFacilities {
		if id == facilityID && facilityID != 0 {
			return true
		}
	}

	// Subcounty access
	for _, id := range accessSubcounties {
		if id == subcountyID && subcountyID != 0 {
			return true
		}
	}

	// District access
	for _, id := range accessDistricts {
		if id == districtID && districtID != 0 {
			return true
		}
	}

	// Region access
	for _, id := range accessRegions {
		if id == regionID && regionID != 0 {
			return true
		}
	}

	// If user has no specific admin area assigned, deny access
	if len(accessRegions) == 0 && len(accessDistricts) == 0 && len(accessSubcounties) == 0 && len(accessFacilities) == 0 {
		return false
	}

	return false
}

func getRegionForDistrict(districtID int) (int, error) {
	var regionID int
	err := database.DB.QueryRow(`SELECT region_id FROM districts WHERE id = $1`, districtID).Scan(&regionID)
	return regionID, err
}
