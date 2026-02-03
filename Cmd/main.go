package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	config "marryo/Config"
	controller "marryo/Internal/Controllers"
	repositories "marryo/Internal/Repositories"
	routes "marryo/Internal/Routes"
	services "marryo/Internal/Services"
	"time"
	// "github.com/gin-contrib/cors"
)


func main(){
	
	if err := config.InitRedis(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	err := godotenv.Load()
	if err != nil {
		log.Println("env file not found, using system environment variables")
	}

	app := fiber.New()
	config.InitDB()

	// redis := config.Redis()
	// authservice:=services.NewAuthService(repo, redis)

	app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:5173",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Authorization",
        AllowCredentials: true,
        MaxAge:           int((12 * time.Hour).Seconds()),
    }))

	repo := repositories.NewRepo(config.DB)

	authservice:=services.NewAuthService(repo, config.Redis)
	userService := services.NewUserService(repo)


	authcontroller := controller.NewAuthController(authservice)
	userController := controller.NewUserController(userService)

	//Routes
	routes.Routes(app, authcontroller)
	routes.UserRoutes(app, userController)

	 app.Listen(":3000")
}