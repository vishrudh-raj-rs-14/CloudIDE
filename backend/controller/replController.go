package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/models"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserRepls(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if(userID == nil){
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	userIdHex, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
		})
	}
	repls, err := utils.GetReplsByUserID(userIdHex)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get repls",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Repls fetched successfully",
		"repls":   repls,
	})
}

func CreateRepl(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if(userID == nil){
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	userIdHex, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
		})
	}
	var newRepl models.Repl
	if err := c.BodyParser(&newRepl); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	if newRepl.Name == "" || newRepl.Description == "" || newRepl.Language == "" || newRepl.Framework == "" || newRepl.Visibility == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "All fields are required",
		})
	}
	newRepl.OwnerID = userIdHex
	newRepl.CreatedAt = time.Now()
	newRepl.UpdatedAt = time.Now()
	newRepl.Status = "STOPPED"

	if(newRepl.ImageName==""){
		newRepl.ImageName="express-docker-app"
	}

	insertedID, err := utils.CreateRepl(&newRepl)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
		})
	}
	newRepl.ID = insertedID
	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Repl created successfully",
		"repl": newRepl,
	})
}

func DeleteRepl(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	
	replID := c.Query("id")
	if replID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Repl ID is required",
		})
	}
	replIDHex, err := primitive.ObjectIDFromHex(replID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Repl ID",
			"error":   err.Error(),
		})
	}

	// Check if the repl exists
	exists, err := utils.ReplExists(replIDHex)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Repl does not exist",
		})
	}

	err = utils.DeleteRepl(replIDHex)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete repl",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Repl deleted successfully",
	})
}



func PerformAction(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	action := c.Query("action")
	if action == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Action is required",
		})
	}
	if action == "start" {
		return StartRepl(c)
	}else if action == "stop" {
		return StopRepl(c)
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  "error",
		"message": "Invalid action",
	})
}

func StartRepl(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	replID := c.Params("replID")
	if replID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Repl ID is required",
		})
	}
	replIDHex, err := primitive.ObjectIDFromHex(replID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Repl ID",
			"error":   err.Error(),
		})
	}
	repl, err := utils.GetRepl(replIDHex)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	if repl.Status == "RUNNING" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Repl already running",
		})
	}
	if(repl.ImageName==""){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No Repl image specified",
		})
	}
	repl.Status = "RUNNING"
	containerID, port, err := utils.SpawnContainerFromLocalImage(repl.ImageName, "3000")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	repl.ContainerID = containerID
	repl.ContainerPort = port
	err = utils.UpdateRepl(replIDHex, repl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Repl started successfully",
		"repl": repl,
	})
}

func StopRepl(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	replID := c.Params("replID")
	if replID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Repl ID is required",
		})
	}
	replIDHex, err := primitive.ObjectIDFromHex(replID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Repl ID",
			"error":   err.Error(),
		})
	}
	repl, err := utils.GetRepl(replIDHex)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	if repl.Status == "STOPPED" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Repl already stopped",
		})
	}
	repl.Status = "STOPPED"
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	err = utils.StopContainer(repl.ContainerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	repl.ContainerID = ""
	repl.ContainerPort = ""
	err = utils.UpdateRepl(replIDHex, repl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Repl stopped successfully",
		"repl": repl,
	})
}