package main

import (
	"github.com/gorilla/websocket"
)

type batclient struct {
	socket *websocket.Conn

	send chan []byte

	batcave *batcave
}

func (c *batclient) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.batcave.forward <- msg
		} else {
			break
		}
	}

	c.socket.Close()
}

func (c *batclient) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}

	c.socket.Close()
}
