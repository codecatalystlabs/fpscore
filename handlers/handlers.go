package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
)

// GetUserAdminAreasWithDetails returns user's admin areas with full hierarchy details
func GetUserAdminAreasWithDetails(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	result := map[string]interface{}{
		"isAdmin": isAdmin,
		"restrictions": map[string]interface{}{
			"regionId":         nil,
			"regionName":       nil,
			"districtId":       nil,
			"districtName":     nil,
			"subcountyId":      nil,
			"subcountyName":    nil,
			"facilityId":       nil,
			"facilityName":     nil,
			"restrictionLevel": nil, // "region", "district", "subcounty", "facility", or null
		},
	}

	if isAdmin {
		return c.JSON(result)
	}

	// Get all user's admin areas (not just the first row)
	rows, err := database.DB.Query(`
		SELECT region_id, district_id, subcounty_id, facility_id
		FROM user_admin_areas
		WHERE user_id = $1
	`, userID)
	if err != nil {
		fmt.Printf("ERROR: Failed to query user_admin_areas for user %d: %v\n", userID, err)
		return err
	}
	defer rows.Close()

	restrictions := result["restrictions"].(map[string]interface{})
	regionSet := map[int]bool{}
	districtSet := map[int]bool{}
	subcountySet := map[int]bool{}
	facilitySet := map[int]bool{}

	for rows.Next() {
		var regionID, districtID, subcountyID, facilityID sql.NullInt64
		if err := rows.Scan(&regionID, &districtID, &subcountyID, &facilityID); err != nil {
			fmt.Printf("ERROR: Failed to scan user_admin_areas for user %d: %v\n", userID, err)
			return err
		}
		if regionID.Valid {
			regionSet[int(regionID.Int64)] = true
		}
		if districtID.Valid {
			districtSet[int(districtID.Int64)] = true
		}
		if subcountyID.Valid {
			subcountySet[int(subcountyID.Int64)] = true
		}
		if facilityID.Valid {
			facilitySet[int(facilityID.Int64)] = true
		}
	}

	regionIDs := make([]int, 0, len(regionSet))
	for id := range regionSet {
		regionIDs = append(regionIDs, id)
	}
	districtIDs := make([]int, 0, len(districtSet))
	for id := range districtSet {
		districtIDs = append(districtIDs, id)
	}
	subcountyIDs := make([]int, 0, len(subcountySet))
	for id := range subcountySet {
		subcountyIDs = append(subcountyIDs, id)
	}
	facilityIDs := make([]int, 0, len(facilitySet))
	for id := range facilitySet {
		facilityIDs = append(facilityIDs, id)
	}

	// Expose full sets so frontend can support multi-area restrictions
	restrictions["regionIds"] = regionIDs
	restrictions["districtIds"] = districtIDs
	restrictions["subcountyIds"] = subcountyIDs
	restrictions["facilityIds"] = facilityIDs

	// If user has multiple admin areas at the same/higher level, do not hard-lock to a single area.
	// API endpoints already enforce full area restrictions; this prevents UI from narrowing to only the first entry.
	if len(facilityIDs) > 1 || len(subcountyIDs) > 1 || len(districtIDs) > 1 || len(regionIDs) > 1 {
		return c.JSON(result)
	}

	// Single-scope restriction handling (preserve old behavior for exactly one assigned area)
	if len(facilityIDs) == 1 {
		facilityID := facilityIDs[0]
		var facilityName, subcountyName, districtName, regionName string
		var subcountyIDResolved, districtIDResolved, regionIDResolved int
		err := database.DB.QueryRow(`
			SELECT f.name, s.id, s.name, d.id, d.name, r.id, r.name
			FROM facilities f
			JOIN subcounties s ON f.subcounty_id = s.id
			JOIN districts d ON s.district_id = d.id
			JOIN regions r ON d.region_id = r.id
			WHERE f.id = $1
		`, facilityID).Scan(&facilityName, &subcountyIDResolved, &subcountyName,
			&districtIDResolved, &districtName, &regionIDResolved, &regionName)
		if err == nil {
			restrictions["facilityId"] = facilityID
			restrictions["facilityName"] = facilityName
			restrictions["subcountyId"] = subcountyIDResolved
			restrictions["subcountyName"] = subcountyName
			restrictions["districtId"] = districtIDResolved
			restrictions["districtName"] = districtName
			restrictions["regionId"] = regionIDResolved
			restrictions["regionName"] = regionName
			restrictions["restrictionLevel"] = "facility"
		}
	} else if len(subcountyIDs) == 1 {
		subcountyID := subcountyIDs[0]
		var subcountyName, districtName, regionName string
		var districtIDResolved, regionIDResolved int
		err := database.DB.QueryRow(`
			SELECT s.name, d.id, d.name, r.id, r.name
			FROM subcounties s
			JOIN districts d ON s.district_id = d.id
			JOIN regions r ON d.region_id = r.id
			WHERE s.id = $1
		`, subcountyID).Scan(&subcountyName, &districtIDResolved, &districtName, &regionIDResolved, &regionName)
		if err == nil {
			restrictions["subcountyId"] = subcountyID
			restrictions["subcountyName"] = subcountyName
			restrictions["districtId"] = districtIDResolved
			restrictions["districtName"] = districtName
			restrictions["regionId"] = regionIDResolved
			restrictions["regionName"] = regionName
			restrictions["restrictionLevel"] = "subcounty"
		}
	} else if len(districtIDs) == 1 {
		districtID := districtIDs[0]
		var districtName, regionName string
		var regionIDResolved int
		err := database.DB.QueryRow(`
			SELECT d.name, r.id, r.name
			FROM districts d
			JOIN regions r ON d.region_id = r.id
			WHERE d.id = $1
		`, districtID).Scan(&districtName, &regionIDResolved, &regionName)
		if err == nil {
			restrictions["districtId"] = districtID
			restrictions["districtName"] = districtName
			restrictions["regionId"] = regionIDResolved
			restrictions["regionName"] = regionName
			restrictions["restrictionLevel"] = "district"
		}
	} else if len(regionIDs) == 1 {
		regionID := regionIDs[0]
		var regionName string
		err := database.DB.QueryRow("SELECT name FROM regions WHERE id = $1", regionID).Scan(&regionName)
		if err == nil {
			restrictions["regionId"] = regionID
			restrictions["regionName"] = regionName
			restrictions["restrictionLevel"] = "region"
		}
	}

	return c.JSON(result)
}

