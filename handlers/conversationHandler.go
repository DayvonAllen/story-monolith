package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"story-app-monolith/services"
)

type ConversationHandler struct {
	ConversationService services.ConversationService
}

func (ch *ConversationHandler) FindConversation(c *fiber.Ctx) error {
	c.Accepts("application/json")
	currentUsername := c.Locals("username").(string)


	username := c.Params("username")

	conversation, err := ch.ConversationService.FindConversation(currentUsername, username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": conversation.Messages})
}

func (ch *ConversationHandler) GetConversationPreviews(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUsername := c.Locals("username").(string)

	conversation, err := ch.ConversationService.GetConversationPreviews(currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": conversation})
}
