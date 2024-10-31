package controller

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/vishrudh-raj-rs-14/cloudIDEbackend/socket"
)


func HandleTerminalConnection(c *websocket.Conn){
	client := &socket.Client{
		Conn:c,
		ContainerId: c.Locals("containerId").(string),
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