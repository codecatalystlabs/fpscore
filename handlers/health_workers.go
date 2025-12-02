package handlers

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
)

// Health worker handlers

// GetHealthWorkers returns list of health workers with optional filters
func GetHealthWorkers(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	query := `
		SELECT hw.id, hw.full_name, hw.email, hw.phone_number, hw.facility_id,
		       f.name as facility_name, hw.created_at, hw.updated_at
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	// Apply admin area restrictions if not admin
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(regionIDs) > 0 || len(districtIDs) > 0 || len(subcountyIDs) > 0 || len(facilityIDs) > 0 {
			query += " AND ("
			conditions := []string{}
			if len(regionIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(regionIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, regionIDs[i])
					argIdx++
				}
				conditions = append(conditions, "r.id IN ("+placeholders+")")
			}
			if len(districtIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(districtIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, districtIDs[i])
					argIdx++
				}
				conditions = append(conditions, "d.id IN ("+placeholders+")")
			}
			if len(subcountyIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(subcountyIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, subcountyIDs[i])
					argIdx++
				}
				conditions = append(conditions, "s.id IN ("+placeholders+")")
			}
			if len(facilityIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(facilityIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, facilityIDs[i])
					argIdx++
				}
				conditions = append(conditions, "f.id IN ("+placeholders+")")
			}
			query += strings.Join(conditions, " OR ") + ")"
		} else {
			query += " AND 1=0"
		}
	}

	// Query parameter filters
	if regionID := c.Query("regionId"); regionID != "" {
		query += fmt.Sprintf(" AND r.id = $%d", argIdx)
		args = append(args, regionID)
		argIdx++
	}
	if districtID := c.Query("districtId"); districtID != "" {
		query += fmt.Sprintf(" AND d.id = $%d", argIdx)
		args = append(args, districtID)
		argIdx++
	}
	if subcountyID := c.Query("subcountyId"); subcountyID != "" {
		query += fmt.Sprintf(" AND s.id = $%d", argIdx)
		args = append(args, subcountyID)
		argIdx++
	}
	if facilityID := c.Query("facilityId"); facilityID != "" {
		query += fmt.Sprintf(" AND hw.facility_id = $%d", argIdx)
		args = append(args, facilityID)
		argIdx++
	}
	if search := c.Query("search"); search != "" {
		query += fmt.Sprintf(" AND (hw.full_name ILIKE $%d OR hw.email ILIKE $%d OR hw.phone_number ILIKE $%d)", argIdx, argIdx, argIdx)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern)
		argIdx++
	}

	query += " ORDER BY hw.full_name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var healthWorkers []map[string]interface{}
	for rows.Next() {
		var hwID, facilityID int
		var fullName, facilityName, createdAt, updatedAt string
		var email, phoneNumber sql.NullString
		if err := rows.Scan(&hwID, &fullName, &email, &phoneNumber, &facilityID, &facilityName, &createdAt, &updatedAt); err != nil {
			return err
		}

		hw := map[string]interface{}{
			"id":           hwID,
			"fullName":     fullName,
			"facilityId":   facilityID,
			"facilityName": facilityName,
			"createdAt":    createdAt,
			"updatedAt":    updatedAt,
		}
		if email.Valid {
			hw["email"] = email.String
		}
		if phoneNumber.Valid {
			hw["phoneNumber"] = phoneNumber.String
		}
		healthWorkers = append(healthWorkers, hw)
	}

	return c.JSON(healthWorkers)
}

// GetHealthWorker returns a single health worker by ID
func GetHealthWorker(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	query := `
		SELECT hw.id, hw.full_name, hw.email, hw.phone_number, hw.facility_id,
		       f.name as facility_name, hw.created_at, hw.updated_at
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE hw.id = $1
	`
	args := []interface{}{id}
	argIdx := 2

	// Apply admin area restrictions if not admin
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(regionIDs) > 0 || len(districtIDs) > 0 || len(subcountyIDs) > 0 || len(facilityIDs) > 0 {
			query += " AND ("
			conditions := []string{}
			if len(regionIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(regionIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, regionIDs[i])
					argIdx++
				}
				conditions = append(conditions, "r.id IN ("+placeholders+")")
			}
			if len(districtIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(districtIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, districtIDs[i])
					argIdx++
				}
				conditions = append(conditions, "d.id IN ("+placeholders+")")
			}
			if len(subcountyIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(subcountyIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, subcountyIDs[i])
					argIdx++
				}
				conditions = append(conditions, "s.id IN ("+placeholders+")")
			}
			if len(facilityIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(facilityIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, facilityIDs[i])
					argIdx++
				}
				conditions = append(conditions, "f.id IN ("+placeholders+")")
			}
			query += strings.Join(conditions, " OR ") + ")"
		} else {
			return fiber.NewError(404, "Health worker not found")
		}
	}

	var hwID, facilityID int
	var fullName, facilityName, createdAt, updatedAt string
	var email, phoneNumber sql.NullString
	err := database.DB.QueryRow(query, args...).Scan(&hwID, &fullName, &email, &phoneNumber, &facilityID, &facilityName, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Health worker not found")
		}
		return err
	}

	hw := map[string]interface{}{
		"id":           hwID,
		"fullName":     fullName,
		"facilityId":   facilityID,
		"facilityName": facilityName,
		"createdAt":    createdAt,
		"updatedAt":    updatedAt,
	}
	if email.Valid {
		hw["email"] = email.String
	}
	if phoneNumber.Valid {
		hw["phoneNumber"] = phoneNumber.String
	}

	return c.JSON(hw)
}

// CreateHealthWorkerRequest represents the request body for creating a health worker
type CreateHealthWorkerRequest struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	FacilityID  int    `json:"facilityId"`
}

// CreateHealthWorker creates a new health worker
func CreateHealthWorker(c *fiber.Ctx) error {
	var req CreateHealthWorkerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "Invalid request body")
	}

	// Validate: at least one of email or phone must be provided
	if req.Email == "" && req.PhoneNumber == "" {
		return fiber.NewError(400, "Either email or phone number must be provided")
	}

	// Validate: full name is required
	if req.FullName == "" {
		return fiber.NewError(400, "Full name is required")
	}

	// Validate: facility exists
	var facilityExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM facilities WHERE id = $1)", req.FacilityID).Scan(&facilityExists)
	if err != nil {
		return err
	}
	if !facilityExists {
		return fiber.NewError(400, "Facility not found")
	}

	// Check for duplicate email or phone
	var existingID int
	checkQuery := "SELECT id FROM health_workers WHERE "
	checkArgs := []interface{}{}
	checkIdx := 1
	conditions := []string{}
	if req.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", checkIdx))
		checkArgs = append(checkArgs, req.Email)
		checkIdx++
	}
	if req.PhoneNumber != "" {
		if len(conditions) > 0 {
			conditions = append(conditions, "OR")
		}
		conditions = append(conditions, fmt.Sprintf("phone_number = $%d", checkIdx))
		checkArgs = append(checkArgs, req.PhoneNumber)
		checkIdx++
	}
	checkQuery += strings.Join(conditions, " ")

	err = database.DB.QueryRow(checkQuery, checkArgs...).Scan(&existingID)
	if err == nil {
		return fiber.NewError(400, "A health worker with this email or phone number already exists")
	} else if err != sql.ErrNoRows {
		return err
	}

	// Insert health worker
	var id int
	var email, phoneNumber sql.NullString
	if req.Email != "" {
		email = sql.NullString{String: req.Email, Valid: true}
	}
	if req.PhoneNumber != "" {
		phoneNumber = sql.NullString{String: req.PhoneNumber, Valid: true}
	}

	err = database.DB.QueryRow(`
		INSERT INTO health_workers (full_name, email, phone_number, facility_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, req.FullName, email, phoneNumber, req.FacilityID).Scan(&id)
	if err != nil {
		return err
	}

	// Get created health worker with facility name
	var fullName, facilityName, createdAt, updatedAt string
	err = database.DB.QueryRow(`
		SELECT hw.full_name, f.name, hw.created_at, hw.updated_at
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		WHERE hw.id = $1
	`, id).Scan(&fullName, &facilityName, &createdAt, &updatedAt)
	if err != nil {
		return err
	}

	hw := map[string]interface{}{
		"id":           id,
		"fullName":     fullName,
		"facilityId":   req.FacilityID,
		"facilityName": facilityName,
		"createdAt":    createdAt,
		"updatedAt":    updatedAt,
	}
	if email.Valid {
		hw["email"] = email.String
	}
	if phoneNumber.Valid {
		hw["phoneNumber"] = phoneNumber.String
	}

	return c.Status(201).JSON(hw)
}

