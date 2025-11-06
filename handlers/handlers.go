package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
)

// Geographic hierarchy handlers
func GetRegions(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT id, name FROM regions ORDER BY name")
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
	rows, err := database.DB.Query("SELECT id, name FROM districts WHERE region_id = $1 ORDER BY name", regionId)
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
	rows, err := database.DB.Query("SELECT id, name FROM subcounties WHERE district_id = $1 ORDER BY name", districtId)
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
	rows, err := database.DB.Query("SELECT id, name FROM facilities WHERE subcounty_id = $1 ORDER BY name", subcountyId)
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
	rows, err := database.DB.Query("SELECT id, name, code FROM assessment_types ORDER BY name")
	if err != nil {
		return err
	}
	defer rows.Close()

	var types []map[string]interface{}
	for rows.Next() {
		var id int
		var name, code string
		if err := rows.Scan(&id, &name, &code); err != nil {
			return err
		}
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
		SELECT id, name, display_order 
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
		SELECT id, question_text, score_weight, is_critical, is_important, display_order 
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
		// Overall possible always includes all questions
		totalPossible += weight

		if _, exists := thematicScores[thematicID]; !exists {
			thematicScores[thematicID] = struct{ possible, achieved int }{0, 0}
		}
		// For thematic possible, we'll initially include all, and subtract later if response is NA
		ts := thematicScores[thematicID]
		ts.possible += weight
		thematicScores[thematicID] = ts
	}
	rows.Close()

	// Calculate percentage
	percentage := 0.0
	if totalPossible > 0 {
		percentage = (float64(achieved) / float64(totalPossible)) * 100
	}

	// Determine performance level
	performanceLevel := "Not Acceptable"
	if percentage > 90 {
		performanceLevel = "Proficient"
	} else if percentage >= 70 {
		performanceLevel = "Competent"
	}

	// Insert assessment first
	var assessmentID int
	err = tx.QueryRow(`
		INSERT INTO assessments 
		(facility_id, assessment_type_id, assessor_name, client_name, notes, 
		 total_possible_score, achieved_score, percentage_score, performance_level)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, req.FacilityID, req.AssessmentTypeID, req.AssessorName, req.ClientName, req.Notes,
		totalPossible, achieved, percentage, performanceLevel).Scan(&assessmentID)
	if err != nil {
		return err
	}

	// Calculate achieved scores from responses and store them
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
			pointsEarned = q.weight
			achieved += q.weight

			ts := thematicScores[q.thematicID]
			ts.achieved += q.weight
			thematicScores[q.thematicID] = ts
		} else if response == "NA" {
			// Exclude NA from thematic-area possible score
			ts := thematicScores[q.thematicID]
			if ts.possible >= q.weight {
				ts.possible -= q.weight
			}
			thematicScores[q.thematicID] = ts
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

	// Recalculate and update assessment with correct achieved score
	percentage = 0.0
	if totalPossible > 0 {
		percentage = (float64(achieved) / float64(totalPossible)) * 100
	}

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

	// Build query with filters - join through hierarchy to get region/district info
	query := `
		SELECT a.id, a.created_at, a.percentage_score, a.performance_level,
		       f.name as facility_name, at.name as assessment_type,
		       a.assessor_name, a.client_name, 
		       r.id as region_id, d.id as district_id, s.id as subcounty_id, f.id as facility_id
		FROM assessments a
		JOIN facilities f ON a.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		JOIN assessment_types at ON a.assessment_type_id = at.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	// Apply admin area restrictions
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
	if assessmentTypeID := c.Query("assessmentTypeId"); assessmentTypeID != "" {
		query += fmt.Sprintf(" AND a.assessment_type_id = $%d", argIdx)
		args = append(args, assessmentTypeID)
		argIdx++
	}
	if thematicAreaID := c.Query("thematicAreaId"); thematicAreaID != "" {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM assessment_responses ar JOIN questions q ON ar.question_id = q.id WHERE ar.assessment_id = a.id AND q.thematic_area_id = $%d)", argIdx)
		args = append(args, thematicAreaID)
		argIdx++
	}

	query += " ORDER BY a.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var assessments []map[string]interface{}
	for rows.Next() {
		var id int
		var createdAt, assessorName, clientName sql.NullString
		var percentage float64
		var performanceLevel, facilityName, assessmentType string
		var regionID, districtID, subcountyID, facilityID sql.NullInt64
		if err := rows.Scan(&id, &createdAt, &percentage, &performanceLevel, &facilityName, &assessmentType, &assessorName, &clientName, &regionID, &districtID, &subcountyID, &facilityID); err != nil {
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
		})
	}

	return c.JSON(assessments)
}

func GetAssessment(c *fiber.Ctx) error {
	id := c.Params("id")

	var assessment struct {
		ID               int
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
	}

	err := database.DB.QueryRow(`
		SELECT a.id, f.name, at.name, a.assessor_name, a.client_name, a.notes,
		       a.total_possible_score, a.achieved_score, a.percentage_score,
		       a.performance_level, a.created_at
		FROM assessments a
		JOIN facilities f ON a.facility_id = f.id
		JOIN assessment_types at ON a.assessment_type_id = at.id
		WHERE a.id = $1
	`, id).Scan(
		&assessment.ID, &assessment.FacilityName, &assessment.AssessmentType,
		&assessment.AssessorName, &assessment.ClientName, &assessment.Notes,
		&assessment.TotalPossible, &assessment.Achieved, &assessment.Percentage,
		&assessment.PerformanceLevel, &assessment.CreatedAt,
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

	// Get assessment basic info
	var assessment struct {
		ID               int
		Percentage       float64
		PerformanceLevel string
		TotalPossible    int
		Achieved         int
	}

	err := database.DB.QueryRow(`
		SELECT id, percentage_score, performance_level, total_possible_score, achieved_score
		FROM assessments
		WHERE id = $1
	`, id).Scan(&assessment.ID, &assessment.Percentage, &assessment.PerformanceLevel,
		&assessment.TotalPossible, &assessment.Achieved)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "Assessment not found")
		}
		return err
	}

	// Get good and bad contributions
	goodRows, err := database.DB.Query(`
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
	defer goodRows.Close()

	var goodContributions []map[string]interface{}
	for goodRows.Next() {
		var text, taName string
		var weight int
		if err := goodRows.Scan(&text, &weight, &taName); err != nil {
			return err
		}
		goodContributions = append(goodContributions, map[string]interface{}{
			"question":     text,
			"points":       weight,
			"thematicArea": taName,
		})
	}

	badRows, err := database.DB.Query(`
		SELECT q.question_text, q.score_weight, ta.name as thematic_area, ar.response
		FROM assessment_responses ar
		JOIN questions q ON ar.question_id = q.id
		JOIN thematic_areas ta ON q.thematic_area_id = ta.id
		WHERE ar.assessment_id = $1 AND ar.response IN ('No', 'NA')
		ORDER BY q.score_weight DESC, ta.display_order
	`, id)
	if err != nil {
		return err
	}
	defer badRows.Close()

	var badContributions []map[string]interface{}
	for badRows.Next() {
		var text, taName, response string
		var weight int
		if err := badRows.Scan(&text, &weight, &taName, &response); err != nil {
			return err
		}
		badContributions = append(badContributions, map[string]interface{}{
			"question":     text,
			"points":       weight,
			"thematicArea": taName,
			"response":     response,
		})
	}

	return c.JSON(fiber.Map{
		"assessment":        assessment,
		"goodContributions": goodContributions,
		"badContributions":  badContributions,
	})
}
