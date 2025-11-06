package handlers

import (
	"database/sql"
	"strings"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT token and adds user info to context
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header")
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	// Add user info to context
	c.Locals("userID", int(claims["sub"].(float64)))
	c.Locals("userEmail", claims["email"])
	c.Locals("userName", claims["name"])

	return c.Next()
}

// GetUserAdminAreas returns the administrative areas a user has access to
func GetUserAdminAreas(userID int) ([]int, []int, []int, []int, error) {
	var regionIDs, districtIDs, subcountyIDs, facilityIDs []int

	rows, err := database.DB.Query(`
		SELECT region_id, district_id, subcounty_id, facility_id 
		FROM user_admin_areas 
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var regionID, districtID, subcountyID, facilityID sql.NullInt64
		if err := rows.Scan(&regionID, &districtID, &subcountyID, &facilityID); err != nil {
			continue
		}
		if regionID.Valid {
			regionIDs = append(regionIDs, int(regionID.Int64))
		}
		if districtID.Valid {
			districtIDs = append(districtIDs, int(districtID.Int64))
		}
		if subcountyID.Valid {
			subcountyIDs = append(subcountyIDs, int(subcountyID.Int64))
		}
		if facilityID.Valid {
			facilityIDs = append(facilityIDs, int(facilityID.Int64))
		}
	}

	return regionIDs, districtIDs, subcountyIDs, facilityIDs, nil
}
