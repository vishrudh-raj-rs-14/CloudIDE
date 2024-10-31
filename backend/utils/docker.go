package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func SpawnContainerFromLocalImage(imageName string, containerPort string) (string, string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", "", fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Check if the image exists locally
	_, _, err = cli.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return "", "", fmt.Errorf("local image not found: %v", err)
	}

	// Create a container
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(containerPort + "/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "0", // Assign a random port
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			nat.Port(containerPort + "/tcp"): struct{}{},
		},
	}, hostConfig, nil, nil, "")
	if err != nil {
		return "", "", fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", "", fmt.Errorf("failed to start container: %v", err)
	}

	// Get container info to retrieve the assigned port
	containerInfo, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to inspect container: %v", err)
	}

	hostPort := containerInfo.NetworkSettings.Ports[nat.Port(containerPort+"/tcp")][0].HostPort

	return resp.ID, hostPort, nil
}

func StopContainer(containerID string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}
	noWaitTimeout := 0
	err = cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &noWaitTimeout})
	if err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}

	return nil
}

func CheckMyRepl(c *fiber.Ctx) error{
	replId := c.Params("replId");
	userId := c.Locals("userID")

	
	userID, ok := primitive.ObjectIDFromHex(userId.(string))
	if ok!=nil {
		fmt.Println("Couldnt convert")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": ok.Error(),
		})
	}
	repls, err := GetReplsByUserID(userID)
	if(err!=nil){
		fmt.Printf(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User does not exist",
		})
	}

	replPresent := false;
	var curRepl models.Repl
	for _, repl := range repls {
		if repl.ID.Hex() == replId {
			replPresent = true;
			curRepl = repl;
			c.Locals("containerId", repl.ContainerID)
			break
		}
	}
	if(!replPresent){
		fmt.Println("You dont own such a repl")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "You dont own such a repl",
		})
	}
	if(curRepl.Status!="RUNNING"){
		fmt.Println("Container is not running")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Your Current Repl is not active",
		})
	}
	return c.Next();
}