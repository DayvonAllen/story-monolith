package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"story-app-monolith/domain"
	"story-app-monolith/services"
	"story-app-monolith/util"
	"strings"
)

type UserHandler struct {
	UserService services.UserService
}

func (uh *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := c.Query("page", "1")

	currentUsername := c.Locals("username").(string)
	currentUserId := c.Locals("id").(primitive.ObjectID)

	users, err := uh.UserService.GetAllUsers(currentUserId, page, c.Context(), currentUsername)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": users})
}

func (uh *UserHandler) GetAllBlockedUsers(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)
	currentUserId := c.Locals("id").(primitive.ObjectID)

	users, err := uh.UserService.GetAllBlockedUsers(currentUserId, c.Context(), currentUsername)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": users})
}

func (uh *UserHandler) GetCurrentUserProfile(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)

	user, err := uh.UserService.GetCurrentUserProfile(currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (uh *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	username := c.Params("username")
	currentUsername := c.Locals("username").(string)

	user, err := uh.UserService.GetUserProfile(username, currentUsername)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (uh *UserHandler) CreateUser(c *fiber.Ctx) error {
	c.Accepts("application/json")
	createUserDto := new(domain.CreateUserDto)

	err := c.BodyParser(createUserDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	if !util.IsEmail(createUserDto.Email) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid email")})
	}

	if len(createUserDto.Username) <= 1 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid username")})
	}

	user := util.CreateUser(createUserDto)

	user.Following = make([]string,0, 0)
	user.Followers = make([]string,0, 0)
	user.DisplayFollowerCount = true

	err = uh.UserService.CreateUser(user)

	if err != nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateProfileVisibility(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	userDto := new(domain.UpdateProfileVisibility)

	err := c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateProfileVisibility(currentUserId, userDto, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateDisplayFollowerCount(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	userDto := new(domain.UpdateDisplayFollowerCount)

	err := c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateDisplayFollowerCount(currentUserId, userDto)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateMessageAcceptance(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	userDto := new(domain.UpdateMessageAcceptance)

	err := c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateMessageAcceptance(currentUserId, userDto, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}

		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateCurrentBadge(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	userDto := new(domain.UpdateCurrentBadge)

	err := c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateCurrentBadge(currentUserId, userDto, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateProfilePicture(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	userDto := new(domain.UpdateProfilePicture)

	err := c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateProfilePicture(currentUserId, userDto, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateProfileBackgroundPicture(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	userDto := new(domain.UpdateProfileBackgroundPicture)

	err := c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateProfileBackgroundPicture(currentUserId, userDto, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateCurrentTagline(c *fiber.Ctx) error {
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	userDto := new(domain.UpdateCurrentTagline)

	err := c.BodyParser(userDto)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = uh.UserService.UpdateCurrentTagline(currentUserId, userDto, c.Context())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UpdateFlagCount(c *fiber.Ctx) error {
	username := c.Params("username")
	c.Accepts("application/json")

	currentUserId := c.Locals("id").(primitive.ObjectID)

	flag := new(domain.Flag)

	err := c.BodyParser(flag)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	flag.FlaggedUsername = strings.ToLower(username)
	flag.FlaggerID = currentUserId

	err = uh.UserService.UpdateFlagCount(flag)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) DeleteByID(c *fiber.Ctx) error {
	currentUsername := c.Locals("username").(string)
	currentUserId := c.Locals("id").(primitive.ObjectID)

	err := uh.UserService.DeleteByID(currentUserId, c.Context(), currentUsername)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) FollowUser(c *fiber.Ctx) error {
	username := c.Locals("username").(string)

	currentUsername := c.Params("username")

	err := uh.UserService.FollowUser(strings.ToLower(currentUsername), username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UnfollowUser(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	currentUsername := c.Params("username")

	err := uh.UserService.UnfollowUser(strings.ToLower(currentUsername), username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) BlockUser(c *fiber.Ctx) error {
	username := c.Params("username")
	currentUsername := c.Locals("username").(string)
	currentUserId := c.Locals("id").(primitive.ObjectID)

	err := uh.UserService.BlockUser(currentUserId, strings.ToLower(username), c.Context(), currentUsername)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (uh *UserHandler) UnblockUser(c *fiber.Ctx) error {
	username := c.Params("username")
	currentUsername := c.Locals("username").(string)
	currentUserId := c.Locals("id").(primitive.ObjectID)


	err := uh.UserService.UnblockUser(currentUserId, strings.ToLower(username), c.Context(), currentUsername)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