// Geographic hierarchy handlers
func GetRegions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	query := "SELECT id, name FROM regions"
	args := []interface{}{}

	if !isAdmin {
		regionIDs, _, _, _, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(regionIDs) > 0 {
			placeholders := ""
			for i, id := range regionIDs {
				if i > 0 {
					placeholders += ","
				}
				placeholders += fmt.Sprintf("$%d", i+1)
				args = append(args, id)
			}
			query += fmt.Sprintf(" WHERE id IN (%s)", placeholders)
		} else {
			// User has no admin areas - return empty
			return c.JSON([]map[string]interface{}{})
		}
	}

	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var regions []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		regions = append(regions, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	return c.JSON(regions)
}

func GetDistricts(c *fiber.Ctx) error {
	regionId := c.Params("regionId")
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	query := "SELECT id, name FROM districts WHERE region_id = $1"
	args := []interface{}{regionId}
	argIdx := 2

	if !isAdmin {
		_, districtIDs, _, _, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(districtIDs) > 0 {
			placeholders := ""
			for i, id := range districtIDs {
				if i > 0 {
					placeholders += ","
				}
				placeholders += fmt.Sprintf("$%d", argIdx)
				args = append(args, id)
				argIdx++
			}
			query += fmt.Sprintf(" AND id IN (%s)", placeholders)
		} else {
			// User has no admin areas - return empty
			return c.JSON([]map[string]interface{}{})
		}
	}

	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var districts []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		districts = append(districts, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	return c.JSON(districts)
}

func GetSubcounties(c *fiber.Ctx) error {
	districtId := c.Params("districtId")
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	query := "SELECT id, name FROM subcounties WHERE district_id = $1"
	args := []interface{}{districtId}
	argIdx := 2

	if !isAdmin {
		_, districtIDs, subcountyIDs, _, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}

		// Check if user has specific subcounty restrictions
		if len(subcountyIDs) > 0 {
			// User is restricted to specific subcounties - filter by those
			placeholders := ""
			for i, id := range subcountyIDs {
				if i > 0 {
					placeholders += ","
				}
				placeholders += fmt.Sprintf("$%d", argIdx)
				args = append(args, id)
				argIdx++
			}
			query += fmt.Sprintf(" AND id IN (%s)", placeholders)
		} else if len(districtIDs) > 0 {
			// User is restricted to districts but not specific subcounties
			// Check if the requested district is in their allowed districts
			districtIdInt, err := strconv.Atoi(districtId)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid district ID")
			}
			districtAllowed := false
			for _, id := range districtIDs {
				if id == districtIdInt {
					districtAllowed = true
					break
				}
			}
			if !districtAllowed {
				// User doesn't have access to this district
				return c.JSON([]map[string]interface{}{})
			}
			// User has access to this district - return all subcounties in it (no additional filtering)
		} else {
			// User has no admin areas - return empty
			return c.JSON([]map[string]interface{}{})
		}
	}

	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var subcounties []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		subcounties = append(subcounties, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	return c.JSON(subcounties)
}

