package utils

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)


func CheckSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}