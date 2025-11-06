package handlers

import (
	"bytes"
	"fmt"
	"time"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

// ExportAssessmentsPDF generates a simple PDF report for assessments
func ExportAssessmentsPDF(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	_, _, _, _, err := GetUserAdminAreas(userID)
	if err != nil {
		return err
	}

	// Basic query (reuse GetAssessments filter params)
	query := `
        SELECT a.created_at, f.name as facility_name, at.name as assessment_type, a.percentage_score, a.performance_level
        FROM assessments a
        JOIN facilities f ON a.facility_id = f.id
        JOIN assessment_types at ON a.assessment_type_id = at.id
        ORDER BY a.created_at DESC
        LIMIT 200
    `

	rows, err := database.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Assessments Report", false)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Assessments Report")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(190, 6, time.Now().Format("02 Jan 2006 15:04"))
	pdf.Ln(10)

	// Table header
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(40, 8, "Date", "1", 0, "L", false, 0, "")
	pdf.CellFormat(55, 8, "Facility", "1", 0, "L", false, 0, "")
	pdf.CellFormat(45, 8, "Type", "1", 0, "L", false, 0, "")
	pdf.CellFormat(25, 8, "Score %", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Level", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	for rows.Next() {
		var createdAt string
		var facility, atype, level string
		var percentage float64
		if err := rows.Scan(&createdAt, &facility, &atype, &percentage, &level); err != nil {
			continue
		}
		pdf.CellFormat(40, 7, createdAt[:19], "1", 0, "L", false, 0, "")
		pdf.CellFormat(55, 7, facility, "1", 0, "L", false, 0, "")
		pdf.CellFormat(45, 7, atype, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 7, fmt.Sprintf("%.1f", percentage), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 7, level, "1", 1, "C", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return err
	}
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=assessments.pdf")
	return c.SendStream(&buf)
}
