package controller

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/socket"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/utils"
)


func HandleTerminalConnection(c *websocket.Conn){
	containerId :=  c.Locals("containerId").(string);
	executor, _ := utils.NewDockerExecutor(containerId)
	client := &socket.Client{
		Conn:c,
		ContainerId:containerId,
		DockerExecutor: executor,
		Send: make(chan []byte),
	}

	var wg sync.WaitGroup
    wg.Add(2) // We have two goroutines

    go func() {
        defer wg.Done()
        client.ReadPump()
    }()

    go func() {
        defer wg.Done()
        client.WritePump()
    }()

    wg.Wait()

}