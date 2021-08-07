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
	//mh := handlers.MessageHandler{MessageService: services.NewMessageService(repo.NewMessageRepoImpl())}
	conh := handlers.ConversationHandler{ConversationService: services.NewConversationService(repo.NewConversationRepoImpl())}

	app.Use(recover.New())
	api := app.Group("", logger.New())

	auth := api.Group("/auth")
	auth.Post("/login", ah.Login)
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
	profile.Get("/", middleware.IsLoggedIn, uh.GetCurrentUserProfile)

	stories := api.Group("/stories")
	stories.Post("/", middleware.IsLoggedIn, sh.CreateStory)
	stories.Put("/:id", middleware.IsLoggedIn, sh.UpdateStory)
	stories.Put("/like/:id", middleware.IsLoggedIn, sh.LikeStory)
	stories.Put("/dislike/:id", middleware.IsLoggedIn, sh.DisLikeStory)
	stories.Put("/flag/:id", middleware.IsLoggedIn, sh.UpdateFlagCount)
	stories.Get("/featured", sh.FeaturedStories)
	stories.Get("/:id", middleware.IsLoggedIn, sh.FindStory)
	stories.Delete("/:id", middleware.IsLoggedIn, sh.DeleteStory)
	stories.Get("/", middleware.IsLoggedIn, sh.FindAll)

	comments := api.Group("/comment")
	comments.Post("/:id", middleware.IsLoggedIn, ch.CreateCommentOnStory)
	comments.Put("/like/:id", middleware.IsLoggedIn, ch.LikeComment)
	comments.Put("/dislike/:id", middleware.IsLoggedIn, ch.DisLikeComment)
	comments.Put("/flag/:id", middleware.IsLoggedIn, ch.UpdateFlagCount)
	comments.Put("/:id", middleware.IsLoggedIn, ch.UpdateById)
	comments.Delete("/:id", middleware.IsLoggedIn, ch.DeleteById)

	reply := api.Group("/reply")
	reply.Post("/:id", middleware.IsLoggedIn, reh.CreateReply)
	reply.Put("/like/:id", middleware.IsLoggedIn, reh.LikeReply)
	reply.Put("/dislike/:id", middleware.IsLoggedIn, reh.DisLikeReply)
	reply.Put("/flag/:id", middleware.IsLoggedIn, reh.UpdateFlagCount)
	reply.Put("/:id", middleware.IsLoggedIn, reh.UpdateById)
	reply.Delete("/:id", middleware.IsLoggedIn, reh.DeleteById)

	readLater := api.Group("/read")
	readLater.Post("/:id", middleware.IsLoggedIn, rh.Create)
	readLater.Get("/", middleware.IsLoggedIn, rh.GetByUsername)
	readLater.Delete("/:id", middleware.IsLoggedIn, rh.Delete)

	//messages := api.Group("/messages")
	//messages.Post("/", middleware.IsLoggedIn, mh.CreateMessage)
	//messages.Delete("/multi", middleware.IsLoggedIn, mh.DeleteByIDs)
	//messages.Delete("/", middleware.IsLoggedIn, mh.DeleteByID)

	conversations := api.Group("/conversation")
	conversations.Get("/:username", middleware.IsLoggedIn, conh.FindConversation)
	conversations.Get("/", middleware.IsLoggedIn, conh.GetConversationPreviews)
}

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	SetupRoutes(app)

	return app
}
