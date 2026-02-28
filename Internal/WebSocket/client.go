package websocket

import "github.com/gofiber/websocket/v2"

type Client struct {
	UserID uint
	Conn *websocket.Conn
	Send chan MessagePayload
}