// UpdateHealthWorkerRequest represents the request body for updating a health worker
type UpdateHealthWorkerRequest struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	FacilityID  int    `json:"facilityId"`
}

// UpdateHealthWorker updates an existing health worker
func UpdateHealthWorker(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateHealthWorkerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "Invalid request body")
	}

	// Validate: at least one of email or phone must be provided
	if req.Email == "" && req.PhoneNumber == "" {
		return fiber.NewError(400, "Either email or phone number must be provided")
	}

	// Validate: full name is required
	if req.FullName == "" {
		return fiber.NewError(400, "Full name is required")
	}

	// Validate: facility exists
	var facilityExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM facilities WHERE id = $1)", req.FacilityID).Scan(&facilityExists)
	if err != nil {
		return err
	}
	if !facilityExists {
		return fiber.NewError(400, "Facility not found")
	}

	// Check for duplicate email or phone (excluding current health worker)
	var existingID int
	checkQuery := "SELECT id FROM health_workers WHERE id != $1 AND ("
	checkArgs := []interface{}{id}
	checkIdx := 2
	conditions := []string{}
	if req.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", checkIdx))
		checkArgs = append(checkArgs, req.Email)
		checkIdx++
	}
	if req.PhoneNumber != "" {
		if len(conditions) > 0 {
			conditions = append(conditions, "OR")
		}
		conditions = append(conditions, fmt.Sprintf("phone_number = $%d", checkIdx))
		checkArgs = append(checkArgs, req.PhoneNumber)
		checkIdx++
	}
	checkQuery += strings.Join(conditions, " ") + ")"

	err = database.DB.QueryRow(checkQuery, checkArgs...).Scan(&existingID)
	if err == nil {
		return fiber.NewError(400, "A health worker with this email or phone number already exists")
	} else if err != sql.ErrNoRows {
		return err
	}

	// Update health worker
	var email, phoneNumber sql.NullString
	if req.Email != "" {
		email = sql.NullString{String: req.Email, Valid: true}
	}
	if req.PhoneNumber != "" {
		phoneNumber = sql.NullString{String: req.PhoneNumber, Valid: true}
	}

	_, err = database.DB.Exec(`
		UPDATE health_workers
		SET full_name = $1, email = $2, phone_number = $3, facility_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`, req.FullName, email, phoneNumber, req.FacilityID, id)
	if err != nil {
		return err
	}

	// Get updated health worker
	var fullName, facilityName, createdAt, updatedAt string
	err = database.DB.QueryRow(`
		SELECT hw.full_name, f.name, hw.created_at, hw.updated_at
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		WHERE hw.id = $1
	`, id).Scan(&fullName, &facilityName, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Health worker not found")
		}
		return err
	}

	hw := map[string]interface{}{
		"id":           id,
		"fullName":     fullName,
		"facilityId":   req.FacilityID,
		"facilityName": facilityName,
		"createdAt":    createdAt,
		"updatedAt":    updatedAt,
	}
	if email.Valid {
		hw["email"] = email.String
	}
	if phoneNumber.Valid {
		hw["phoneNumber"] = phoneNumber.String
	}

	return c.JSON(hw)
}

