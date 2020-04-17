package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward channel that holds incomig msgs
	//that should be forwarded to other clients

	forward chan []byte

	//join is the channel for clients wishing to
	//join this channel

	join chan *client

	//leave is the channel for
	// clients wishing to leave the room
	leave chan *client

	//clients holds all the current clients in this room
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {

		case client := <-r.join:
			//joining clients
			r.clients[client] = true

		case client := <-r.leave:
			//leaving the room
			delete(r.clients, client)
			close(client.send)

		case msg := <-r.forward:
			//forward message to all clients
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