func GetFacilities(c *fiber.Ctx) error {
	subcountyId := c.Params("subcountyId")
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	query := "SELECT id, name FROM facilities WHERE subcounty_id = $1"
	args := []interface{}{subcountyId}
	argIdx := 2

	if !isAdmin {
		_, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}

		// Check if user has specific facility restrictions
		if len(facilityIDs) > 0 {
			// User is restricted to specific facilities - filter by those
			placeholders := ""
			for i, id := range facilityIDs {
				if i > 0 {
					placeholders += ","
				}
				placeholders += fmt.Sprintf("$%d", argIdx)
				args = append(args, id)
				argIdx++
			}
			query += fmt.Sprintf(" AND id IN (%s)", placeholders)
		} else if len(subcountyIDs) > 0 {
			// User is restricted to subcounties but not specific facilities
			// Check if the requested subcounty is in their allowed subcounties
			subcountyIdInt, err := strconv.Atoi(subcountyId)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid subcounty ID")
			}
			subcountyAllowed := false
			for _, id := range subcountyIDs {
				if id == subcountyIdInt {
					subcountyAllowed = true
					break
				}
			}
			if !subcountyAllowed {
				// User doesn't have access to this subcounty
				return c.JSON([]map[string]interface{}{})
			}
			// User has access to this subcounty - return all facilities in it (no additional filtering)
		} else if len(districtIDs) > 0 {
			// User is restricted to districts - need to check if subcounty belongs to their district
			subcountyIdInt, err := strconv.Atoi(subcountyId)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid subcounty ID")
			}
			// Get the district for this subcounty
			var subcountyDistrictID int
			err = database.DB.QueryRow("SELECT district_id FROM subcounties WHERE id = $1", subcountyIdInt).Scan(&subcountyDistrictID)
			if err != nil {
				return fiber.NewError(fiber.StatusNotFound, "Subcounty not found")
			}
			// Check if this district is in user's allowed districts
			districtAllowed := false
			for _, id := range districtIDs {
				if id == subcountyDistrictID {
					districtAllowed = true
					break
				}
			}
			if !districtAllowed {
				// User doesn't have access to this district/subcounty
				return c.JSON([]map[string]interface{}{})
			}
			// User has access to this district - return all facilities in the subcounty (no additional filtering)
		} else {
			// User has no admin areas - return empty
			return c.JSON([]map[string]interface{}{})
		}
	}

	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var facilities []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		facilities = append(facilities, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	return c.JSON(facilities)
}

// Assessment type handlers
func GetAssessmentTypes(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT DISTINCT id, name, code FROM assessment_types ORDER BY name")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Use a map to deduplicate by id in case of any duplicates
	seenIds := make(map[int]bool)
	var types []map[string]interface{}
	for rows.Next() {
		var id int
		var name, code string
		if err := rows.Scan(&id, &name, &code); err != nil {
			return err
		}
		// Skip if we've already seen this ID
		if seenIds[id] {
			continue
		}
		seenIds[id] = true
		types = append(types, map[string]interface{}{
			"id":   id,
			"name": name,
			"code": code,
		})
	}

	return c.JSON(types)
}