// MoveHealthWorkerRequest represents the request body for moving a health worker
type MoveHealthWorkerRequest struct {
	FacilityID int `json:"facilityId"`
}

// MoveHealthWorker moves a health worker to a different facility
func MoveHealthWorker(c *fiber.Ctx) error {
	id := c.Params("id")
	var req MoveHealthWorkerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "Invalid request body")
	}

	// Validate: facility exists
	var facilityExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM facilities WHERE id = $1)", req.FacilityID).Scan(&facilityExists)
	if err != nil {
		return err
	}
	if !facilityExists {
		return fiber.NewError(400, "Facility not found")
	}

	// Update health worker facility
	_, err = database.DB.Exec(`
		UPDATE health_workers
		SET facility_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, req.FacilityID, id)
	if err != nil {
		return err
	}

	// Get updated health worker
	var fullName, facilityName, createdAt, updatedAt string
	var email, phoneNumber sql.NullString
	err = database.DB.QueryRow(`
		SELECT hw.full_name, hw.email, hw.phone_number, f.name, hw.created_at, hw.updated_at
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		WHERE hw.id = $1
	`, id).Scan(&fullName, &email, &phoneNumber, &facilityName, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Health worker not found")
		}
		return err
	}

	hw := map[string]interface{}{
		"id":           id,
		"fullName":     fullName,
		"facilityId":   req.FacilityID,
		"facilityName": facilityName,
		"createdAt":    createdAt,
		"updatedAt":    updatedAt,
	}
	if email.Valid {
		hw["email"] = email.String
	}
	if phoneNumber.Valid {
		hw["phoneNumber"] = phoneNumber.String
	}

	return c.JSON(hw)
}

// DeleteHealthWorker deletes a health worker
func DeleteHealthWorker(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if health worker has assessments
	var assessmentCount int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM assessments WHERE health_worker_id = $1", id).Scan(&assessmentCount)
	if err != nil {
		return err
	}
	if assessmentCount > 0 {
		return fiber.NewError(400, "Cannot delete health worker with existing assessments")
	}

	// Delete health worker
	result, err := database.DB.Exec("DELETE FROM health_workers WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fiber.NewError(404, "Health worker not found")
	}

	return c.SendStatus(204)
}

// GetHealthWorkerAssessments returns assessments grouped by health worker with thematic area scores
func GetHealthWorkerAssessments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	healthWorkerID := c.Query("healthWorkerId")
	if healthWorkerID == "" {
		return fiber.NewError(400, "healthWorkerId query parameter is required")
	}

	// Verify health worker exists and user has access
	query := `
		SELECT hw.id, hw.full_name, hw.email, hw.phone_number, hw.facility_id, f.name as facility_name
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE hw.id = $1
	`
	args := []interface{}{healthWorkerID}
	argIdx := 2

	// Apply admin area restrictions if not admin
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(regionIDs) > 0 || len(districtIDs) > 0 || len(subcountyIDs) > 0 || len(facilityIDs) > 0 {
			query += " AND ("
			conditions := []string{}
			if len(regionIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(regionIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, regionIDs[i])
					argIdx++
				}
				conditions = append(conditions, "r.id IN ("+placeholders+")")
			}
			if len(districtIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(districtIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, districtIDs[i])
					argIdx++
				}
				conditions = append(conditions, "d.id IN ("+placeholders+")")
			}
			if len(subcountyIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(subcountyIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, subcountyIDs[i])
					argIdx++
				}
				conditions = append(conditions, "s.id IN ("+placeholders+")")
			}
			if len(facilityIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(facilityIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, facilityIDs[i])
					argIdx++
				}
				conditions = append(conditions, "f.id IN ("+placeholders+")")
			}
			query += strings.Join(conditions, " OR ") + ")"
		} else {
			return fiber.NewError(404, "Health worker not found")
		}
	}

	var hwID, facilityID int
	var fullName, facilityName string
	var email, phoneNumber sql.NullString
	err := database.DB.QueryRow(query, args...).Scan(&hwID, &fullName, &email, &phoneNumber, &facilityID, &facilityName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Health worker not found")
		}
		return err
	}

	// Get all assessments for this health worker
	assessmentRows, err := database.DB.Query(`
		SELECT a.id, a.assessment_type_id, at.name as assessment_type, a.percentage_score,
		       a.performance_level, a.created_at, f.name as facility_name
		FROM assessments a
		JOIN assessment_types at ON a.assessment_type_id = at.id
		JOIN facilities f ON a.facility_id = f.id
		WHERE a.health_worker_id = $1
		ORDER BY a.created_at DESC
	`, healthWorkerID)
	if err != nil {
		return err
	}
	defer assessmentRows.Close()

	type AssessmentData struct {
		ID               int
		AssessmentTypeID int
		AssessmentType   string
		Percentage       float64
		PerformanceLevel string
		CreatedAt        string
		FacilityName     string
		ThematicScores   []map[string]interface{}
	}

	var assessments []AssessmentData
	for assessmentRows.Next() {
		var a AssessmentData
		if err := assessmentRows.Scan(&a.ID, &a.AssessmentTypeID, &a.AssessmentType, &a.Percentage, &a.PerformanceLevel, &a.CreatedAt, &a.FacilityName); err != nil {
			return err
		}

		// Get thematic area scores for this assessment
		thematicRows, err := database.DB.Query(`
			SELECT ta.id, ta.name, tas.possible_score, tas.achieved_score, tas.percentage_score
			FROM thematic_area_scores tas
			JOIN thematic_areas ta ON tas.thematic_area_id = ta.id
			WHERE tas.assessment_id = $1
			ORDER BY ta.display_order
		`, a.ID)
		if err != nil {
			return err
		}

		var thematicScores []map[string]interface{}
		for thematicRows.Next() {
			var taID, possible, achieved int
			var taName string
			var percentage float64
			if err := thematicRows.Scan(&taID, &taName, &possible, &achieved, &percentage); err != nil {
				thematicRows.Close()
				return err
			}
			thematicScores = append(thematicScores, map[string]interface{}{
				"id":         taID,
				"name":       taName,
				"possible":   possible,
				"achieved":   achieved,
				"percentage": percentage,
			})
		}
		thematicRows.Close()
		a.ThematicScores = thematicScores
		assessments = append(assessments, a)
	}

	// Calculate average score across all thematic areas (only for assessed areas)
	var totalThematicPercentage float64
	var thematicCount int
	for _, assessment := range assessments {
		for _, thematic := range assessment.ThematicScores {
			if percentage, ok := thematic["percentage"].(float64); ok {
				totalThematicPercentage += percentage
				thematicCount++
			}
		}
	}
	avgThematicScore := 0.0
	if thematicCount > 0 {
		avgThematicScore = totalThematicPercentage / float64(thematicCount)
	}

	hwData := map[string]interface{}{
		"id":               hwID,
		"fullName":         fullName,
		"facilityId":       facilityID,
		"facilityName":     facilityName,
		"assessments":      assessments,
		"avgThematicScore": avgThematicScore,
	}
	if email.Valid {
		hwData["email"] = email.String
	}
	if phoneNumber.Valid {
		hwData["phoneNumber"] = phoneNumber.String
	}

	return c.JSON(hwData)
}

// GetHealthWorkersWithAssessments returns health workers who have assessments, with assessment counts
func GetHealthWorkersWithAssessments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	// Get date range filters
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	query := `
		SELECT DISTINCT hw.id, hw.full_name, hw.email, hw.phone_number,
		       hw.facility_id, f.name as facility_name,
		       r.name as region_name, d.name as district_name, s.name as subcounty_name,
		       COUNT(a.id) as assessment_count,
		       MIN(a.created_at) as first_assessment,
		       MAX(a.created_at) as last_assessment
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		JOIN assessments a ON a.health_worker_id = hw.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	// Apply date range filters
	if startDate != "" {
		query += fmt.Sprintf(" AND a.created_at >= $%d", argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		query += fmt.Sprintf(" AND a.created_at <= $%d", argIdx)
		args = append(args, endDate+" 23:59:59")
		argIdx++
	}

	// Apply admin area restrictions if not admin
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(regionIDs) > 0 || len(districtIDs) > 0 || len(subcountyIDs) > 0 || len(facilityIDs) > 0 {
			query += " AND ("
			conditions := []string{}
			if len(regionIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(regionIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, regionIDs[i])
					argIdx++
				}
				conditions = append(conditions, "r.id IN ("+placeholders+")")
			}
			if len(districtIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(districtIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, districtIDs[i])
					argIdx++
				}
				conditions = append(conditions, "d.id IN ("+placeholders+")")
			}
			if len(subcountyIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(subcountyIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, subcountyIDs[i])
					argIdx++
				}
				conditions = append(conditions, "s.id IN ("+placeholders+")")
			}
			if len(facilityIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(facilityIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, facilityIDs[i])
					argIdx++
				}
				conditions = append(conditions, "f.id IN ("+placeholders+")")
			}
			query += strings.Join(conditions, " OR ") + ")"
		} else {
			query += " AND 1=0"
		}
	}

	query += " GROUP BY hw.id, hw.full_name, hw.email, hw.phone_number, hw.facility_id, f.name, r.name, d.name, s.name"
	query += " ORDER BY hw.full_name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var healthWorkers []map[string]interface{}
	for rows.Next() {
		var hwID, facilityID, assessmentCount int
		var fullName, facilityName, regionName, districtName, subcountyName, firstAssessment, lastAssessment string
		var email, phoneNumber sql.NullString

		if err := rows.Scan(&hwID, &fullName, &email, &phoneNumber, &facilityID, &facilityName,
			&regionName, &districtName, &subcountyName, &assessmentCount, &firstAssessment, &lastAssessment); err != nil {
			return err
		}

		hw := map[string]interface{}{
			"id":              hwID,
			"fullName":        fullName,
			"facilityId":      facilityID,
			"facilityName":    facilityName,
			"regionName":      regionName,
			"districtName":    districtName,
			"subcountyName":   subcountyName,
			"assessmentCount": assessmentCount,
			"firstAssessment": firstAssessment,
			"lastAssessment":  lastAssessment,
		}
		if email.Valid {
			hw["email"] = email.String
		}
		if phoneNumber.Valid {
			hw["phoneNumber"] = phoneNumber.String
		}
		healthWorkers = append(healthWorkers, hw)
	}

	return c.JSON(healthWorkers)
}

