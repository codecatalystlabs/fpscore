package handlers

import (
	"database/sql"
	"fmt"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
)

type moveDistrictsRequest struct {
	DistrictIDs []int `json:"districtIds"`
	NewRegionID int   `json:"newRegionId"`
}

type moveSubcountiesRequest struct {
	SubcountyIDs  []int `json:"subcountyIds"`
	NewDistrictID int   `json:"newDistrictId"`
}

type moveFacilitiesRequest struct {
	FacilityIDs    []int `json:"facilityIds"`
	NewSubcountyID int   `json:"newSubcountyId"`
}

func MoveDistricts(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	var req moveDistrictsRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if len(req.DistrictIDs) == 0 || req.NewRegionID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "districtIds and newRegionId are required")
	}

	if err := ensureRegionExists(req.NewRegionID); err != nil {
		return err
	}
	if !isAdmin && !userHasAccessToArea(userID, req.NewRegionID, 0, 0, 0) {
		return fiber.NewError(fiber.StatusForbidden, "You do not have access to the target region")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, districtID := range req.DistrictIDs {
		var currentRegionID int
		if err := tx.QueryRow(`SELECT region_id FROM districts WHERE id = $1`, districtID).Scan(&currentRegionID); err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("District %d not found", districtID))
			}
			return err
		}

		if !isAdmin && !userHasAccessToArea(userID, currentRegionID, districtID, 0, 0) {
			return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("You do not have access to district %d", districtID))
		}

		if _, err := tx.Exec(`UPDATE districts SET region_id = $1 WHERE id = $2`, req.NewRegionID, districtID); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message":       "Districts moved successfully",
		"districtCount": len(req.DistrictIDs),
	})
}

func MoveSubcounties(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	var req moveSubcountiesRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if len(req.SubcountyIDs) == 0 || req.NewDistrictID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "subcountyIds and newDistrictId are required")
	}

	newRegionID, err := getRegionForDistrict(req.NewDistrictID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusBadRequest, "Target district not found")
		}
		return err
	}
	if !isAdmin && !userHasAccessToArea(userID, newRegionID, req.NewDistrictID, 0, 0) {
		return fiber.NewError(fiber.StatusForbidden, "You do not have access to the target district")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, subcountyID := range req.SubcountyIDs {
		currentRegionID, currentDistrictID, err := getRegionDistrictForSubcounty(subcountyID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Subcounty %d not found", subcountyID))
			}
			return err
		}

		if !isAdmin && !userHasAccessToArea(userID, currentRegionID, currentDistrictID, subcountyID, 0) {
			return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("You do not have access to subcounty %d", subcountyID))
		}

		if _, err := tx.Exec(`UPDATE subcounties SET district_id = $1 WHERE id = $2`, req.NewDistrictID, subcountyID); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message":        "Subcounties moved successfully",
		"subcountyCount": len(req.SubcountyIDs),
	})
}

func MoveFacilities(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	isAdmin := IsAdmin(userID)

	var req moveFacilitiesRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if len(req.FacilityIDs) == 0 || req.NewSubcountyID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "facilityIds and newSubcountyId are required")
	}

	newRegionID, newDistrictID, err := getRegionDistrictForSubcounty(req.NewSubcountyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusBadRequest, "Target subcounty not found")
		}
		return err
	}
	if !isAdmin && !userHasAccessToArea(userID, newRegionID, newDistrictID, req.NewSubcountyID, 0) {
		return fiber.NewError(fiber.StatusForbidden, "You do not have access to the target subcounty")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, facilityID := range req.FacilityIDs {
		currentRegionID, currentDistrictID, currentSubcountyID, err := getAreaForFacility(facilityID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Facility %d not found", facilityID))
			}
			return err
		}

		if !isAdmin && !userHasAccessToArea(userID, currentRegionID, currentDistrictID, currentSubcountyID, facilityID) {
			return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("You do not have access to facility %d", facilityID))
		}

		if _, err := tx.Exec(`UPDATE facilities SET subcounty_id = $1 WHERE id = $2`, req.NewSubcountyID, facilityID); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message":       "Facilities moved successfully",
		"facilityCount": len(req.FacilityIDs),
	})
}

func ensureRegionExists(regionID int) error {
	var name string
	if err := database.DB.QueryRow(`SELECT name FROM regions WHERE id = $1`, regionID).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusBadRequest, "Target region not found")
		}
		return err
	}
	return nil
}
