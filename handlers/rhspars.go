package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
)

// RH SPARS scoring (strict):
//   Yes -> earn score_weight; count toward possible
//   No  -> earn 0; count toward possible
//   NA  -> ignored (not in possible or achieved)
// Section %  = achieved/possible * 100
// Domain %   = average of section percentages in that domain (sections with possible>0)
// Grand %    = average of the 6 domain percentages
// Performance level (same as FP proficiency): Proficient >90%, Competent 70–89%, Not Acceptable <70%

type RHSParsCreateRequest struct {
	FacilityID               int               `json:"facilityId"`
	SupervisionType          string            `json:"supervisionType"`
	SupervisionDate          string            `json:"supervisionDate"`
	NextSupervisionDate      string            `json:"nextSupervisionDate"`
	RHSupplier               string            `json:"rhSupplier"`
	FacilityTypeDetail       string            `json:"facilityTypeDetail"`
	FacilityOwnership        string            `json:"facilityOwnership"`
	PrimaryContactName       string            `json:"primaryContactName"`
	PrimaryContactPhone      string            `json:"primaryContactPhone"`
	PrimaryContactEmail      string            `json:"primaryContactEmail"`
	CompletedByName          string            `json:"completedByName"`
	CompletedByDesignation   string            `json:"completedByDesignation"`
	CompletedByInstitution   string            `json:"completedByInstitution"`
	AssessmentCompletionDate string            `json:"assessmentCompletionDate"`
	Notes                    string            `json:"notes"`
	Responses                map[string]string `json:"responses"` // question_id -> Yes/No/NA
	Attendees                []struct {
		PersonType    string `json:"personType"`
		FullName      string `json:"fullName"`
		CadrePosition string `json:"cadrePosition"`
		Affiliation   string `json:"affiliation"`
		Contact       string `json:"contact"`
	} `json:"attendees"`
	ActionPlan []struct {
		GapIdentified        string `json:"gapIdentified"`
		ActionRecommendation string `json:"actionRecommendation"`
		MeansOfVerification  string `json:"meansOfVerification"`
		ResponsiblePerson    string `json:"responsiblePerson"`
		TimeFrame            string `json:"timeFrame"`
	} `json:"actionPlan"`
}

func nullDate(s string) interface{} {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return s
}

func GetRHSParsDomains(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
		SELECT id, code, name, display_order
		FROM rhspars_domains
		ORDER BY display_order, name
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	out := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, order int
		var code, name string
		if err := rows.Scan(&id, &code, &name, &order); err != nil {
			return err
		}
		out = append(out, map[string]interface{}{
			"id": id, "code": code, "name": name, "displayOrder": order,
		})
	}
	return c.JSON(out)
}

func GetRHSParsThematicAreas(c *fiber.Ctx) error {
	domainID := c.Params("domainId")
	rows, err := database.DB.Query(`
		SELECT id, name, display_order
		FROM rhspars_thematic_areas
		WHERE domain_id = $1
		ORDER BY display_order, name
	`, domainID)
	if err != nil {
		return err
	}
	defer rows.Close()

	out := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, order int
		var name string
		if err := rows.Scan(&id, &name, &order); err != nil {
			return err
		}
		out = append(out, map[string]interface{}{
			"id": id, "name": name, "displayOrder": order,
		})
	}
	return c.JSON(out)
}

func GetRHSParsQuestions(c *fiber.Ctx) error {
	thematicAreaID := c.Params("thematicAreaId")
	rows, err := database.DB.Query(`
		SELECT id, question_text, score_weight, display_order
		FROM rhspars_questions
		WHERE thematic_area_id = $1
		ORDER BY display_order, id
	`, thematicAreaID)
	if err != nil {
		return err
	}
	defer rows.Close()

	out := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, weight, order int
		var text string
		if err := rows.Scan(&id, &text, &weight, &order); err != nil {
			return err
		}
		out = append(out, map[string]interface{}{
			"id": id, "text": text, "scoreWeight": weight, "displayOrder": order,
		})
	}
	return c.JSON(out)
}

