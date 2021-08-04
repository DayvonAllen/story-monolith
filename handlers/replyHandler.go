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
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	replyDto := new(domain.CreateReply)

	err = c.BodyParser(replyDto)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	reply := new(domain.Reply)

	reply.Likes = make([]string, 0, 0)
	reply.Dislikes = make([]string, 0, 0)
	reply.Id = primitive.NewObjectID()
	reply.AuthorUsername = u.Username
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
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	reply := new(domain.Reply)

	err = c.BodyParser(reply)

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

	err = rh.ReplyService.UpdateById(id, reply.Content, reply.Edited, reply.UpdatedAt, u.Username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (rh *ReplyHandler) LikeReply(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = rh.ReplyService.LikeReplyById(id, u.Username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (rh *ReplyHandler) DisLikeReply(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = rh.ReplyService.DisLikeReplyById(id, u.Username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (rh *ReplyHandler) UpdateFlagCount(c *fiber.Ctx) error {
	c.Accepts("application/json")
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	flag := new(domain.Flag)

	err = c.BodyParser(flag)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	flag.FlaggedResource = id
	flag.FlaggerID = u.Id

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
	token := c.Get("Authorization")

	var auth domain.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": "Unauthorized user"})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = rh.ReplyService.DeleteById(id, u.Username)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
