package main

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
