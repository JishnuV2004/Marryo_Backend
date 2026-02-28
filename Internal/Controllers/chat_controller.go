package controller

import (
	services "marryo/Internal/Services"
	utils "marryo/Internal/Utils"
	websocket "marryo/Internal/WebSocket"
	"strconv"

	"github.com/gofiber/fiber/v2"
	ws "github.com/gofiber/websocket/v2"

)

type ChatController struct {
	Hub         *websocket.Hub
	ChatService *services.ChatService
}

func NewChatController(hub *websocket.Hub, chatService *services.ChatService) *ChatController {
	return &ChatController{Hub: hub, ChatService: chatService}
}

func (c *ChatController) ChatWebSocket(conn *ws.Conn) {

	tokenStr := conn.Cookies("access")
	if tokenStr == "" {
		conn.Close()
		return
	}

	// token, err := utils.Parse(tokenStr)
	// if err != nil || !token.Valid {
	// 	conn.Close()
	// 	return
	// }

	// claims := token.Claims.(jwt.MapClaims)
	// userID := uint(claims["userID"].(float64))

	claims, err := utils.VerifyToken(tokenStr)
if err != nil {
	conn.Close()
	return
}

userID := claims.UserID

	client := &websocket.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan websocket.MessagePayload),
	}

	c.Hub.Register <- client
	defer func() {
		c.Hub.Unregister <- client
		conn.Close()
	}()

	// Writer goroutine
	go func() {
		for msg := range client.Send {
			conn.WriteJSON(msg)
		}
	}()

	// Reader loop
	for {
		var payload websocket.MessagePayload
		if err := conn.ReadJSON(&payload); err != nil {
			break
		}

		//Only matched users can chat
		if !c.ChatService.Chat(userID, payload.ToUserID) {
			continue
		}

		//Save message
		err := c.ChatService.SaveMessages(payload.MatchID, userID, payload.Message)
		if err != nil {
			break
		}

		payload.FromUserID = userID
		c.Hub.Broadcast <- payload
	}
}

// Get Chats
func (c *ChatController) GetChats(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	data, err := c.ChatService.GetChats(userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": data})
}

// Get Messages (GET /chats/:match_id/messages)
func (c *ChatController) GetMessages(ctx *fiber.Ctx) error {
	matchID, _ := strconv.Atoi(ctx.Params("match_id"))

	msgs, err := c.ChatService.GetMessages(uint(matchID), 50)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"messages": msgs})
}

// Save Message (POST /chats/:match_id/messages)
func (c *ChatController) SendMessage(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	matchID, _ := strconv.Atoi(ctx.Params("match_id"))

	var body struct {
		Content string `json:"content"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := c.ChatService.SaveMessages(uint(matchID), userID, body.Content); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "sent"})
}