func GetThematicAreas(c *fiber.Ctx) error {
	typeId := c.Params("typeId")
	rows, err := database.DB.Query(`
		SELECT 
			id, 
			name, 
			COALESCE(display_order, 0) AS display_order
		FROM thematic_areas 
		WHERE assessment_type_id = $1 
		ORDER BY display_order, name
	`, typeId)
	if err != nil {
		return err
	}
	defer rows.Close()

	var areas []map[string]interface{}
	for rows.Next() {
		var id, order int
		var name string
		if err := rows.Scan(&id, &name, &order); err != nil {
			return err
		}
		areas = append(areas, map[string]interface{}{
			"id":           id,
			"name":         name,
			"displayOrder": order,
		})
	}

	return c.JSON(areas)
}

func GetQuestions(c *fiber.Ctx) error {
	thematicAreaId := c.Params("thematicAreaId")
	rows, err := database.DB.Query(`
		SELECT 
			id, 
			question_text, 
			COALESCE(score_weight, 0)     AS score_weight, 
			is_critical, 
			is_important, 
			COALESCE(display_order, 0)    AS display_order
		FROM questions 
		WHERE thematic_area_id = $1 
		ORDER BY display_order
	`, thematicAreaId)
	if err != nil {
		return err
	}
	defer rows.Close()

	var questions []map[string]interface{}
	for rows.Next() {
		var id, weight, order int
		var text string
		var critical, important bool
		if err := rows.Scan(&id, &text, &weight, &critical, &important, &order); err != nil {
			return err
		}
		questions = append(questions, map[string]interface{}{
			"id":           id,
			"text":         text,
			"scoreWeight":  weight,
			"isCritical":   critical,
			"isImportant":  important,
			"displayOrder": order,
		})
	}

	return c.JSON(questions)
}

// Assessment handlers
type CreateAssessmentRequest struct {
	HealthWorkerID   int               `json:"healthWorkerId"`
	FacilityID       int               `json:"facilityId"`
	AssessmentTypeID int               `json:"assessmentTypeId"`
	AssessorName     string            `json:"assessorName"`
	ClientName       string            `json:"clientName"`
	Notes            string            `json:"notes"`
	Responses        map[string]string `json:"responses"` // question_id -> "Yes"/"No"/"NA"
}

