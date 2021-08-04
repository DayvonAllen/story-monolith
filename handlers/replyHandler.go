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

type ReplyHandler struct {
	ReplyService services.ReplyService
}

func (rh *ReplyHandler) CreateReply(c *fiber.Ctx) error {
	c.Accepts("application/json")
	currentUsername := c.Locals("username").(string)

	replyDto := new(domain.CreateReply)

	err := c.BodyParser(replyDto)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	reply := new(domain.Reply)

	reply.Likes = make([]string, 0, 0)
	reply.Dislikes = make([]string, 0, 0)
	reply.Id = primitive.NewObjectID()
	reply.AuthorUsername = currentUsername
	reply.Content = replyDto.Content
	reply.ResourceId = id
	reply.CreatedAt = time.Now()
	reply.UpdatedAt = time.Now()
	reply.CreatedDate = reply.CreatedAt.Format("January 2, 2006 at 3:04pm")
	reply.UpdatedDate = reply.UpdatedAt.Format("January 2, 2006 at 3:04pm")

	err = rh.ReplyService.Create(reply)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (rh *ReplyHandler) UpdateById(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUsername := c.Locals("username").(string)

	reply := new(domain.Reply)

	err := c.BodyParser(reply)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	reply.UpdatedAt = time.Now()
	reply.Edited = true
	reply.UpdatedDate = reply.UpdatedAt.Format("January 2, 2006 at 3:04pm")

	err = rh.ReplyService.UpdateById(id, reply.Content, reply.Edited, reply.UpdatedAt, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (rh *ReplyHandler) LikeReply(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = rh.ReplyService.LikeReplyById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (rh *ReplyHandler) DisLikeReply(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = rh.ReplyService.DisLikeReplyById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (rh *ReplyHandler) UpdateFlagCount(c *fiber.Ctx) error {
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

	err = rh.ReplyService.UpdateFlagCount(flag)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (rh *ReplyHandler) DeleteById(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = rh.ReplyService.DeleteById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