func GetRHSParsStructure(c *fiber.Ctx) error {
	// Full nested structure for one-shot form load
	domainRows, err := database.DB.Query(`
		SELECT id, code, name, display_order FROM rhspars_domains ORDER BY display_order, name
	`)
	if err != nil {
		return err
	}
	defer domainRows.Close()

	domains := make([]map[string]interface{}, 0)
	for domainRows.Next() {
		var dID, dOrder int
		var code, name string
		if err := domainRows.Scan(&dID, &code, &name, &dOrder); err != nil {
			return err
		}

		taRows, err := database.DB.Query(`
			SELECT id, name, display_order
			FROM rhspars_thematic_areas
			WHERE domain_id = $1
			ORDER BY display_order, name
		`, dID)
		if err != nil {
			return err
		}

		areas := make([]map[string]interface{}, 0)
		for taRows.Next() {
			var taID, taOrder int
			var taName string
			if err := taRows.Scan(&taID, &taName, &taOrder); err != nil {
				taRows.Close()
				return err
			}
			qRows, err := database.DB.Query(`
				SELECT id, question_text, score_weight, display_order
				FROM rhspars_questions
				WHERE thematic_area_id = $1
				ORDER BY display_order, id
			`, taID)
			if err != nil {
				taRows.Close()
				return err
			}
			questions := make([]map[string]interface{}, 0)
			for qRows.Next() {
				var qID, weight, qOrder int
				var text string
				if err := qRows.Scan(&qID, &text, &weight, &qOrder); err != nil {
					qRows.Close()
					taRows.Close()
					return err
				}
				questions = append(questions, map[string]interface{}{
					"id": qID, "text": text, "scoreWeight": weight, "displayOrder": qOrder,
				})
			}
			qRows.Close()
			areas = append(areas, map[string]interface{}{
				"id": taID, "name": taName, "displayOrder": taOrder, "questions": questions,
			})
		}
		taRows.Close()

		domains = append(domains, map[string]interface{}{
			"id": dID, "code": code, "name": name, "displayOrder": dOrder, "thematicAreas": areas,
		})
	}
	return c.JSON(domains)
}

