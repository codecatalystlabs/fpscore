package handlers

import (
	"database/sql"
	"strings"
	"time"

	"fpscore/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("change-this-secret")

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	var id int
	var name, dbEmail, passwordHash string
	var isActive bool
	err := database.DB.QueryRow("SELECT id, name, email, password_hash, is_active FROM users WHERE email=$1", email).Scan(&id, &name, &dbEmail, &passwordHash, &isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}
		return err
	}
	if !isActive {
		return fiber.NewError(fiber.StatusUnauthorized, "User inactive")
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":   id,
		"name":  name,
		"email": dbEmail,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"token": signed, "user": fiber.Map{"id": id, "name": name, "email": email}})
}

type BootstrapRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// BootstrapAdmin creates an initial admin user when the system is empty
func BootstrapAdmin(c *fiber.Ctx) error {
	var count int
	if err := database.DB.QueryRow("SELECT COUNT(1) FROM users").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Users already exist")
	}

	var req BootstrapRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Name, email and password are required")
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var userID int
	if err := database.DB.QueryRow(`INSERT INTO users(name,email,password_hash,is_active) VALUES($1,$2,$3,true) RETURNING id`, req.Name, email, string(hash)).Scan(&userID); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"id": userID, "message": "Bootstrap admin created"})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// ChangePassword allows the current user to change their password
func ChangePassword(c *fiber.Ctx) error {
	uid := c.Locals("userID").(int)
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}
	var currentHash string
	if err := database.DB.QueryRow("SELECT password_hash FROM users WHERE id=$1", uid).Scan(&currentHash); err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword)) != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Old password is incorrect")
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if _, err := database.DB.Exec("UPDATE users SET password_hash=$1 WHERE id=$2", string(newHash), uid); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Password changed"})
}
