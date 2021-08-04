package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"story-app-monolith/handlers"
	"story-app-monolith/middleware"
	"story-app-monolith/repo"
	"story-app-monolith/services"
)

func SetupRoutes(app *fiber.App) {
	uh := handlers.UserHandler{UserService: services.NewUserService(repo.NewUserRepoImpl())}
	ah := handlers.AuthHandler{AuthService: services.NewAuthService(repo.NewAuthRepoImpl())}
	ch := handlers.CommentHandler{CommentService: services.NewCommentService(repo.NewCommentRepoImpl())}
	sh := handlers.StoryHandler{StoryService: services.NewStoryService(repo.NewStoryRepoImpl())}
	rh := handlers.ReadLaterHandler{ReadLaterService: services.NewReadLaterService(repo.NewReadLaterRepoImpl())}
	reh := handlers.ReplyHandler{ReplyService: services.NewReplyService(repo.NewReplyRepoImpl())}

	app.Use(recover.New())
	api := app.Group("", logger.New())

	auth := api.Group("/auth")
	auth.Post("/login", ah.Login)
	auth.Post("/logout", ah.Logout)
	auth.Post("/reset", ah.ResetPasswordQuery)
	auth.Put("/reset/:token", ah.ResetPassword)
	auth.Get("/account/:code", ah.VerifyCode)

	user := api.Group("/users")
	user.Get("/", middleware.IsLoggedIn, uh.GetAllUsers)
	user.Get("/blocked", middleware.IsLoggedIn, uh.GetAllBlockedUsers)
	user.Post("flag/:username", middleware.IsLoggedIn, uh.UpdateFlagCount)
	user.Post("/", uh.CreateUser)
	user.Put("/profile-visibility", middleware.IsLoggedIn, uh.UpdateProfileVisibility)
	user.Put("/follower-count", middleware.IsLoggedIn, uh.UpdateDisplayFollowerCount)
	user.Put("/message-acceptance", middleware.IsLoggedIn, uh.UpdateMessageAcceptance)
	user.Put("/current-badge", middleware.IsLoggedIn, uh.UpdateCurrentBadge)
	user.Put("/profile-photo", middleware.IsLoggedIn, uh.UpdateProfilePicture)
	user.Put("/background-photo", middleware.IsLoggedIn, uh.UpdateProfileBackgroundPicture)
	user.Put("/current-tagline", middleware.IsLoggedIn, uh.UpdateCurrentTagline)
	user.Put("/block/:username", middleware.IsLoggedIn, uh.BlockUser)
	user.Put("/unblock/:username", middleware.IsLoggedIn, uh.UnblockUser)
	user.Put("/follow/:username", middleware.IsLoggedIn, uh.FollowUser)
	user.Put("/unfollow/:username", middleware.IsLoggedIn, uh.UnfollowUser)
	user.Delete("/delete", middleware.IsLoggedIn, uh.DeleteByID)

	profile := api.Group("/profile")
	profile.Get("/:username", middleware.IsLoggedIn, uh.GetUserProfile)
	profile.Get("/", uh.GetCurrentUserProfile)

	stories := api.Group("/stories")
	stories.Post("/", sh.CreateStory)
	stories.Put("/:id", sh.UpdateStory)
	stories.Put("/like/:id", sh.LikeStory)
	stories.Put("/dislike/:id", sh.DisLikeStory)
	stories.Put("/flag/:id", sh.UpdateFlagCount)
	stories.Get("/featured", sh.FeaturedStories)
	stories.Get("/:id", sh.FindStory)
	stories.Delete("/:id", sh.DeleteStory)
	stories.Get("/", middleware.IsLoggedIn, sh.FindAll)

	comments := api.Group("/comment")
	comments.Post("/:id", ch.CreateCommentOnStory)
	comments.Put("/like/:id", ch.LikeComment)
	comments.Put("/dislike/:id", ch.DisLikeComment)
	comments.Put("/flag/:id", ch.UpdateFlagCount)
	comments.Put("/:id", ch.UpdateById)
	comments.Delete("/:id", ch.DeleteById)

	reply := api.Group("/reply")
	reply.Post("/:id", reh.CreateReply)
	reply.Put("/like/:id", reh.LikeReply)
	reply.Put("/dislike/:id", reh.DisLikeReply)
	reply.Put("/flag/:id", reh.UpdateFlagCount)
	reply.Put("/:id", reh.UpdateById)
	reply.Delete("/:id", reh.DeleteById)

	readLater := api.Group("/read")
	readLater.Post("/:id", rh.Create)
	readLater.Get("/", rh.GetByUsername)
	readLater.Delete("/:id", rh.Delete)
}

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	SetupRoutes(app)

	return app
}
