package utils

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
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

type DockerExecutor struct {
    client      *client.Client
    containerID string
    workDir     string
}

// NewDockerExecutor creates a new executor for running commands in a container
func NewDockerExecutor(containerID string) (*DockerExecutor, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, fmt.Errorf("failed to create docker client: %v", err)
    }

    executor := &DockerExecutor{
        client:      cli,
        containerID: containerID,
        workDir:     "/usr/src/app",
    }

    // Get initial working directory

    return executor, nil
}

// getCurrentWorkDir executes pwd to get the current working directory
func (d *DockerExecutor) getCurrentWorkDir() (string, error) {
    ctx := context.Background()

    execConfig := types.ExecConfig{
        AttachStdout: true,
        AttachStderr: true,
        WorkingDir:   d.workDir,
        Cmd:         []string{"/bin/sh", "-c", "pwd"},
    }

    execID, err := d.client.ContainerExecCreate(ctx, d.containerID, execConfig)
    if err != nil {
        return "", fmt.Errorf("failed to create pwd exec: %v", err)
    }

    resp, err := d.client.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
    if err != nil {
        return "", fmt.Errorf("failed to attach to pwd exec: %v", err)
    }
    defer resp.Close()

    output, err := io.ReadAll(resp.Reader)
    if err != nil {
        return "", fmt.Errorf("failed to read pwd output: %v", err)
    }

    // Trim whitespace and newlines from pwd output
    return strings.TrimSpace(string(output)), nil
}

// executeRawCommand runs a single command and returns its output
func (d *DockerExecutor) executeRawCommand(cmd string) (string, error) {
    ctx := context.Background()
	wrappedCmd := fmt.Sprintf(`
        cd "%s" && {
            %s
            echo "::PWD::$(pwd)"
        }
    `, d.workDir, cmd)
    execConfig := types.ExecConfig{
        AttachStdout: true,
        AttachStderr: true,
        WorkingDir:   d.workDir,
        Cmd:         []string{"/bin/sh", "-c", wrappedCmd},
    }

    execID, err := d.client.ContainerExecCreate(ctx, d.containerID, execConfig)
    if err != nil {
        return "", fmt.Errorf("failed to create exec: %v", err)
    }

    resp, err := d.client.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
    if err != nil {
        return "", fmt.Errorf("failed to attach to exec: %v", err)
    }
    defer resp.Close()

    output, err := io.ReadAll(resp.Reader)
    if err != nil {
        return "", fmt.Errorf("failed to read output: %v", err)
    }

    return string(output), nil
}

// ExecuteCommand runs a command and updates the working directory
func (d *DockerExecutor) ExecuteCommand(cmd string) (string, error) {
    // Execute the actual command
    output, err := d.executeRawCommand(cmd)
    if err != nil {
        return output, err
    }
    if pwdIndex := strings.LastIndex(output, "::PWD::"); pwdIndex != -1 {
        newPwd := strings.TrimSpace(output[pwdIndex+7:])
		output = strings.TrimSpace(output[:pwdIndex])
        if newPwd != "" {
            d.workDir = newPwd
        }
    }

    return output, nil
}

// GetCurrentWorkDir returns the current working directory
func (d *DockerExecutor) GetCurrentWorkDir() string {
    return d.workDir
}

