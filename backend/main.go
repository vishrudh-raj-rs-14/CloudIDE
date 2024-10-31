// main.go

package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/middleware"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/models"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/routes"
)



func main() {
	err := godotenv.Load()
	if(err!=nil){
		log.Fatal("Error loading .env file")
	}
	dbName := os.Getenv("DB_NAME");
	mongo_URI := os.Getenv("MONGO_URI");
    // Create a new Fiber instance
    app := fiber.New()
	err = models.Connect(dbName, mongo_URI)
	defer func() {
		if err := models.Mg.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	if(err !=nil){
		log.Fatal("Error Connecting to Database")
	}else{
		log.Println("Successfully Connected to Database")
	}
    // Define your routes
    app.Get( "/", middleware.Protect ,func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

	app.Route("/api/auth", routes.Auth_router)
	app.Route("/api/repl", routes.Repl_router)
	app.Route("/api/repl/terminal/:replId", routes.Terminal_Router)
	

    // Start the server
    log.Fatal(app.Listen(":3000"))
}
