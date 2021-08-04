package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"story-app-monolith/services"
)

type ReadLaterHandler struct {
	ReadLaterService services.ReadLaterService
}

func (r *ReadLaterHandler) Create(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = r.ReadLaterService.Create(currentUsername, id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (r *ReadLaterHandler) GetByUsername(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	items, err := r.ReadLaterService.GetByUsername(currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": items})
}

func (r *ReadLaterHandler) Delete(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = r.ReadLaterService.Delete(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