func CreateRHSParsAssessment(c *fiber.Ctx) error {
	var req RHSParsCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "Invalid request body")
	}
	if req.FacilityID == 0 {
		return fiber.NewError(400, "facilityId is required")
	}

	userID, _ := c.Locals("userID").(int)

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var assessmentID int
	err = tx.QueryRow(`
		INSERT INTO rhspars_assessments (
			facility_id, supervision_type, supervision_date, next_supervision_date,
			rh_supplier, facility_type_detail, facility_ownership,
			primary_contact_name, primary_contact_phone, primary_contact_email,
			completed_by_name, completed_by_designation, completed_by_institution,
			assessment_completion_date, notes, created_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING id
	`, req.FacilityID, nullIfEmpty(req.SupervisionType), nullDate(req.SupervisionDate), nullDate(req.NextSupervisionDate),
		nullIfEmpty(req.RHSupplier), nullIfEmpty(req.FacilityTypeDetail), nullIfEmpty(req.FacilityOwnership),
		nullIfEmpty(req.PrimaryContactName), nullIfEmpty(req.PrimaryContactPhone), nullIfEmpty(req.PrimaryContactEmail),
		nullIfEmpty(req.CompletedByName), nullIfEmpty(req.CompletedByDesignation), nullIfEmpty(req.CompletedByInstitution),
		nullDate(req.AssessmentCompletionDate), nullIfEmpty(req.Notes), nullInt(userID),
	).Scan(&assessmentID)
	if err != nil {
		return err
	}

	for _, a := range req.Attendees {
		if strings.TrimSpace(a.FullName) == "" {
			continue
		}
		_, err = tx.Exec(`
			INSERT INTO rhspars_attendees (assessment_id, person_type, full_name, cadre_position, affiliation, contact)
			VALUES ($1,$2,$3,$4,$5,$6)
		`, assessmentID, a.PersonType, a.FullName, a.CadrePosition, a.Affiliation, a.Contact)
		if err != nil {
			return err
		}
	}

	for i, ap := range req.ActionPlan {
		if strings.TrimSpace(ap.GapIdentified) == "" && strings.TrimSpace(ap.ActionRecommendation) == "" {
			continue
		}
		_, err = tx.Exec(`
			INSERT INTO rhspars_action_plan_items
			(assessment_id, gap_identified, action_recommendation, means_of_verification, responsible_person, time_frame, display_order)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
		`, assessmentID, ap.GapIdentified, ap.ActionRecommendation, ap.MeansOfVerification, ap.ResponsiblePerson, ap.TimeFrame, i+1)
		if err != nil {
			return err
		}
	}

	// Load all questions with domain/thematic mapping
	qRows, err := tx.Query(`
		SELECT q.id, q.thematic_area_id, q.score_weight, ta.domain_id, d.code
		FROM rhspars_questions q
		JOIN rhspars_thematic_areas ta ON q.thematic_area_id = ta.id
		JOIN rhspars_domains d ON ta.domain_id = d.id
	`)
	if err != nil {
		return err
	}

	type qMeta struct {
		thematicID int
		weight     int
		domainID   int
		domainCode string
	}
	questions := map[int]qMeta{}
	for qRows.Next() {
		var id, taID, weight, domainID int
		var domainCode string
		if err := qRows.Scan(&id, &taID, &weight, &domainID, &domainCode); err != nil {
			qRows.Close()
			return err
		}
		questions[id] = qMeta{taID, weight, domainID, domainCode}
	}
	qRows.Close()

	sectionPossible := map[int]int{}
	sectionAchieved := map[int]int{}
	sectionDomain := map[int]int{}
	sectionDomainCode := map[int]string{}
	totalPossible := 0
	totalAchieved := 0

	for qIDStr, response := range req.Responses {
		qID, err := strconv.Atoi(qIDStr)
		if err != nil {
			continue
		}
		meta, ok := questions[qID]
		if !ok {
			continue
		}
		if response != "Yes" && response != "No" && response != "NA" {
			continue
		}

		points, counts := ScoreQuestionResponse(response, meta.weight)
		if counts {
			sectionPossible[meta.thematicID] += meta.weight
			totalPossible += meta.weight
			if response == "Yes" {
				sectionAchieved[meta.thematicID] += meta.weight
				totalAchieved += meta.weight
			}
		}
		sectionDomain[meta.thematicID] = meta.domainID
		sectionDomainCode[meta.thematicID] = meta.domainCode

		_, err = tx.Exec(`
			INSERT INTO rhspars_responses (assessment_id, question_id, response, points_earned)
			VALUES ($1,$2,$3,$4)
		`, assessmentID, qID, response, points)
		if err != nil {
			return err
		}
	}

	domainSectionPcts := map[int][]float64{}
	domainCodes := map[int]string{}

	for taID, possible := range sectionPossible {
		achieved := sectionAchieved[taID]
		pct := PercentageScore(achieved, possible)
		_, err = tx.Exec(`
			INSERT INTO rhspars_section_scores
			(assessment_id, thematic_area_id, possible_score, achieved_score, percentage_score)
			VALUES ($1,$2,$3,$4,$5)
		`, assessmentID, taID, possible, achieved, pct)
		if err != nil {
			return err
		}
		dID := sectionDomain[taID]
		domainSectionPcts[dID] = append(domainSectionPcts[dID], pct)
		domainCodes[dID] = sectionDomainCode[taID]
	}

	domainPct := map[string]float64{}
	for dID, pcts := range domainSectionPcts {
		sum := 0.0
		for _, p := range pcts {
			sum += p
		}
		avg := 0.0
		if len(pcts) > 0 {
			avg = sum / float64(len(pcts))
		}
		_, err = tx.Exec(`
			INSERT INTO rhspars_domain_scores (assessment_id, domain_id, percentage_score)
			VALUES ($1,$2,$3)
		`, assessmentID, dID, avg)
		if err != nil {
			return err
		}
		domainPct[domainCodes[dID]] = avg
	}

	mat := domainPct["maternity"]
	fp := domainPct["fp"]
	anc := domainPct["anc"]
	com := domainPct["commodities"]
	hi := domainPct["health_information"]
	gf := domainPct["general_facility"]
	grand := (mat + fp + anc + com + hi + gf) / 6.0
	performanceLevel := PerformanceLevelFromPercentage(grand)

	_, err = tx.Exec(`
		UPDATE rhspars_assessments SET
			maternity_pct=$1, fp_pct=$2, anc_pct=$3, commodities_pct=$4,
			health_info_pct=$5, general_facility_pct=$6, grand_percentage=$7,
			total_possible_score=$8, achieved_score=$9, performance_level=$10,
			updated_at=CURRENT_TIMESTAMP
		WHERE id=$11
	`, mat, fp, anc, com, hi, gf, grand, totalPossible, totalAchieved, performanceLevel, assessmentID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id":               assessmentID,
		"maternityPct":     mat,
		"fpPct":            fp,
		"ancPct":           anc,
		"commoditiesPct":   com,
		"healthInfoPct":    hi,
		"generalFacilityPct": gf,
		"grandPercentage":  grand,
		"totalPossible":    totalPossible,
		"achieved":         totalAchieved,
		"performanceLevel": performanceLevel,
	})
}

