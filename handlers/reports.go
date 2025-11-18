package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

// ExportAssessmentsPDF generates a professional PDF report for assessments
func ExportAssessmentsPDF(c *fiber.Ctx) error {
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

	// Build query with filters - similar to GetAssessments
	query := `
		SELECT a.id, a.created_at, a.percentage_score, a.performance_level,
		       f.name as facility_name, at.name as assessment_type,
		       r.name as region_name, d.name as district_name,
		       a.assessor_name, a.client_name
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
	var filterRegionName, filterDistrictName, filterTypeName, filterThematicName string
	if regionID := c.Query("regionId"); regionID != "" {
		query += fmt.Sprintf(" AND r.id = $%d", argIdx)
		args = append(args, regionID)
		// Get region name for display
		var name string
		if err := database.DB.QueryRow("SELECT name FROM regions WHERE id = $1", regionID).Scan(&name); err == nil {
			filterRegionName = name
		}
		argIdx++
	}
	if districtID := c.Query("districtId"); districtID != "" {
		query += fmt.Sprintf(" AND d.id = $%d", argIdx)
		args = append(args, districtID)
		// Get district name for display
		var name string
		if err := database.DB.QueryRow("SELECT name FROM districts WHERE id = $1", districtID).Scan(&name); err == nil {
			filterDistrictName = name
		}
		argIdx++
	}
	if assessmentTypeID := c.Query("assessmentTypeId"); assessmentTypeID != "" {
		query += fmt.Sprintf(" AND a.assessment_type_id = $%d", argIdx)
		args = append(args, assessmentTypeID)
		// Get type name for display
		var name string
		if err := database.DB.QueryRow("SELECT name FROM assessment_types WHERE id = $1", assessmentTypeID).Scan(&name); err == nil {
			filterTypeName = name
		}
		argIdx++
	}
	if thematicAreaID := c.Query("thematicAreaId"); thematicAreaID != "" {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM assessment_responses ar JOIN questions q ON ar.question_id = q.id WHERE ar.assessment_id = a.id AND q.thematic_area_id = $%d)", argIdx)
		args = append(args, thematicAreaID)
		// Get thematic area name for display
		var name string
		if err := database.DB.QueryRow("SELECT name FROM thematic_areas WHERE id = $1", thematicAreaID).Scan(&name); err == nil {
			filterThematicName = name
		}
		argIdx++
	}

	query += " ORDER BY a.created_at DESC LIMIT 500"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Collect data and calculate statistics
	type AssessmentRow struct {
		ID               int
		CreatedAt        string
		Percentage       float64
		PerformanceLevel string
		FacilityName     string
		AssessmentType   string
		RegionName       string
		DistrictName     string
		AssessorName     string
		ClientName       string
	}

	var assessments []AssessmentRow
	var totalScore float64
	levelCounts := map[string]int{
		"Proficient":     0,
		"Competent":      0,
		"Not Acceptable": 0,
	}

	for rows.Next() {
		var row AssessmentRow
		var createdAt, regionName, districtName string
		var assessorNameNull, clientNameNull sql.NullString
		if err := rows.Scan(&row.ID, &createdAt, &row.Percentage, &row.PerformanceLevel,
			&row.FacilityName, &row.AssessmentType, &regionName, &districtName,
			&assessorNameNull, &clientNameNull); err != nil {
			continue
		}

		// Format date safely
		if len(createdAt) >= 19 {
			row.CreatedAt = createdAt[:19]
		} else {
			row.CreatedAt = createdAt
		}
		row.RegionName = regionName
		row.DistrictName = districtName
		if assessorNameNull.Valid {
			row.AssessorName = assessorNameNull.String
		}
		if clientNameNull.Valid {
			row.ClientName = clientNameNull.String
		}

		assessments = append(assessments, row)
		totalScore += row.Percentage
		if count, ok := levelCounts[row.PerformanceLevel]; ok {
			levelCounts[row.PerformanceLevel] = count + 1
		} else {
			levelCounts[row.PerformanceLevel] = 1
		}
	}

	// Calculate statistics
	totalAssessments := len(assessments)
	avgScore := 0.0
	if totalAssessments > 0 {
		avgScore = totalScore / float64(totalAssessments)
	}

	// Create PDF with professional styling
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Family Planning Assessment Report", false)
	pdf.SetAuthor("FP Score Tool", false)
	pdf.SetCreator("FP Score Tool", false)

	// Define colors (RGB values 0-255, converted to 0-1)
	pdf.SetDrawColor(13, 110, 253) // Primary blue
	pdf.SetFillColor(13, 110, 253)
	pdf.SetTextColor(255, 255, 255) // White text

	// Header
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.SetFillColor(13, 110, 253)
	pdf.Rect(10, 10, 190, 25, "F")
	pdf.SetXY(15, 15)
	pdf.Cell(180, 10, "Family Planning Assessment Report")

	// Report metadata
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(15, 40)
	pdf.Cell(90, 6, fmt.Sprintf("Generated: %s", time.Now().Format("02 January 2006 at 15:04")))

	// Filter information
	yPos := 50.0
	if filterRegionName != "" || filterDistrictName != "" || filterTypeName != "" || filterThematicName != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.SetXY(15, yPos)
		pdf.Cell(90, 6, "Filters Applied:")
		yPos += 7
		pdf.SetFont("Arial", "", 9)
		if filterRegionName != "" {
			pdf.SetXY(20, yPos)
			pdf.Cell(90, 5, fmt.Sprintf("Region: %s", filterRegionName))
			yPos += 6
		}
		if filterDistrictName != "" {
			pdf.SetXY(20, yPos)
			pdf.Cell(90, 5, fmt.Sprintf("District: %s", filterDistrictName))
			yPos += 6
		}
		if filterTypeName != "" {
			pdf.SetXY(20, yPos)
			pdf.Cell(90, 5, fmt.Sprintf("Assessment Type: %s", filterTypeName))
			yPos += 6
		}
		if filterThematicName != "" {
			pdf.SetXY(20, yPos)
			pdf.Cell(90, 5, fmt.Sprintf("Thematic Area: %s", filterThematicName))
			yPos += 6
		}
		yPos += 5
	}

	// Summary Statistics Box
	pdf.SetDrawColor(200, 200, 200)
	pdf.SetFillColor(248, 249, 250)
	pdf.Rect(10, yPos, 190, 35, "FD")

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(15, yPos+5)
	pdf.Cell(90, 6, "Summary Statistics")

	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(15, yPos+12)
	pdf.Cell(45, 5, fmt.Sprintf("Total Assessments: %d", totalAssessments))
	pdf.SetXY(15, yPos+18)
	pdf.Cell(45, 5, fmt.Sprintf("Average Score: %.1f%%", avgScore))

	pdf.SetXY(110, yPos+12)
	pdf.Cell(45, 5, fmt.Sprintf("Proficient: %d", levelCounts["Proficient"]))
	pdf.SetXY(110, yPos+18)
	pdf.Cell(45, 5, fmt.Sprintf("Competent: %d", levelCounts["Competent"]))
	pdf.SetXY(110, yPos+24)
	pdf.Cell(45, 5, fmt.Sprintf("Not Acceptable: %d", levelCounts["Not Acceptable"]))

	yPos += 45

	// Table Header
	pdf.SetDrawColor(13, 110, 253)
	pdf.SetFillColor(13, 110, 253)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 9)

	headerHeight := 8.0
	pdf.Rect(10, yPos, 190, headerHeight, "F")
	pdf.SetXY(12, yPos+2)
	pdf.Cell(25, 6, "Date")
	pdf.SetXY(37, yPos+2)
	pdf.Cell(40, 6, "Facility")
	pdf.SetXY(77, yPos+2)
	pdf.Cell(30, 6, "Region")
	pdf.SetXY(107, yPos+2)
	pdf.Cell(30, 6, "Type")
	pdf.SetXY(137, yPos+2)
	pdf.Cell(20, 6, "Score %")
	pdf.SetXY(157, yPos+2)
	pdf.Cell(20, 6, "Level")
	pdf.SetXY(177, yPos+2)
	pdf.Cell(20, 6, "ID")

	yPos += headerHeight

	// Table rows
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 8)
	rowHeight := 6.0
	alternate := false

	// If no assessments, show a message
	if len(assessments) == 0 {
		pdf.SetFont("Arial", "I", 10)
		pdf.SetXY(15, yPos+10)
		pdf.Cell(180, 6, "No assessments found matching the selected criteria.")
		yPos += 20
	}

	for i, assessment := range assessments {
		// Check if we need a new page
		if yPos+rowHeight > 270 {
			pdf.AddPage()
			// Redraw header on new page
			pdf.SetDrawColor(13, 110, 253)
			pdf.SetFillColor(13, 110, 253)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFont("Arial", "B", 9)
			pdf.Rect(10, 10, 190, headerHeight, "F")
			pdf.SetXY(12, 12)
			pdf.Cell(25, 6, "Date")
			pdf.SetXY(37, 12)
			pdf.Cell(40, 6, "Facility")
			pdf.SetXY(77, 12)
			pdf.Cell(30, 6, "Region")
			pdf.SetXY(107, 12)
			pdf.Cell(30, 6, "Type")
			pdf.SetXY(137, 12)
			pdf.Cell(20, 6, "Score %")
			pdf.SetXY(157, 12)
			pdf.Cell(20, 6, "Level")
			pdf.SetXY(177, 12)
			pdf.Cell(20, 6, "ID")
			yPos = 10.0 + headerHeight
			alternate = false
		}

		// Alternate row colors
		if alternate {
			pdf.SetFillColor(248, 249, 250)
			pdf.Rect(10, yPos, 190, rowHeight, "F")
		} else {
			pdf.SetFillColor(255, 255, 255)
			pdf.Rect(10, yPos, 190, rowHeight, "F")
		}
		alternate = !alternate

		// Level color coding
		levelColor := []int{0, 0, 0} // Default black
		if assessment.PerformanceLevel == "Proficient" {
			levelColor = []int{25, 135, 84} // Green
		} else if assessment.PerformanceLevel == "Competent" {
			levelColor = []int{255, 193, 7} // Yellow/Orange
		} else {
			levelColor = []int{220, 53, 69} // Red
		}

		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(12, yPos+1)
		pdf.Cell(25, 5, assessment.CreatedAt)
		pdf.SetXY(37, yPos+1)
		pdf.Cell(40, 5, truncateString(assessment.FacilityName, 25))
		pdf.SetXY(77, yPos+1)
		pdf.Cell(30, 5, truncateString(assessment.RegionName, 20))
		pdf.SetXY(107, yPos+1)
		pdf.Cell(30, 5, truncateString(assessment.AssessmentType, 20))
		pdf.SetXY(137, yPos+1)
		pdf.CellFormat(20, 5, fmt.Sprintf("%.1f", assessment.Percentage), "0", 0, "C", false, 0, "")

		// Level with color
		pdf.SetTextColor(levelColor[0], levelColor[1], levelColor[2])
		pdf.SetXY(157, yPos+1)
		pdf.CellFormat(20, 5, truncateString(assessment.PerformanceLevel, 15), "0", 0, "C", false, 0, "")

		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(177, yPos+1)
		pdf.CellFormat(20, 5, strconv.Itoa(assessment.ID), "0", 0, "C", false, 0, "")

		// Draw border
		pdf.SetDrawColor(200, 200, 200)
		pdf.Line(10, yPos+rowHeight, 200, yPos+rowHeight)

		yPos += rowHeight

		// Limit to prevent huge PDFs
		if i >= 499 {
			break
		}
	}

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.SetXY(10, 285)
	pdf.CellFormat(190, 5, fmt.Sprintf("Page %d - Generated by FP Score Tool", pdf.PageNo()), "0", 0, "C", false, 0, "")

	// Output PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return fiber.NewError(500, fmt.Sprintf("Failed to generate PDF: %v", err))
	}

	// Check if PDF is empty
	if buf.Len() == 0 {
		return fiber.NewError(500, "Generated PDF is empty")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=assessments-report-%s.pdf", time.Now().Format("20060102-150405")))
	c.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	return c.Send(buf.Bytes())
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
