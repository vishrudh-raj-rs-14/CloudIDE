package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/controller"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/middleware"
)


func Repl_router(app fiber.Router){
	app.Get("/", middleware.Protect, controller.GetUserRepls)
	app.Post("/", middleware.Protect, controller.CreateRepl)
	app.Delete("/", middleware.Protect, controller.DeleteRepl)
	app.Get("/:replID", middleware.Protect, controller.PerformAction)
	// app.Post("/run", controller.Run)
}