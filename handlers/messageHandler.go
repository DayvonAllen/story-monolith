package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"story-app-monolith/domain"
	"story-app-monolith/services"
	"time"
)

type MessageHandler struct {
	MessageService services.MessageService
}

func (mh *MessageHandler) CreateMessage(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUsername := c.Locals("username").(string)

	message := new(domain.Message)

	message.From = currentUsername

	err := c.BodyParser(message)

	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	conversation, err := mh.MessageService.Create(message)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": conversation})
}

func (mh *MessageHandler) DeleteByID(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUsername := c.Locals("username").(string)

	message := new(domain.DeleteMessage)
	err := c.BodyParser(message)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = mh.MessageService.DeleteByID(currentUsername, message.Id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (mh *MessageHandler) DeleteByIDs(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUsername := c.Locals("username").(string)

	messages := new(domain.DeleteMessages)
	err := c.BodyParser(messages)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = mh.MessageService.DeleteAllByIDs(currentUsername, messages.MessageIds)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}