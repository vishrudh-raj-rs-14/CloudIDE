package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/controller"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/middleware"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/utils"
)


func Terminal_Router(app fiber.Router){
	app.Use(utils.CheckSocket, middleware.Protect, utils.CheckMyRepl)
	app.Get("/", websocket.New(func(c *websocket.Conn) {
		controller.HandleTerminalConnection(c)
	}))
	// app.Post("/run", controller.Run)
}