package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type batcave struct {
	forward chan []byte

	join chan *batclient

	leave chan *batclient

	batclients map[*batclient]bool
}

func (bcave *batcave) run() {
	for {
		select {
		case batclient := <-bcave.join:
			// joining
			bcave.batclients[batclient] = true
		case batclient := <-bcave.leave:
			// leaving
			delete(bcave.batclients, batclient)
			close(batclient.send)
		case msg := <-bcave.forward:
			// all clients
			for batclient := range bcave.batclients {
				select {
				case batclient.send <- msg:
					// send the message
				default:
					// failed
					delete(bcave.batclients, batclient)
					close(batclient.send)
				}
			}

		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (bcave *batcave) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServeHTTP", err)
		return
	}

	batclient := &batclient{
		socket:  socket,
		send:    make(chan []byte, messageBufferSize),
		batcave: bcave,
	}

	bcave.join <- batclient

	defer func() {
		bcave.leave <- batclient
	}()

	go batclient.write()
	batclient.read()
}
