package routes

import (
	controller "marryo/Internal/Controllers"
	middleware "marryo/Internal/MiddleWare"

	"github.com/gofiber/fiber/v2"
	ws "github.com/gofiber/websocket/v2"
)

func ChatRoute(app *fiber.App, chatController *controller.ChatController) {

	// WebSocket upgrade check
	app.Use("/chat/ws", func(c *fiber.Ctx) error {
		if ws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	chat := app.Group("/chat")
	chat.Get("/ws", ws.New(chatController.ChatWebSocket))

	//Get Chat Datas
	chats := app.Group("/chats")
	chats.Use(middleware.AuthMiddleWare())

	chats.Get("/", chatController.GetChats)
	chats.Get("/:match_id/messages", chatController.GetMessages)
	chats.Post("/:match_id/messages", chatController.SendMessage)

}
