package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"story-app-monolith/domain"
	"story-app-monolith/services"
	"strconv"
	"time"
)

type StoryHandler struct {
	StoryService services.StoryService
}

func (s *StoryHandler) CreateStory(c *fiber.Ctx) error {
	c.Accepts("application/json")

	storyDto := new(domain.CreateStoryDto)

	currentUsername := c.Locals("username").(string)

	err := c.BodyParser(storyDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	storyDto.AuthorUsername = currentUsername
	storyDto.CreatedAt = time.Now()
	storyDto.UpdatedAt = time.Now()
	storyDto.Likes = make([]string, 0)
	storyDto.Dislikes = make([]string, 0)
	storyDto.CreatedDate = storyDto.CreatedAt.Format("January 2, 2006")
	storyDto.UpdatedDate = storyDto.UpdatedAt.Format("January 2, 2006")

	err = storyDto.Tag.ValidateTag(&storyDto.Tag)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = s.StoryService.Create(storyDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (s *StoryHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newStoriesQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newStoriesQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	stories, err := s.StoryService.FindAll(page, isNew)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": stories})
}

func (s *StoryHandler) FeaturedStories(c *fiber.Ctx) error {
	stories, err := s.StoryService.FeaturedStories()

	if err != nil {
		fmt.Println("err")
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": stories})
}

func (s *StoryHandler) UpdateStory(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUsername := c.Locals("username").(string)

	storyDto := new(domain.UpdateStoryDto)

	err := c.BodyParser(storyDto)

	storyDto.UpdatedAt = time.Now()
	storyDto.Updated = true
	storyDto.UpdatedDate = storyDto.UpdatedAt.Format("January 2, 2006")

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = storyDto.Tag.ValidateTag(&storyDto.Tag)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = s.StoryService.UpdateById(id, storyDto.Content, storyDto.Title, currentUsername, &storyDto.Tag, storyDto.Updated)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (s *StoryHandler) FindStory(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	story, err := s.StoryService.FindById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": &story})
}

func (s *StoryHandler) LikeStory(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = s.StoryService.LikeStoryById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (s *StoryHandler) DisLikeStory(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = s.StoryService.DisLikeStoryById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (s *StoryHandler) UpdateFlagCount(c *fiber.Ctx) error {
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

	err = s.StoryService.UpdateFlagCount(flag)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (s *StoryHandler) DeleteStory(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = s.StoryService.DeleteById(id, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