func nullIfEmpty(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func nullInt(v int) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

func GetRHSParsAssessments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	query := `
		SELECT a.id, a.supervision_date, a.grand_percentage,
		       a.maternity_pct, a.fp_pct, a.anc_pct, a.commodities_pct,
		       a.health_info_pct, a.general_facility_pct,
		       a.total_possible_score, a.achieved_score, a.performance_level,
		       f.name as facility_name, a.created_at,
		       r.id as region_id, d.id as district_id
		FROM rhspars_assessments a
		JOIN facilities f ON a.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if !isAdmin {
		regionIDs, districtIDs, subcountyIDs, facilityIDs, err := GetUserAdminAreas(userID)
		if err != nil {
			return err
		}
		if len(regionIDs) == 0 && len(districtIDs) == 0 && len(subcountyIDs) == 0 && len(facilityIDs) == 0 {
			return c.JSON([]map[string]interface{}{})
		}
		conds := []string{}
		if len(regionIDs) > 0 {
			ph := makePlaceholders(argIdx, len(regionIDs))
			for _, id := range regionIDs {
				args = append(args, id)
				argIdx++
			}
			conds = append(conds, "r.id IN ("+ph+")")
		}
		if len(districtIDs) > 0 {
			ph := makePlaceholders(argIdx, len(districtIDs))
			for _, id := range districtIDs {
				args = append(args, id)
				argIdx++
			}
			conds = append(conds, "d.id IN ("+ph+")")
		}
		if len(subcountyIDs) > 0 {
			ph := makePlaceholders(argIdx, len(subcountyIDs))
			for _, id := range subcountyIDs {
				args = append(args, id)
				argIdx++
			}
			conds = append(conds, "s.id IN ("+ph+")")
		}
		if len(facilityIDs) > 0 {
			ph := makePlaceholders(argIdx, len(facilityIDs))
			for _, id := range facilityIDs {
				args = append(args, id)
				argIdx++
			}
			conds = append(conds, "f.id IN ("+ph+")")
		}
		query += " AND (" + strings.Join(conds, " OR ") + ")"
	}

	if facilityID := c.Query("facilityId"); facilityID != "" {
		query += fmt.Sprintf(" AND f.id = $%d", argIdx)
		args = append(args, facilityID)
		argIdx++
	}
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

	query += " ORDER BY a.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	out := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var supervisionDate, createdAt sql.NullString
		var grand, mat, fp, anc, com, hi, gf float64
		var totalPossible, achieved int
		var performanceLevel string
		var facilityName string
		var regionID, districtID sql.NullInt64
		if err := rows.Scan(&id, &supervisionDate, &grand, &mat, &fp, &anc, &com, &hi, &gf,
			&totalPossible, &achieved, &performanceLevel,
			&facilityName, &createdAt, &regionID, &districtID); err != nil {
			return err
		}
		out = append(out, map[string]interface{}{
			"id": id,
			"supervisionDate": supervisionDate.String,
			"grandPercentage": grand,
			"maternityPct": mat,
			"fpPct": fp,
			"ancPct": anc,
			"commoditiesPct": com,
			"healthInfoPct": hi,
			"generalFacilityPct": gf,
			"totalPossible": totalPossible,
			"achieved": achieved,
			"performanceLevel": performanceLevel,
			"facilityName": facilityName,
			"createdAt": createdAt.String,
		})
	}
	return c.JSON(out)
}

func makePlaceholders(start, n int) string {
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = fmt.Sprintf("$%d", start+i)
	}
	return strings.Join(parts, ",")
}

func GetRHSParsAssessment(c *fiber.Ctx) error {
	id := c.Params("id")

	var assessment struct {
		ID                 int
		FacilityID         int
		FacilityName       string
		SupervisionType    sql.NullString
		SupervisionDate    sql.NullString
		NextSupervisionDate sql.NullString
		RHSupplier         sql.NullString
		FacilityTypeDetail sql.NullString
		FacilityOwnership  sql.NullString
		PrimaryContactName sql.NullString
		PrimaryContactPhone sql.NullString
		PrimaryContactEmail sql.NullString
		CompletedByName    sql.NullString
		CompletedByDesignation sql.NullString
		CompletedByInstitution sql.NullString
		AssessmentCompletionDate sql.NullString
		Notes              sql.NullString
		MaternityPct       float64
		FPPct              float64
		ANCPct             float64
		CommoditiesPct     float64
		HealthInfoPct      float64
		GeneralFacilityPct float64
		GrandPercentage    float64
		TotalPossible      int
		Achieved           int
		PerformanceLevel   string
		CreatedAt          string
	}

	err := database.DB.QueryRow(`
		SELECT a.id, a.facility_id, f.name, a.supervision_type, a.supervision_date, a.next_supervision_date,
		       a.rh_supplier, a.facility_type_detail, a.facility_ownership,
		       a.primary_contact_name, a.primary_contact_phone, a.primary_contact_email,
		       a.completed_by_name, a.completed_by_designation, a.completed_by_institution,
		       a.assessment_completion_date, a.notes,
		       a.maternity_pct, a.fp_pct, a.anc_pct, a.commodities_pct,
		       a.health_info_pct, a.general_facility_pct, a.grand_percentage,
		       a.total_possible_score, a.achieved_score, a.performance_level, a.created_at
		FROM rhspars_assessments a
		JOIN facilities f ON a.facility_id = f.id
		WHERE a.id = $1
	`, id).Scan(
		&assessment.ID, &assessment.FacilityID, &assessment.FacilityName,
		&assessment.SupervisionType, &assessment.SupervisionDate, &assessment.NextSupervisionDate,
		&assessment.RHSupplier, &assessment.FacilityTypeDetail, &assessment.FacilityOwnership,
		&assessment.PrimaryContactName, &assessment.PrimaryContactPhone, &assessment.PrimaryContactEmail,
		&assessment.CompletedByName, &assessment.CompletedByDesignation, &assessment.CompletedByInstitution,
		&assessment.AssessmentCompletionDate, &assessment.Notes,
		&assessment.MaternityPct, &assessment.FPPct, &assessment.ANCPct, &assessment.CommoditiesPct,
		&assessment.HealthInfoPct, &assessment.GeneralFacilityPct, &assessment.GrandPercentage,
		&assessment.TotalPossible, &assessment.Achieved, &assessment.PerformanceLevel, &assessment.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(404, "RH SPARS assessment not found")
		}
		return err
	}

	domainRows, err := database.DB.Query(`
		SELECT d.id, d.code, d.name, COALESCE(ds.percentage_score, 0)
		FROM rhspars_domains d
		LEFT JOIN rhspars_domain_scores ds ON ds.domain_id = d.id AND ds.assessment_id = $1
		ORDER BY d.display_order
	`, id)
	if err != nil {
		return err
	}
	defer domainRows.Close()
	domainScores := make([]map[string]interface{}, 0)
	for domainRows.Next() {
		var dID int
		var code, name string
		var pct float64
		if err := domainRows.Scan(&dID, &code, &name, &pct); err != nil {
			return err
		}
		domainScores = append(domainScores, map[string]interface{}{
			"id": dID, "code": code, "name": name, "percentage": pct,
			"performanceLevel": PerformanceLevelFromPercentage(pct),
		})
	}

	sectionRows, err := database.DB.Query(`
		SELECT ta.id, ta.name, d.code, ss.possible_score, ss.achieved_score, ss.percentage_score
		FROM rhspars_section_scores ss
		JOIN rhspars_thematic_areas ta ON ss.thematic_area_id = ta.id
		JOIN rhspars_domains d ON ta.domain_id = d.id
		WHERE ss.assessment_id = $1
		ORDER BY d.display_order, ta.display_order
	`, id)
	if err != nil {
		return err
	}
	defer sectionRows.Close()
	sectionScores := make([]map[string]interface{}, 0)
	for sectionRows.Next() {
		var taID, possible, achieved int
		var taName, domainCode string
		var pct float64
		if err := sectionRows.Scan(&taID, &taName, &domainCode, &possible, &achieved, &pct); err != nil {
			return err
		}
		sectionScores = append(sectionScores, map[string]interface{}{
			"thematicAreaId": taID, "name": taName, "domainCode": domainCode,
			"possible": possible, "achieved": achieved, "percentage": pct,
			"performanceLevel": PerformanceLevelFromPercentage(pct),
		})
	}

	respRows, err := database.DB.Query(`
		SELECT q.id, q.question_text, q.score_weight, ar.response, ar.points_earned,
		       ta.id, ta.name, d.code
		FROM rhspars_responses ar
		JOIN rhspars_questions q ON ar.question_id = q.id
		JOIN rhspars_thematic_areas ta ON q.thematic_area_id = ta.id
		JOIN rhspars_domains d ON ta.domain_id = d.id
		WHERE ar.assessment_id = $1
		ORDER BY d.display_order, ta.display_order, q.display_order
	`, id)
	if err != nil {
		return err
	}
	defer respRows.Close()
	responses := make([]map[string]interface{}, 0)
	for respRows.Next() {
		var qID, weight, points, taID int
		var text, response, taName, domainCode string
		if err := respRows.Scan(&qID, &text, &weight, &response, &points, &taID, &taName, &domainCode); err != nil {
			return err
		}
		responses = append(responses, map[string]interface{}{
			"questionId": qID, "questionText": text, "scoreWeight": weight,
			"response": response, "pointsEarned": points,
			"thematicAreaId": taID, "thematicAreaName": taName, "domainCode": domainCode,
		})
	}

	apRows, err := database.DB.Query(`
		SELECT gap_identified, action_recommendation, means_of_verification, responsible_person, time_frame
		FROM rhspars_action_plan_items WHERE assessment_id = $1 ORDER BY display_order, id
	`, id)
	if err != nil {
		return err
	}
	defer apRows.Close()
	actionPlan := make([]map[string]interface{}, 0)
	for apRows.Next() {
		var gap, action, mov, person, tf sql.NullString
		if err := apRows.Scan(&gap, &action, &mov, &person, &tf); err != nil {
			return err
		}
		actionPlan = append(actionPlan, map[string]interface{}{
			"gapIdentified": gap.String, "actionRecommendation": action.String,
			"meansOfVerification": mov.String, "responsiblePerson": person.String, "timeFrame": tf.String,
		})
	}

	return c.JSON(fiber.Map{
		"assessment": map[string]interface{}{
			"id": assessment.ID,
			"facilityId": assessment.FacilityID,
			"facilityName": assessment.FacilityName,
			"supervisionType": assessment.SupervisionType.String,
			"supervisionDate": assessment.SupervisionDate.String,
			"nextSupervisionDate": assessment.NextSupervisionDate.String,
			"rhSupplier": assessment.RHSupplier.String,
			"facilityTypeDetail": assessment.FacilityTypeDetail.String,
			"facilityOwnership": assessment.FacilityOwnership.String,
			"primaryContactName": assessment.PrimaryContactName.String,
			"primaryContactPhone": assessment.PrimaryContactPhone.String,
			"primaryContactEmail": assessment.PrimaryContactEmail.String,
			"completedByName": assessment.CompletedByName.String,
			"completedByDesignation": assessment.CompletedByDesignation.String,
			"completedByInstitution": assessment.CompletedByInstitution.String,
			"assessmentCompletionDate": assessment.AssessmentCompletionDate.String,
			"notes": assessment.Notes.String,
			"maternityPct": assessment.MaternityPct,
			"fpPct": assessment.FPPct,
			"ancPct": assessment.ANCPct,
			"commoditiesPct": assessment.CommoditiesPct,
			"healthInfoPct": assessment.HealthInfoPct,
			"generalFacilityPct": assessment.GeneralFacilityPct,
			"grandPercentage": assessment.GrandPercentage,
			"totalPossible": assessment.TotalPossible,
			"achieved": assessment.Achieved,
			"performanceLevel": assessment.PerformanceLevel,
			"createdAt": assessment.CreatedAt,
		},
		"domainScores": domainScores,
		"sectionScores": sectionScores,
		"responses": responses,
		"actionPlan": actionPlan,
	})
}
