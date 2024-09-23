package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/controller"
)



func Auth_router(app fiber.Router){
	app.Post("/register", controller.Register)
	app.Post("/login", controller.Login)
	app.Post("/logout", controller.Logout)
}