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
	"github.com/xuri/excelize/v2"
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
		       hw.full_name as health_worker_name,
		       f.name as facility_name, at.name as assessment_type,
		       r.name as region_name, d.name as district_name,
		       a.assessor_name, a.client_name
        FROM assessments a
        JOIN health_workers hw ON a.health_worker_id = hw.id
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

	fmt.Printf("Executing PDF query with %d args\n", len(args))
	fmt.Printf("Query: %s\n", query)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		return err
	}
	defer rows.Close()

	// Collect data and calculate statistics
	type AssessmentRow struct {
		ID               int
		CreatedAt        string
		Percentage       float64
		PerformanceLevel string
		HealthWorkerName string
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
			&row.HealthWorkerName, &row.FacilityName, &row.AssessmentType, &regionName, &districtName,
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

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		fmt.Printf("Error iterating rows: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("Error reading assessment data: %v", err),
		})
	}

	// Calculate statistics
	totalAssessments := len(assessments)
	avgScore := 0.0
	if totalAssessments > 0 {
		avgScore = totalScore / float64(totalAssessments)
	}

	fmt.Printf("Found %d assessments for PDF export\n", totalAssessments)

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
	pdf.Cell(20, 6, "Date")
	pdf.SetXY(32, yPos+2)
	pdf.Cell(35, 6, "Health Worker")
	pdf.SetXY(67, yPos+2)
	pdf.Cell(30, 6, "Facility")
	pdf.SetXY(97, yPos+2)
	pdf.Cell(25, 6, "Region")
	pdf.SetXY(122, yPos+2)
	pdf.Cell(25, 6, "Type")
	pdf.SetXY(147, yPos+2)
	pdf.Cell(18, 6, "Score %")
	pdf.SetXY(165, yPos+2)
	pdf.Cell(20, 6, "Level")
	pdf.SetXY(185, yPos+2)
	pdf.Cell(15, 6, "ID")

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
			pdf.Cell(20, 6, "Date")
			pdf.SetXY(32, 12)
			pdf.Cell(35, 6, "Health Worker")
			pdf.SetXY(67, 12)
			pdf.Cell(30, 6, "Facility")
			pdf.SetXY(97, 12)
			pdf.Cell(25, 6, "Region")
			pdf.SetXY(122, 12)
			pdf.Cell(25, 6, "Type")
			pdf.SetXY(147, 12)
			pdf.Cell(18, 6, "Score %")
			pdf.SetXY(165, 12)
			pdf.Cell(20, 6, "Level")
			pdf.SetXY(185, 12)
			pdf.Cell(15, 6, "ID")
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
		pdf.Cell(20, 5, assessment.CreatedAt)
		pdf.SetXY(32, yPos+1)
		pdf.Cell(35, 5, truncateString(assessment.HealthWorkerName, 20))
		pdf.SetXY(67, yPos+1)
		pdf.Cell(30, 5, truncateString(assessment.FacilityName, 20))
		pdf.SetXY(97, yPos+1)
		pdf.Cell(25, 5, truncateString(assessment.RegionName, 18))
		pdf.SetXY(122, yPos+1)
		pdf.Cell(25, 5, truncateString(assessment.AssessmentType, 18))
		pdf.SetXY(147, yPos+1)
		pdf.CellFormat(18, 5, fmt.Sprintf("%.1f", assessment.Percentage), "0", 0, "C", false, 0, "")

		// Level with color
		pdf.SetTextColor(levelColor[0], levelColor[1], levelColor[2])
		pdf.SetXY(165, yPos+1)
		pdf.CellFormat(20, 5, truncateString(assessment.PerformanceLevel, 15), "0", 0, "C", false, 0, "")

		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(185, yPos+1)
		pdf.CellFormat(15, 5, strconv.Itoa(assessment.ID), "0", 0, "C", false, 0, "")

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

	// Output PDF to buffer
	var buf bytes.Buffer
	fmt.Printf("About to call pdf.Output()...\n")
	err = pdf.Output(&buf)
	if err != nil {
		fmt.Printf("PDF Output() error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to generate PDF: %v", err),
		})
	}

	// Check if PDF is empty
	pdfSize := buf.Len()
	fmt.Printf("PDF buffer size after Output(): %d bytes\n", pdfSize)

	if pdfSize == 0 {
		fmt.Printf("ERROR: PDF buffer is empty after generation. This should not happen.\n")
		fmt.Printf("Attempting to generate a minimal test PDF to verify gofpdf is working...\n")

		// Try to generate a minimal PDF as fallback
		pdf2 := gofpdf.New("P", "mm", "A4", "")
		pdf2.AddPage()
		pdf2.SetFont("Arial", "B", 16)
		pdf2.Cell(40, 10, "Error: PDF generation failed")
		var buf2 bytes.Buffer
		if err2 := pdf2.Output(&buf2); err2 != nil {
			fmt.Printf("Fallback PDF generation also failed: %v\n", err2)
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("PDF generation failed: %v. Fallback also failed: %v", err, err2),
			})
		}
		fallbackSize := buf2.Len()
		fmt.Printf("Fallback PDF generated, size: %d bytes\n", fallbackSize)

		if fallbackSize == 0 {
			fmt.Printf("ERROR: Even fallback PDF is empty. This suggests a problem with gofpdf or the environment.\n")
			return c.Status(500).JSON(fiber.Map{
				"error": "PDF generation is failing. Please check server logs and ensure gofpdf is properly installed.",
			})
		}

		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=assessments-report-%s.pdf", time.Now().Format("20060102-150405")))
		c.Set("Content-Length", fmt.Sprintf("%d", fallbackSize))
		c.Status(200)
		return c.Send(buf2.Bytes())
	}

	fmt.Printf("PDF generated successfully, size: %d bytes\n", pdfSize)

	// Get the bytes from the buffer
	pdfBytes := buf.Bytes()
	if len(pdfBytes) == 0 {
		fmt.Printf("ERROR: PDF bytes are empty even though buffer length is %d\n", pdfSize)
		return c.Status(500).JSON(fiber.Map{
			"error": "PDF bytes are empty. Please check server logs for details.",
		})
	}

	// Verify PDF header (PDF files start with %PDF)
	if len(pdfBytes) < 4 {
		fmt.Printf("ERROR: PDF bytes are too short (%d bytes). First %d bytes: %v\n", len(pdfBytes), len(pdfBytes), pdfBytes)
		return c.Status(500).JSON(fiber.Map{
			"error": "Generated PDF is too short. Please check server logs for details.",
		})
	}

	if string(pdfBytes[0:4]) != "%PDF" {
		fmt.Printf("WARNING: PDF bytes do not start with %%PDF header. First 20 bytes: %v\n", pdfBytes[:min(20, len(pdfBytes))])
		fmt.Printf("This might indicate the PDF was not generated correctly.\n")
	} else {
		fmt.Printf("PDF header verified: starts with %%PDF\n")
	}

	// Set headers first
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=assessments-report-%s.pdf", time.Now().Format("20060102-150405")))
	c.Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	fmt.Printf("Sending PDF response, %d bytes\n", len(pdfBytes))
	fmt.Printf("PDF first 10 bytes: %v\n", pdfBytes[:min(10, len(pdfBytes))])

	// Use Status(200).Send() to ensure both status and body are set together
	// This is the recommended way in Fiber to send binary data with a specific status
	return c.Status(200).Send(pdfBytes)
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ExportAssessmentsXLS generates an Excel file with assessments grouped by health worker
func ExportAssessmentsXLS(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
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

	// Build query to get assessments with health worker info
	query := `
		SELECT a.id, a.created_at, a.percentage_score, a.performance_level,
		       hw.id as health_worker_id, hw.full_name as health_worker_name,
		       hw.email, hw.phone_number,
		       f.name as facility_name, f.id as facility_id,
		       r.name as region_name, d.name as district_name, s.name as subcounty_name,
		       at.name as assessment_type, at.id as assessment_type_id,
		       a.assessor_name, a.client_name
		FROM assessments a
		JOIN health_workers hw ON a.health_worker_id = hw.id
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
	if facilityID := c.Query("facilityId"); facilityID != "" {
		query += fmt.Sprintf(" AND f.id = $%d", argIdx)
		args = append(args, facilityID)
		argIdx++
	}
	if assessmentTypeID := c.Query("assessmentTypeId"); assessmentTypeID != "" {
		query += fmt.Sprintf(" AND a.assessment_type_id = $%d", argIdx)
		args = append(args, assessmentTypeID)
		argIdx++
	}
	if healthWorkerID := c.Query("healthWorkerId"); healthWorkerID != "" {
		query += fmt.Sprintf(" AND hw.id = $%d", argIdx)
		args = append(args, healthWorkerID)
		argIdx++
	}

	query += " ORDER BY hw.full_name, a.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing Excel file: %v\n", err)
		}
	}()

	sheetName := "Assessments"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Set headers
	headers := []string{"Health Worker", "Email", "Phone", "Facility", "Region", "District", "Subcounty",
		"Assessment Type", "Date", "Score %", "Performance Level", "Assessor", "Client"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style headers
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), headerStyle)
	}

	rowNum := 2
	healthWorkerThematicScores := make(map[int]map[int]struct {
		possible, achieved int
		count              int
	})

	for rows.Next() {
		var id, hwID, facilityID, assessmentTypeID int
		var createdAt, performanceLevel, hwName, facilityName, regionName, districtName, subcountyName, assessmentType string
		var percentage float64
		var email, phoneNumber, assessorName, clientName sql.NullString

		if err := rows.Scan(&id, &createdAt, &percentage, &performanceLevel, &hwID, &hwName,
			&email, &phoneNumber, &facilityName, &facilityID, &regionName, &districtName, &subcountyName,
			&assessmentType, &assessmentTypeID, &assessorName, &clientName); err != nil {
			continue
		}

		// Write assessment row
		row := []interface{}{
			hwName,
			func() string {
				if email.Valid {
					return email.String
				}
				return ""
			}(),
			func() string {
				if phoneNumber.Valid {
					return phoneNumber.String
				}
				return ""
			}(),
			facilityName,
			regionName,
			districtName,
			subcountyName,
			assessmentType,
			createdAt,
			percentage,
			performanceLevel,
			func() string {
				if assessorName.Valid {
					return assessorName.String
				}
				return ""
			}(),
			func() string {
				if clientName.Valid {
					return clientName.String
				}
				return ""
			}(),
		}

		for i, val := range row {
			cell := fmt.Sprintf("%c%d", 'A'+i, rowNum)
			f.SetCellValue(sheetName, cell, val)
		}
		rowNum++

		// Get thematic area scores for this assessment
		taRows, err := database.DB.Query(`
			SELECT ta.id, tas.possible_score, tas.achieved_score
			FROM thematic_area_scores tas
			JOIN thematic_areas ta ON tas.thematic_area_id = ta.id
			WHERE tas.assessment_id = $1
		`, id)
		if err == nil {
			if healthWorkerThematicScores[hwID] == nil {
				healthWorkerThematicScores[hwID] = make(map[int]struct {
					possible, achieved int
					count              int
				})
			}
			for taRows.Next() {
				var taID, possible, achieved int
				if err := taRows.Scan(&taID, &possible, &achieved); err == nil {
					ta := healthWorkerThematicScores[hwID][taID]
					ta.possible += possible
					ta.achieved += achieved
					ta.count++
					healthWorkerThematicScores[hwID][taID] = ta
				}
			}
			taRows.Close()
		}
	}

	// Add summary sheet with health worker scores by thematic area
	summarySheet := "Summary by Health Worker"
	summaryIndex, err := f.NewSheet(summarySheet)
	if err != nil {
		return err
	}

	// Get all thematic areas
	taRows, err := database.DB.Query(`
		SELECT DISTINCT ta.id, ta.name, ta.display_order, at.name as assessment_type
		FROM thematic_areas ta
		JOIN assessment_types at ON ta.assessment_type_id = at.id
		ORDER BY at.name, ta.display_order
	`)
	if err != nil {
		return err
	}
	defer taRows.Close()

	thematicAreas := []struct {
		ID             int
		Name           string
		AssessmentType string
	}{}
	for taRows.Next() {
		var ta struct {
			ID             int
			Name           string
			AssessmentType string
		}
		var displayOrder int
		if err := taRows.Scan(&ta.ID, &ta.Name, &displayOrder, &ta.AssessmentType); err == nil {
			thematicAreas = append(thematicAreas, ta)
		}
	}

	// Summary headers
	summaryHeaders := []string{"Health Worker", "Email", "Phone", "Facility", "Region", "District", "Subcounty"}
	for _, ta := range thematicAreas {
		summaryHeaders = append(summaryHeaders, ta.Name+" (%)")
	}
	summaryHeaders = append(summaryHeaders, "Average Score (%)")

	for i, header := range summaryHeaders {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(summarySheet, cell, header)
	}
	if err == nil {
		f.SetCellStyle(summarySheet, "A1", fmt.Sprintf("%c1", 'A'+len(summaryHeaders)-1), headerStyle)
	}

	// Get unique health workers
	hwRows, err := database.DB.Query(`
		SELECT DISTINCT hw.id, hw.full_name, hw.email, hw.phone_number,
		       f.name as facility_name, r.name as region_name, d.name as district_name, s.name as subcounty_name
		FROM health_workers hw
		JOIN facilities f ON hw.facility_id = f.id
		JOIN subcounties s ON f.subcounty_id = s.id
		JOIN districts d ON s.district_id = d.id
		JOIN regions r ON d.region_id = r.id
		WHERE EXISTS (SELECT 1 FROM assessments a WHERE a.health_worker_id = hw.id)
		ORDER BY hw.full_name
	`)
	if err != nil {
		return err
	}
	defer hwRows.Close()

	summaryRowNum := 2
	for hwRows.Next() {
		var hwID int
		var hwName, facilityName, regionName, districtName, subcountyName string
		var email, phoneNumber sql.NullString
		if err := hwRows.Scan(&hwID, &hwName, &email, &phoneNumber, &facilityName, &regionName, &districtName, &subcountyName); err != nil {
			continue
		}

		row := []interface{}{
			hwName,
			func() string {
				if email.Valid {
					return email.String
				}
				return ""
			}(),
			func() string {
				if phoneNumber.Valid {
					return phoneNumber.String
				}
				return ""
			}(),
			facilityName,
			regionName,
			districtName,
			subcountyName,
		}

		// Add thematic area scores
		var totalPercentage float64
		var count int
		for _, ta := range thematicAreas {
			if scores, exists := healthWorkerThematicScores[hwID][ta.ID]; exists && scores.count > 0 {
				percentage := 0.0
				if scores.possible > 0 {
					percentage = (float64(scores.achieved) / float64(scores.possible)) * 100
				}
				row = append(row, percentage)
				totalPercentage += percentage
				count++
			} else {
				row = append(row, "")
			}
		}

		// Add average score (only for assessed thematic areas)
		avgScore := 0.0
		if count > 0 {
			avgScore = totalPercentage / float64(count)
		}
		row = append(row, avgScore)

		for i, val := range row {
			cell := fmt.Sprintf("%c%d", 'A'+i, summaryRowNum)
			f.SetCellValue(summarySheet, cell, val)
		}
		summaryRowNum++
	}

	f.SetActiveSheet(summaryIndex)

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return err
	}

	// Set headers
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=assessments-report-%s.xlsx", time.Now().Format("20060102-150405")))
	c.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	return c.Status(200).Send(buf.Bytes())
}
