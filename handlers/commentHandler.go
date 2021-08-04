package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"story-app-monolith/domain"
	"story-app-monolith/services"
	"time"
)

type CommentHandler struct {
	CommentService services.CommentService
}

func (ch *CommentHandler) CreateCommentOnStory(c *fiber.Ctx) error {
	c.Accepts("application/json")
	currentUsername := c.Locals("username").(string)

	comment := new(domain.Comment)

	err := c.BodyParser(comment)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	comment.Likes = make([]string, 0, 0)
	comment.Dislikes = make([]string, 0, 0)
	comment.Id = primitive.NewObjectID()
	comment.AuthorUsername = currentUsername
	comment.ResourceId = id
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	comment.CreatedDate = comment.CreatedAt.Format("January 2, 2006 at 3:04pm")
	comment.UpdatedDate = comment.UpdatedAt.Format("January 2, 2006 at 3:04pm")

	err = ch.CommentService.Create(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ch *CommentHandler) UpdateById(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUsername := c.Locals("username").(string)

	comment := new(domain.Comment)

	err := c.BodyParser(comment)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	comment.UpdatedAt = time.Now()
	comment.Edited = true
	comment.UpdatedDate = comment.UpdatedAt.Format("January 2, 2006 at 3:04pm")

	err = ch.CommentService.UpdateById(id, comment.Content, comment.Edited, comment.UpdatedAt, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ch *CommentHandler) LikeComment(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = ch.CommentService.LikeCommentById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ch *CommentHandler) DisLikeComment(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = ch.CommentService.DisLikeCommentById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ch *CommentHandler) UpdateFlagCount(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	flag := new(domain.Flag)

	err := c.BodyParser(flag)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	flag.FlaggedResource = id
	flag.FlaggerID = currentUserId

	err = ch.CommentService.UpdateFlagCount(flag)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ch *CommentHandler) DeleteById(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = ch.CommentService.DeleteById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
