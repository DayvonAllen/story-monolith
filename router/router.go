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
}

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	SetupRoutes(app)

	return app
}