// GetHealthWorkerPerformance returns performance by thematic area for a health worker
func GetHealthWorkerPerformance(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	healthWorkerID := c.Params("id")
	if healthWorkerID == "" {
		return fiber.NewError(400, "Health worker ID is required")
	}

	// Get date range filters
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Verify health worker exists and user has access
	query := `
		SELECT hw.id, hw.full_name, hw.email, hw.phone_number, hw.facility_id, f.name as facility_name,
		       r.name as region_name, d.name as district_name, s.name as subcounty_name
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE hw.id = $1
	`
	args := []interface{}{healthWorkerID}
	argIdx := 2

	// Apply admin area restrictions if not admin
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(regionIDs) > 0 || len(districtIDs) > 0 || len(subcountyIDs) > 0 || len(facilityIDs) > 0 {
			query += " AND ("
			conditions := []string{}
			if len(regionIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(regionIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, regionIDs[i])
					argIdx++
				}
				conditions = append(conditions, "r.id IN ("+placeholders+")")
			}
			if len(districtIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(districtIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, districtIDs[i])
					argIdx++
				}
				conditions = append(conditions, "d.id IN ("+placeholders+")")
			}
			if len(subcountyIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(subcountyIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, subcountyIDs[i])
					argIdx++
				}
				conditions = append(conditions, "s.id IN ("+placeholders+")")
			}
			if len(facilityIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(facilityIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, facilityIDs[i])
					argIdx++
				}
				conditions = append(conditions, "f.id IN ("+placeholders+")")
			}
			query += strings.Join(conditions, " OR ") + ")"
		} else {
			return fiber.NewError(404, "Health worker not found")
		}
	}

	var hwID, facilityID int
	var fullName, facilityName, regionName, districtName, subcountyName string
	var email, phoneNumber sql.NullString
	err := database.DB.QueryRow(query, args...).Scan(&hwID, &fullName, &email, &phoneNumber, &facilityID, &facilityName,
		&regionName, &districtName, &subcountyName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Health worker not found")
		}
		return err
	}

	// Get thematic area scores for this health worker (only assessed thematic areas)
	// Only show actual thematic areas (Pre-Procedure, Procedure, Post-Procedure patterns)
	// Exclude detailed sections like "Maintains privacy and confidentiality"
	// Group by name and assessment type to ensure each thematic area appears only once
	taQuery := `
		SELECT MIN(ta.id) as id, ta.name, ta.assessment_type_id, at.name as assessment_type_name, MIN(ta.display_order) as display_order,
		       SUM(tas.possible_score) as total_possible,
		       SUM(tas.achieved_score) as total_achieved,
		       COUNT(DISTINCT tas.assessment_id) as assessment_count
		FROM thematic_areas ta
		JOIN assessment_types at ON ta.assessment_type_id = at.id
		JOIN thematic_area_scores tas ON ta.id = tas.thematic_area_id
		JOIN assessments a ON tas.assessment_id = a.id
		WHERE a.health_worker_id = $1
		AND at.name IN (
			'Counselling',
			'Implant Insertion',
			'Implant Removal',
			'Injectable Progesterone Only',
			'IUD/IUS Insertion',
			'IUD/IUS Removal',
			'Mini-Laparotomy Tubal Ligation',
			'Progesterone only pill and combined oral contraceptive pill',
			'Vasectomy'
		)
		AND (
			ta.name LIKE 'Pre-Procedure%' OR
			ta.name LIKE 'Procedure:%' OR
			ta.name LIKE 'Post-Procedure%' OR
			ta.name LIKE 'Pre-removal%' OR
			ta.name LIKE 'Removal of%' OR
			ta.name LIKE 'Initial Assessment%' OR
			ta.name LIKE 'Client Eligibility%' OR
			ta.name LIKE 'Additional Eligibility%' OR
			ta.name LIKE 'Instructional and Documentation%' OR
			ta.name LIKE 'Getting Ready%' OR
			ta.name LIKE 'Pre-operative%' OR
			ta.name LIKE 'For post partum%' OR
			ta.name LIKE 'Pre insertion%' OR
			ta.name LIKE 'Injection Procedure%' OR
			ta.name LIKE 'Client Information and Eligibility%' OR
			ta.name LIKE 'Administering Injection%'
		)
	`
	taArgs := []interface{}{healthWorkerID}
	taArgIdx := 2

	if startDate != "" {
		taQuery += fmt.Sprintf(" AND a.created_at >= $%d", taArgIdx)
		taArgs = append(taArgs, startDate)
		taArgIdx++
	}
	if endDate != "" {
		taQuery += fmt.Sprintf(" AND a.created_at <= $%d", taArgIdx)
		taArgs = append(taArgs, endDate+" 23:59:59")
		taArgIdx++
	}

	// Group by name and assessment type to ensure uniqueness
	taQuery += " GROUP BY ta.name, ta.assessment_type_id, at.name ORDER BY at.name, MIN(ta.display_order)"

	taRows, err := database.DB.Query(taQuery, taArgs...)
	if err != nil {
		return err
	}
	defer taRows.Close()

	type ThematicAreaScore struct {
		ID               int     `json:"id"`
		Name             string  `json:"name"`
		AssessmentTypeID int     `json:"assessmentTypeId"`
		AssessmentType   string  `json:"assessmentType"`
		TotalPossible    int     `json:"totalPossible"`
		TotalAchieved    int     `json:"totalAchieved"`
		Percentage       float64 `json:"percentage"`
		AssessmentCount  int     `json:"assessmentCount"`
	}

	var thematicScores []ThematicAreaScore
	var totalPercentage float64
	var assessedCount int

	// Use composite key (name + assessment_type_id) to track unique thematic areas
	seenThematicAreas := make(map[string]bool) // Track seen thematic areas by name+type to avoid duplicates

	for taRows.Next() {
		var taID, assessmentTypeID, totalPossible, totalAchieved, assessmentCount int
		var taName, assessmentTypeName string
		var displayOrder int

		if err := taRows.Scan(&taID, &taName, &assessmentTypeID, &assessmentTypeName, &displayOrder,
			&totalPossible, &totalAchieved, &assessmentCount); err != nil {
			continue
		}

		// Create composite key: thematic area name + assessment type
		// This ensures we only show each thematic area once per assessment type
		compositeKey := fmt.Sprintf("%s|%s", taName, assessmentTypeName)

		// Skip if we've already seen this thematic area for this assessment type
		if seenThematicAreas[compositeKey] {
			continue
		}
		seenThematicAreas[compositeKey] = true

		percentage := 0.0
		if totalPossible > 0 {
			percentage = (float64(totalAchieved) / float64(totalPossible)) * 100
		}

		thematicScores = append(thematicScores, ThematicAreaScore{
			ID:               taID,
			Name:             taName,
			AssessmentTypeID: assessmentTypeID,
			AssessmentType:   assessmentTypeName,
			TotalPossible:    totalPossible,
			TotalAchieved:    totalAchieved,
			Percentage:       percentage,
			AssessmentCount:  assessmentCount,
		})

		totalPercentage += percentage
		assessedCount++
	}

	// Calculate average score (only for assessed thematic areas)
	avgScore := 0.0
	if assessedCount > 0 {
		avgScore = totalPercentage / float64(assessedCount)
	}

	result := map[string]interface{}{
		"healthWorker": map[string]interface{}{
			"id":            hwID,
			"fullName":      fullName,
			"facilityId":    facilityID,
			"facilityName":  facilityName,
			"regionName":    regionName,
			"districtName":  districtName,
			"subcountyName": subcountyName,
		},
		"thematicScores": thematicScores,
		"averageScore":   avgScore,
		"assessedCount":  assessedCount,
	}

	if email.Valid {
		result["healthWorker"].(map[string]interface{})["email"] = email.String
	}
	if phoneNumber.Valid {
		result["healthWorker"].(map[string]interface{})["phoneNumber"] = phoneNumber.String
	}

	return c.JSON(result)
}

