package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"story-app-monolith/domain"
	"story-app-monolith/services"
	"strings"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	c.Accepts("application/json")
	details := new(domain.LoginDetails)
	err := c.BodyParser(details)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	if len(details.Email) <= 1 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid username")})
	}

	if len(details.Password) <= 5 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid password")})
	}

	var auth domain.Authentication

	_, token, err := ah.AuthService.Login(strings.ToLower(details.Email), details.Password)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("Authentication failure")})
	}

	signedToken := make([]byte, 0, 100)
	signedToken = append(signedToken, []byte("Bearer " + token + "|")...)
	t, err := auth.SignToken([]byte(token))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	signedToken = append(signedToken, t...)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": string(signedToken)})
}

func (ah *AuthHandler) ResetPasswordQuery(c *fiber.Ctx) error {
	c.Accepts("application/json")
	q := new(domain.ResetPasswordQuery)
	err := c.BodyParser(q)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = ah.AuthService.ResetPasswordQuery(q.Email)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return nil
}

func (ah *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	c.Accepts("application/json")
	p := new(domain.ResetPassword)
	err := c.BodyParser(p)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	token := c.Params("token")

	err = ah.AuthService.ResetPassword(token, p.Password)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return nil
}

func (ah *AuthHandler) VerifyCode(c *fiber.Ctx) error {
	c.Accepts("application/json")

	code := c.Params("code")

	err := ah.AuthService.VerifyCode(code)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return nil
}