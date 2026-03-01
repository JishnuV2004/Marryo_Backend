package main

import (
	// "context"

	"log"
	"os"

	// "os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"

	// "github.com/joho/godotenv"

	// "github.com/redis/go-redis/v9"
	redisStorage "github.com/gofiber/storage/redis/v3"

	// "github.com/joho/godotenv"

	// "github.com/joho/godotenv"

	"marryo/Cmd/seed"
	config "marryo/Config"
	admin "marryo/Internal/Admin"
	controller "marryo/Internal/Controllers"
	repositories "marryo/Internal/Repositories"
	routes "marryo/Internal/Routes"
	services "marryo/Internal/Services"
	websocket "marryo/Internal/WebSocket"
	"time"
	// "github.com/gin-contrib/cors"
)

func main() {

	// 	err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("❌ .env file not found in root")
	// }

	if err := config.InitRedis(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	// log.Println("Loaded project ID:", os.Getenv("FIREBASE_PROJECT_ID"))

	config.InitDB()

	// HTML Engine
	engine := html.New("Web/Templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./Web/Static")

	storage := redisStorage.New(redisStorage.Config{
    URL: os.Getenv("REDIS_URL"),
})

	// Session Store
	store := session.New(session.Config{
		Storage:        storage,
		CookieHTTPOnly: true,
		CookieSecure:   false, // true in production (HTTPS)
		Expiration:     24 * time.Hour,
	})

	hub := websocket.NewHub()
	go hub.Run()

	// firebaseApp, err := config.InitFirebase()
	// if err != nil {
	// 	log.Fatal("Firebase init failed:", err)
	// }

	// fcmClient, err := firebaseApp.Messaging(context.Background())
	// if err != nil {
	// 	log.Fatal("FCM client failed:", err)
	// }
	// redis := config.Redis()
	// authservice:=services.NewAuthService(repo, redis)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Authorization",
		AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()),
	}))

	//User side
	repo := repositories.NewRepo(config.DB)

	// notificationservice := services.NewNotificationService(repo, fcmClient)

	authservice := services.NewAuthService(repo, config.Redis)
	userService := services.NewUserService(repo)
	chatService := services.NewChatService(repo)
	rbacService := services.NewRoleService(repo)

	authcontroller := controller.NewAuthController(authservice)
	userController := controller.NewUserController(userService)
	chatController := controller.NewChatController(hub, chatService)
	rbacController := controller.NewRoleHandler(*rbacService)

	//Routes
	routes.Routes(app, authcontroller)
	routes.UserRoutes(app, userController)
	routes.ChatRoute(app, chatController)
	routes.RBACRoutes(app, rbacController)

	//Admin
	adminservice := services.NewAdminService(repo)

	admincontroller := controller.NewAdminController(adminservice)

	routes.AdminRoutes(app, admincontroller, authcontroller, userController)

	//Admin Side (SSR)
	adminController := admin.NewAdminController(&repo)
	admin.AdminRoutes(app, store, adminController)

	//seed
	seed.SeedPermissions(config.DB)
	seed.SeedRoles(config.DB)
	seed.SeedAdminRole(config.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
