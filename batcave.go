package main

import (
	"log"
	"net/http"

	"github.com/ffhenkes/battrack"
	"github.com/gorilla/websocket"
)

type batcave struct {
	forward chan []byte

	join chan *batclient

	leave chan *batclient

	batclients map[*batclient]bool

	battracker battrack.BatTracker
}

func newCave() *batcave {
	return &batcave{
		forward:    make(chan []byte),
		join:       make(chan *batclient),
		leave:      make(chan *batclient),
		batclients: make(map[*batclient]bool),
	}
}

func (bcave *batcave) run() {

	for {
		select {
		case batclient := <-bcave.join:
			// joining
			bcave.batclients[batclient] = true
			bcave.battracker.Trace("New client joined")
		case batclient := <-bcave.leave:
			// leaving
			delete(bcave.batclients, batclient)
			close(batclient.send)
			bcave.battracker.Trace("Client left")
		case msg := <-bcave.forward:
			// all clients
			for batclient := range bcave.batclients {
				select {
				case batclient.send <- msg:
					bcave.battracker.Trace(" -- sent to client")
					// send the message
				default:
					// failed
					delete(bcave.batclients, batclient)
					close(batclient.send)
					bcave.battracker.Trace(" -- failed to send, cleaned up client")
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