func CreateAssessment(c *fiber.Ctx) error {
	var req CreateAssessmentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "Invalid request body")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Calculate scores
	var totalPossible, achieved int
	thematicScores := make(map[int]struct{ possible, achieved int })

	// Get all questions for this assessment type
	rows, err := tx.Query(`
		SELECT q.id, q.thematic_area_id, q.score_weight
		FROM questions q
		JOIN thematic_areas ta ON q.thematic_area_id = ta.id
		WHERE ta.assessment_type_id = $1
	`, req.AssessmentTypeID)
	if err != nil {
		return err
	}

	questions := make(map[int]struct {
		thematicID int
		weight     int
	})

	for rows.Next() {
		var qID, thematicID, weight int
		if err := rows.Scan(&qID, &thematicID, &weight); err != nil {
			rows.Close()
			return err
		}
		questions[qID] = struct {
			thematicID int
			weight     int
		}{thematicID, weight}

		if _, exists := thematicScores[thematicID]; !exists {
			thematicScores[thematicID] = struct{ possible, achieved int }{0, 0}
		}
	}
	rows.Close()

	// Validate health worker exists and get their facility
	var healthWorkerFacilityID int
	err = tx.QueryRow("SELECT facility_id FROM health_workers WHERE id = $1", req.HealthWorkerID).Scan(&healthWorkerFacilityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(400, "Health worker not found")
		}
		return err
	}
	// Use the health worker's facility (they may have moved, so use current facility)
	req.FacilityID = healthWorkerFacilityID

	// Insert assessment first (we'll update scores after processing responses)
	var assessmentID int
	err = tx.QueryRow(`
		INSERT INTO assessments 
		(health_worker_id, facility_id, assessment_type_id, assessor_name, client_name, notes, 
		 total_possible_score, achieved_score, percentage_score, performance_level)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, req.HealthWorkerID, req.FacilityID, req.AssessmentTypeID, req.AssessorName, req.ClientName, req.Notes,
		0, 0, 0.0, "Not Acceptable").Scan(&assessmentID)
	if err != nil {
		return err
	}

	// Calculate achieved scores from responses and store them
	// Only count questions that are NOT NA toward total possible
	for qIDStr, response := range req.Responses {
		// Parse question ID from string key
		qID, err := strconv.Atoi(qIDStr)
		if err != nil {
			continue
		}

		q, exists := questions[qID]
		if !exists {
			continue
		}

		pointsEarned := 0
		if response == "Yes" {
			// Only "Yes" responses get points
			pointsEarned = q.weight
			achieved += q.weight
			totalPossible += q.weight // Count toward total possible

			ts := thematicScores[q.thematicID]
			ts.achieved += q.weight
			ts.possible += q.weight // Count toward thematic possible
			thematicScores[q.thematicID] = ts
		} else if response == "No" {
			// "No" responses get 0 points but still count toward possible
			pointsEarned = 0
			totalPossible += q.weight // Count toward total possible

			ts := thematicScores[q.thematicID]
			ts.possible += q.weight // Count toward thematic possible
			thematicScores[q.thematicID] = ts
		} else if response == "NA" {
			// NA responses don't count toward possible scores at all
			pointsEarned = 0
			// Don't add to totalPossible or thematic possible
		}

		// Store response
		_, err = tx.Exec(`
			INSERT INTO assessment_responses (assessment_id, question_id, response, points_earned)
			VALUES ($1, $2, $3, $4)
        `, assessmentID, qID, response, pointsEarned)
		if err != nil {
			return err
		}
	}

	// Calculate percentage and performance level
	percentage := 0.0
	if totalPossible > 0 {
		percentage = (float64(achieved) / float64(totalPossible)) * 100
	}

	performanceLevel := "Not Acceptable"
	if percentage > 90 {
		performanceLevel = "Proficient"
	} else if percentage >= 70 {
		performanceLevel = "Competent"
	} else {
		performanceLevel = "Not Acceptable"
	}

	_, err = tx.Exec(`
		UPDATE assessments 
		SET achieved_score = $1, percentage_score = $2, performance_level = $3
		WHERE id = $4
	`, achieved, percentage, performanceLevel, assessmentID)
	if err != nil {
		return err
	}

	// Insert thematic area scores
	for thematicID, scores := range thematicScores {
		taPercentage := 0.0
		if scores.possible > 0 {
			taPercentage = (float64(scores.achieved) / float64(scores.possible)) * 100
		}

		_, err = tx.Exec(`
			INSERT INTO thematic_area_scores 
			(assessment_id, thematic_area_id, possible_score, achieved_score, percentage_score)
			VALUES ($1, $2, $3, $4, $5)
		`, assessmentID, thematicID, scores.possible, scores.achieved, taPercentage)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id":               assessmentID,
		"totalPossible":    totalPossible,
		"achieved":         achieved,
		"percentage":       percentage,
		"performanceLevel": performanceLevel,
	})
}

func GetAssessments(c *fiber.Ctx) error {
	// Get user's admin areas for filtering
	userID := c.Locals("userID").(int)

	// Check if user is admin (has admin role or all permissions)
	isAdmin := IsAdmin(userID)

	// Get user's admin areas (only if not admin)
	var regionIDs, districtIDs, subcountyIDs, facilityIDs []int
	var err error
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err = GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
	}
	// If admin, these arrays remain empty, allowing access to all assessments

	includeThematic := c.Query("includeThematicScores") == "1" || strings.ToLower(c.Query("includeThematicScores")) == "true"

	// Build query with filters - join through hierarchy to get region/district info.
	// Optionally include thematic area scores so the dashboard can visualize by thematic area.
	query := `
		SELECT a.id, a.created_at, a.percentage_score, a.performance_level,
		       f.name as facility_name, at.name as assessment_type,
		       a.assessor_name, a.client_name, 
		       r.id as region_id, d.id as district_id, s.id as subcounty_id, f.id as facility_id,
		       hw.id as health_worker_id, hw.full_name as health_worker_name
	`
	if includeThematic {
		query += `,
		       ta.id as thematic_area_id, ta.name as thematic_area_name,
		       tas.percentage_score as thematic_percentage_score,
		       tas.possible_score as thematic_possible_score,
		       tas.achieved_score as thematic_achieved_score
		`
	}
	query += `
		FROM assessments a
		JOIN health_workers hw ON a.health_worker_id = hw.id
		JOIN facilities f ON a.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		JOIN assessment_types at ON a.assessment_type_id = at.id
	`
	if includeThematic {
		query += `
		JOIN thematic_area_scores tas ON tas.assessment_id = a.id
		JOIN thematic_areas ta ON tas.thematic_area_id = ta.id
		`
	}
	query += `
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	// Apply admin area restrictions
	// If user is not admin, they must have admin areas assigned to see any assessments
	if !isAdmin {
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
			// Non-admin user with no admin areas assigned - return no results
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
		query += fmt.Sprintf(" AND f.id = $%d", argIdx)
		args = append(args, facilityID)
		argIdx++
	}
	if healthWorkerID := c.Query("healthWorkerId"); healthWorkerID != "" {
		query += fmt.Sprintf(" AND hw.id = $%d", argIdx)
		args = append(args, healthWorkerID)
		argIdx++
	}
	if assessmentTypeID := c.Query("assessmentTypeId"); assessmentTypeID != "" {
		query += fmt.Sprintf(" AND a.assessment_type_id = $%d", argIdx)
		args = append(args, assessmentTypeID)
		argIdx++
	}
	if thematicAreaID := c.Query("thematicAreaId"); thematicAreaID != "" {
		if includeThematic {
			query += fmt.Sprintf(" AND ta.id = $%d", argIdx)
		} else {
			query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM assessment_responses ar JOIN questions q ON ar.question_id = q.id WHERE ar.assessment_id = a.id AND q.thematic_area_id = $%d)", argIdx)
		}
		args = append(args, thematicAreaID)
		argIdx++
	}
	if startDate := c.Query("startDate"); startDate != "" {
		query += fmt.Sprintf(" AND DATE(a.created_at) >= $%d", argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate := c.Query("endDate"); endDate != "" {
		query += fmt.Sprintf(" AND DATE(a.created_at) <= $%d", argIdx)
		args = append(args, endDate)
		argIdx++
	}

	query += " ORDER BY a.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	assessments := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var createdAt, assessorName, clientName sql.NullString
		var percentage float64
		var performanceLevel, facilityName, assessmentType, healthWorkerName string
		var regionID, districtID, subcountyID, facilityID, healthWorkerID sql.NullInt64
		if includeThematic {
			var thematicAreaID sql.NullInt64
			var thematicAreaName sql.NullString
			var thematicPercentage sql.NullFloat64
			var thematicPossible, thematicAchieved sql.NullInt64
			if err := rows.Scan(
				&id, &createdAt, &percentage, &performanceLevel,
				&facilityName, &assessmentType,
				&assessorName, &clientName,
				&regionID, &districtID, &subcountyID, &facilityID,
				&healthWorkerID, &healthWorkerName,
				&thematicAreaID, &thematicAreaName,
				&thematicPercentage, &thematicPossible, &thematicAchieved,
			); err != nil {
				return err
			}
			assessments = append(assessments, map[string]interface{}{
				"id":                     id,
				"createdAt":              createdAt.String,
				"percentage":             percentage,
				"performanceLevel":       performanceLevel,
				"facilityName":           facilityName,
				"assessmentType":         assessmentType,
				"assessorName":           assessorName.String,
				"clientName":             clientName.String,
				"healthWorkerId":         healthWorkerID.Int64,
				"healthWorkerName":       healthWorkerName,
				"regionId":               regionID.Int64,
				"districtId":             districtID.Int64,
				"subcountyId":            subcountyID.Int64,
				"facilityId":             facilityID.Int64,
				"thematicAreaId":         thematicAreaID.Int64,
				"thematicArea":           thematicAreaName.String,
				"thematicPercentageScore": func() float64 {
					if thematicPercentage.Valid {
						return thematicPercentage.Float64
					}
					return 0
				}(),
				"thematicPossibleScore": func() int64 {
					if thematicPossible.Valid {
						return thematicPossible.Int64
					}
					return 0
				}(),
				"thematicAchievedScore": func() int64 {
					if thematicAchieved.Valid {
						return thematicAchieved.Int64
					}
					return 0
				}(),
			})
		} else {
			if err := rows.Scan(&id, &createdAt, &percentage, &performanceLevel, &facilityName, &assessmentType, &assessorName, &clientName, &regionID, &districtID, &subcountyID, &facilityID, &healthWorkerID, &healthWorkerName); err != nil {
				return err
			}

			assessments = append(assessments, map[string]interface{}{
				"id":               id,
				"createdAt":        createdAt.String,
				"percentage":       percentage,
				"performanceLevel": performanceLevel,
				"facilityName":     facilityName,
				"assessmentType":   assessmentType,
				"assessorName":     assessorName.String,
				"clientName":       clientName.String,
				"healthWorkerId":   healthWorkerID.Int64,
				"healthWorkerName": healthWorkerName,
				"regionId":         regionID.Int64,
				"districtId":       districtID.Int64,
				"subcountyId":      subcountyID.Int64,
				"facilityId":       facilityID.Int64,
			})
		}
	}

	return c.JSON(assessments)
}

