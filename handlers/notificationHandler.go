package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"story-app-monolith/services"
)

type NotificationHandler struct {
	NotificationService services.NotificationService
}

func (nh *NotificationHandler) GetAllUnreadNotificationByUsername(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	notifications, err := nh.NotificationService.GetAllUnreadNotificationByUsername(currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": notifications})
}