// GetHealthWorkerThematicAreaDetails returns detailed questions and answers for a health worker's thematic area
func GetHealthWorkerThematicAreaDetails(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	healthWorkerID := c.Params("id")
	thematicAreaID := c.Params("thematicAreaId")
	if healthWorkerID == "" || thematicAreaID == "" {
		return fiber.NewError(400, "Health worker ID and thematic area ID are required")
	}

	// Get date range filters
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Verify health worker exists and user has access
	query := `
		SELECT hw.id, hw.full_name
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE hw.id = $1
	`
	args := []interface{}{healthWorkerID}
	argIdx := 2

	// Apply admin area restrictions if not admin
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(regionIDs) > 0 || len(districtIDs) > 0 || len(subcountyIDs) > 0 || len(facilityIDs) > 0 {
			query += " AND ("
			conditions := []string{}
			if len(regionIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(regionIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, regionIDs[i])
					argIdx++
				}
				conditions = append(conditions, "r.id IN ("+placeholders+")")
			}
			if len(districtIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(districtIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, districtIDs[i])
					argIdx++
				}
				conditions = append(conditions, "d.id IN ("+placeholders+")")
			}
			if len(subcountyIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(subcountyIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, subcountyIDs[i])
					argIdx++
				}
				conditions = append(conditions, "s.id IN ("+placeholders+")")
			}
			if len(facilityIDs) > 0 {
				placeholders := ""
				for i := 0; i < len(facilityIDs); i++ {
					if i > 0 {
						placeholders += ","
					}
					placeholders += fmt.Sprintf("$%d", argIdx)
					args = append(args, facilityIDs[i])
					argIdx++
				}
				conditions = append(conditions, "f.id IN ("+placeholders+")")
			}
			query += strings.Join(conditions, " OR ") + ")"
		} else {
			return fiber.NewError(404, "Health worker not found")
		}
	}

	var hwID int
	var fullName string
	err := database.DB.QueryRow(query, args...).Scan(&hwID, &fullName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Health worker not found")
		}
		return err
	}

	// Get thematic area info
	var taName, assessmentTypeName string
	var taID int
	err = database.DB.QueryRow(`
		SELECT ta.id, ta.name, at.name
		FROM thematic_areas ta
		JOIN assessment_types at ON ta.assessment_type_id = at.id
		WHERE ta.id = $1
	`, thematicAreaID).Scan(&taID, &taName, &assessmentTypeName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Thematic area not found")
		}
		return err
	}

	// Find all thematic area IDs with the same name and assessment type
	// This ensures we show all assessments even if there are multiple IDs with the same name
	var allThematicAreaIDs []int
	allIDsRows, err := database.DB.Query(`
		SELECT ta.id
		FROM thematic_areas ta
		JOIN assessment_types at ON ta.assessment_type_id = at.id
		WHERE ta.name = $1 AND at.name = $2
	`, taName, assessmentTypeName)
	if err == nil {
		defer allIDsRows.Close()
		for allIDsRows.Next() {
			var id int
			if err := allIDsRows.Scan(&id); err == nil {
				allThematicAreaIDs = append(allThematicAreaIDs, id)
			}
		}
	}
	// If no additional IDs found, use the original ID
	if len(allThematicAreaIDs) == 0 {
		allThematicAreaIDs = []int{taID}
	}

	// Get only questions that have been answered for this health worker in this thematic area
	// Include all thematic area IDs with the same name
	placeholders := ""
	for i := range allThematicAreaIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", i+1)
	}
	questionsQuery := fmt.Sprintf(`
		SELECT DISTINCT q.id, q.question_text, q.score_weight, q.is_critical, q.is_important, q.display_order
		FROM questions q
		JOIN assessment_responses ar ON q.id = ar.question_id
		JOIN assessments a ON ar.assessment_id = a.id
		WHERE q.thematic_area_id IN (%s)
		AND a.health_worker_id = $%d
	`, placeholders, len(allThematicAreaIDs)+1)
	questionArgs := []interface{}{}
	for _, id := range allThematicAreaIDs {
		questionArgs = append(questionArgs, id)
	}
	questionArgs = append(questionArgs, healthWorkerID)
	questionArgIdx := len(questionArgs) + 1

	if startDate != "" {
		questionsQuery += fmt.Sprintf(" AND a.created_at >= $%d", questionArgIdx)
		questionArgs = append(questionArgs, startDate)
		questionArgIdx++
	}
	if endDate != "" {
		questionsQuery += fmt.Sprintf(" AND a.created_at <= $%d", questionArgIdx)
		questionArgs = append(questionArgs, endDate+" 23:59:59")
		questionArgIdx++
	}

	questionsQuery += " ORDER BY q.display_order"

	questionRows, err := database.DB.Query(questionsQuery, questionArgs...)
	if err != nil {
		return err
	}
	defer questionRows.Close()

	type QuestionDetail struct {
		ID           int    `json:"id"`
		Text         string `json:"text"`
		ScoreWeight  int    `json:"scoreWeight"`
		IsCritical   bool   `json:"isCritical"`
		IsImportant  bool   `json:"isImportant"`
		DisplayOrder int    `json:"displayOrder"`
		Responses    []struct {
			AssessmentID int    `json:"assessmentId"`
			Response     string `json:"response"`
			PointsEarned int    `json:"pointsEarned"`
			Date         string `json:"date"`
		} `json:"responses"`
	}

	var questions []QuestionDetail
	questionMap := make(map[int]*QuestionDetail)

	for questionRows.Next() {
		var q QuestionDetail
		if err := questionRows.Scan(&q.ID, &q.Text, &q.ScoreWeight, &q.IsCritical, &q.IsImportant, &q.DisplayOrder); err != nil {
			continue
		}
		q.Responses = []struct {
			AssessmentID int    `json:"assessmentId"`
			Response     string `json:"response"`
			PointsEarned int    `json:"pointsEarned"`
			Date         string `json:"date"`
		}{}
		questionMap[q.ID] = &q
		questions = append(questions, q)
	}

	// Get all responses for these questions from this health worker's assessments
	// Include all thematic area IDs with the same name
	responsePlaceholders := ""
	for i := range allThematicAreaIDs {
		if i > 0 {
			responsePlaceholders += ","
		}
		responsePlaceholders += fmt.Sprintf("$%d", i+2)
	}
	responsesQuery := fmt.Sprintf(`
		SELECT ar.question_id, ar.response, ar.points_earned, a.id as assessment_id, a.created_at
		FROM assessment_responses ar
		JOIN assessments a ON ar.assessment_id = a.id
		WHERE a.health_worker_id = $1
		AND ar.question_id IN (
			SELECT id FROM questions WHERE thematic_area_id IN (%s)
		)
	`, responsePlaceholders)
	responseArgs := []interface{}{healthWorkerID}
	for _, id := range allThematicAreaIDs {
		responseArgs = append(responseArgs, id)
	}
	responseArgIdx := len(responseArgs) + 1

	if startDate != "" {
		responsesQuery += fmt.Sprintf(" AND a.created_at >= $%d", responseArgIdx)
		responseArgs = append(responseArgs, startDate)
		responseArgIdx++
	}
	if endDate != "" {
		responsesQuery += fmt.Sprintf(" AND a.created_at <= $%d", responseArgIdx)
		responseArgs = append(responseArgs, endDate+" 23:59:59")
		responseArgIdx++
	}

	responsesQuery += " ORDER BY a.created_at DESC, ar.question_id"

	responseRows, err := database.DB.Query(responsesQuery, responseArgs...)
	if err != nil {
		return err
	}
	defer responseRows.Close()

	// Count responses while loading them
	yesCount := 0
	noCount := 0
	naCount := 0

	for responseRows.Next() {
		var questionID, assessmentID, pointsEarned int
		var response, createdAt string
		if err := responseRows.Scan(&questionID, &response, &pointsEarned, &assessmentID, &createdAt); err != nil {
			continue
		}

		// Count responses
		if response == "Yes" {
			yesCount++
		} else if response == "No" {
			noCount++
		} else if response == "NA" {
			naCount++
		}

		// Add response to question
		if q, exists := questionMap[questionID]; exists {
			q.Responses = append(q.Responses, struct {
				AssessmentID int    `json:"assessmentId"`
				Response     string `json:"response"`
				PointsEarned int    `json:"pointsEarned"`
				Date         string `json:"date"`
			}{
				AssessmentID: assessmentID,
				Response:     response,
				PointsEarned: pointsEarned,
				Date:         createdAt,
			})
		}
	}

	// Get all assessments for this thematic area (all IDs with the same name)
	assessmentPlaceholders := ""
	for i := range allThematicAreaIDs {
		if i > 0 {
			assessmentPlaceholders += ","
		}
		assessmentPlaceholders += fmt.Sprintf("$%d", i+2)
	}
	assessmentsQuery := fmt.Sprintf(`
		SELECT DISTINCT a.id, a.created_at, tas.possible_score, tas.achieved_score, tas.percentage_score,
		       a.assessor_name, a.notes, at.name as assessment_type_name
		FROM assessments a
		JOIN thematic_area_scores tas ON a.id = tas.assessment_id
		JOIN assessment_types at ON a.assessment_type_id = at.id
		WHERE a.health_worker_id = $1
		AND tas.thematic_area_id IN (%s)
	`, assessmentPlaceholders)
	assessmentsArgs := []interface{}{healthWorkerID}
	for _, id := range allThematicAreaIDs {
		assessmentsArgs = append(assessmentsArgs, id)
	}
	assessmentsArgIdx := len(assessmentsArgs) + 1

	if startDate != "" {
		assessmentsQuery += fmt.Sprintf(" AND a.created_at >= $%d", assessmentsArgIdx)
		assessmentsArgs = append(assessmentsArgs, startDate)
		assessmentsArgIdx++
	}
	if endDate != "" {
		assessmentsQuery += fmt.Sprintf(" AND a.created_at <= $%d", assessmentsArgIdx)
		assessmentsArgs = append(assessmentsArgs, endDate+" 23:59:59")
		assessmentsArgIdx++
	}

	assessmentsQuery += " ORDER BY a.created_at DESC"

	assessmentRows, err := database.DB.Query(assessmentsQuery, assessmentsArgs...)
	if err != nil {
		return err
	}
	defer assessmentRows.Close()

	type AssessmentDetail struct {
		ID              int     `json:"id"`
		CreatedAt       string  `json:"createdAt"`
		PossibleScore   int     `json:"possibleScore"`
		AchievedScore   int     `json:"achievedScore"`
		PercentageScore float64 `json:"percentageScore"`
		AssessorName    string  `json:"assessorName"`
		Notes           string  `json:"notes"`
		AssessmentType  string  `json:"assessmentType"`
	}

	var assessments []AssessmentDetail
	var totalPossible, totalAchieved int
	var assessmentCount int

	for assessmentRows.Next() {
		var a AssessmentDetail
		var notes sql.NullString
		if err := assessmentRows.Scan(&a.ID, &a.CreatedAt, &a.PossibleScore, &a.AchievedScore, &a.PercentageScore, &a.AssessorName, &notes, &a.AssessmentType); err != nil {
			continue
		}
		if notes.Valid {
			a.Notes = notes.String
		}
		assessments = append(assessments, a)
		totalPossible += a.PossibleScore
		totalAchieved += a.AchievedScore
		assessmentCount++
	}

	percentage := 0.0
	if totalPossible > 0 {
		percentage = (float64(totalAchieved) / float64(totalPossible)) * 100
	}

	// Note: yesCount, noCount, naCount are already calculated above while loading responses

	// Group questions by their latest response type (each question appears only once)
	// Use questions from the map since that's where responses were added
	yesQuestions := []QuestionDetail{}
	noQuestions := []QuestionDetail{}
	naQuestions := []QuestionDetail{}
	filteredQuestions := []QuestionDetail{}

	for _, qPtr := range questionMap {
		// Only include questions that have responses
		if len(qPtr.Responses) == 0 {
			continue
		}

		q := *qPtr
		filteredQuestions = append(filteredQuestions, q)

		// Use latest response for grouping
		latestResponse := q.Responses[0].Response
		if latestResponse == "Yes" {
			yesQuestions = append(yesQuestions, q)
		} else if latestResponse == "No" {
			noQuestions = append(noQuestions, q)
		} else if latestResponse == "NA" {
			naQuestions = append(naQuestions, q)
		}
	}

	// Sort by display order
	sort.Slice(filteredQuestions, func(i, j int) bool {
		return filteredQuestions[i].DisplayOrder < filteredQuestions[j].DisplayOrder
	})
	sort.Slice(yesQuestions, func(i, j int) bool {
		return yesQuestions[i].DisplayOrder < yesQuestions[j].DisplayOrder
	})
	sort.Slice(noQuestions, func(i, j int) bool {
		return noQuestions[i].DisplayOrder < noQuestions[j].DisplayOrder
	})
	sort.Slice(naQuestions, func(i, j int) bool {
		return naQuestions[i].DisplayOrder < naQuestions[j].DisplayOrder
	})

	result := map[string]interface{}{
		"thematicArea": map[string]interface{}{
			"id":             taID,
			"name":           taName,
			"assessmentType": assessmentTypeName,
		},
		"performance": map[string]interface{}{
			"percentage":      percentage,
			"totalPossible":   totalPossible,
			"totalAchieved":   totalAchieved,
			"yesCount":        yesCount,
			"noCount":         noCount,
			"naCount":         naCount,
			"assessmentCount": assessmentCount,
		},
		"assessments":  assessments,
		"allQuestions": questions,
		"yesQuestions": yesQuestions,
		"noQuestions":  noQuestions,
		"naQuestions":  naQuestions,
	}

	return c.JSON(result)
}