func GetAssessment(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(int)

	// Check if user is admin
	isAdmin := IsAdmin(userID)

	// Get user's admin areas (only if not admin)
	var regionIDs, districtIDs, subcountyIDs, facilityIDs []int
	var err error
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err = GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
	}

	var assessment struct {
		ID               int
		HealthWorkerName string
		FacilityName     string
		AssessmentType   string
		AssessorName     sql.NullString
		ClientName       sql.NullString
		Notes            sql.NullString
		TotalPossible    int
		Achieved         int
		Percentage       float64
		PerformanceLevel string
		CreatedAt        string
		RegionID         sql.NullInt64
		DistrictID       sql.NullInt64
		SubcountyID      sql.NullInt64
		FacilityID       sql.NullInt64
	}

	// Build query with admin area check
	query := `
		SELECT a.id, hw.full_name, f.name, at.name, a.assessor_name, a.client_name, a.notes,
		       a.total_possible_score, a.achieved_score, a.percentage_score,
		       a.performance_level, a.created_at,
		       r.id as region_id, d.id as district_id, s.id as subcounty_id, f.id as facility_id
		FROM assessments a
		JOIN health_workers hw ON a.health_worker_id = hw.id
		JOIN facilities f ON a.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		JOIN assessment_types at ON a.assessment_type_id = at.id
		WHERE a.id = $1
	`
	args := []interface{}{id}
	argIdx := 2

	// Apply admin area restrictions
	if !isAdmin {
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
			// Non-admin user with no admin areas assigned - return not found
			return fiber.NewError(404, "Assessment not found")
		}
	}

	err = database.DB.QueryRow(query, args...).Scan(
		&assessment.ID, &assessment.HealthWorkerName, &assessment.FacilityName, &assessment.AssessmentType,
		&assessment.AssessorName, &assessment.ClientName, &assessment.Notes,
		&assessment.TotalPossible, &assessment.Achieved, &assessment.Percentage,
		&assessment.PerformanceLevel, &assessment.CreatedAt,
		&assessment.RegionID, &assessment.DistrictID, &assessment.SubcountyID, &assessment.FacilityID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Assessment not found")
		}
		return err
	}

	// Get responses with questions
	rows, err := database.DB.Query(`
		SELECT q.id, q.question_text, q.score_weight, ar.response, ar.points_earned,
		       ta.id as thematic_area_id, ta.name as thematic_area_name
		FROM assessment_responses ar
		JOIN questions q ON ar.question_id = q.id
		JOIN thematic_areas ta ON q.thematic_area_id = ta.id
		WHERE ar.assessment_id = $1
		ORDER BY ta.display_order, q.display_order
	`, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	var responses []map[string]interface{}
	for rows.Next() {
		var qID, weight, pointsEarned, taID int
		var text, response, taName string
		if err := rows.Scan(&qID, &text, &weight, &response, &pointsEarned, &taID, &taName); err != nil {
			return err
		}
		responses = append(responses, map[string]interface{}{
			"questionId":       qID,
			"questionText":     text,
			"scoreWeight":      weight,
			"response":         response,
			"pointsEarned":     pointsEarned,
			"thematicAreaId":   taID,
			"thematicAreaName": taName,
		})
	}

	// Get thematic area scores
	taRows, err := database.DB.Query(`
		SELECT ta.id, ta.name, tas.possible_score, tas.achieved_score, tas.percentage_score
		FROM thematic_area_scores tas
		JOIN thematic_areas ta ON tas.thematic_area_id = ta.id
		WHERE tas.assessment_id = $1
		ORDER BY ta.display_order
	`, id)
	if err != nil {
		return err
	}
	defer taRows.Close()

	var thematicScores []map[string]interface{}
	for taRows.Next() {
		var taID, possible, achieved int
		var taName string
		var percentage float64
		if err := taRows.Scan(&taID, &taName, &possible, &achieved, &percentage); err != nil {
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

	result := map[string]interface{}{
		"assessment":     assessment,
		"responses":      responses,
		"thematicScores": thematicScores,
	}

	return c.JSON(result)
}

func GetAssessmentSummary(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(int)

	// Check if user is admin
	isAdmin := IsAdmin(userID)

	// Get user's admin areas (only if not admin)
	var regionIDs, districtIDs, subcountyIDs, facilityIDs []int
	var err error
	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err = GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
	}

	// Build query with admin area check
	query := `
		SELECT a.id, a.percentage_score, a.performance_level, a.total_possible_score, a.achieved_score
		FROM assessments a
		JOIN facilities f ON a.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE a.id = $1
	`
	args := []interface{}{id}
	argIdx := 2

	// Apply admin area restrictions
	if !isAdmin {
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
			// Non-admin user with no admin areas assigned - return not found
			return fiber.NewError(404, "Assessment not found")
		}
	}

	// Get assessment basic info
	var assessment struct {
		ID               int
		Percentage       float64
		PerformanceLevel string
		TotalPossible    int
		Achieved         int
	}

	err = database.DB.QueryRow(query, args...).Scan(&assessment.ID, &assessment.Percentage, &assessment.PerformanceLevel,
		&assessment.TotalPossible, &assessment.Achieved)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Assessment not found")
		}
		return err
	}

	// Get Yes responses
	yesRows, err := database.DB.Query(`
		SELECT q.question_text, q.score_weight, ta.name as thematic_area
		FROM assessment_responses ar
		JOIN questions q ON ar.question_id = q.id
		JOIN thematic_areas ta ON q.thematic_area_id = ta.id
		WHERE ar.assessment_id = $1 AND ar.response = 'Yes'
		ORDER BY q.score_weight DESC, ta.display_order
	`, id)
	if err != nil {
		return err
	}
	defer yesRows.Close()

	var yesResponses []map[string]interface{}
	for yesRows.Next() {
		var text, taName string
		var weight int
		if err := yesRows.Scan(&text, &weight, &taName); err != nil {
			return err
		}
		yesResponses = append(yesResponses, map[string]interface{}{
			"question":     text,
			"points":       weight,
			"thematicArea": taName,
		})
	}

	// Get No responses
	noRows, err := database.DB.Query(`
		SELECT q.question_text, q.score_weight, ta.name as thematic_area
		FROM assessment_responses ar
		JOIN questions q ON ar.question_id = q.id
		JOIN thematic_areas ta ON q.thematic_area_id = ta.id
		WHERE ar.assessment_id = $1 AND ar.response = 'No'
		ORDER BY q.score_weight DESC, ta.display_order
	`, id)
	if err != nil {
		return err
	}
	defer noRows.Close()

	var noResponses []map[string]interface{}
	for noRows.Next() {
		var text, taName string
		var weight int
		if err := noRows.Scan(&text, &weight, &taName); err != nil {
			return err
		}
		noResponses = append(noResponses, map[string]interface{}{
			"question":     text,
			"points":       weight,
			"thematicArea": taName,
		})
	}

	// Get N/A responses
	naRows, err := database.DB.Query(`
		SELECT q.question_text, q.score_weight, ta.name as thematic_area
		FROM assessment_responses ar
		JOIN questions q ON ar.question_id = q.id
		JOIN thematic_areas ta ON q.thematic_area_id = ta.id
		WHERE ar.assessment_id = $1 AND ar.response = 'NA'
		ORDER BY q.score_weight DESC, ta.display_order
	`, id)
	if err != nil {
		return err
	}
	defer naRows.Close()

	var naResponses []map[string]interface{}
	for naRows.Next() {
		var text, taName string
		var weight int
		if err := naRows.Scan(&text, &weight, &taName); err != nil {
			return err
		}
		naResponses = append(naResponses, map[string]interface{}{
			"question":     text,
			"points":       weight,
			"thematicArea": taName,
		})
	}

	return c.JSON(fiber.Map{
		"assessment":   assessment,
		"yesResponses": yesResponses,
		"noResponses":  noResponses,
		"naResponses":  naResponses,
	})
}
