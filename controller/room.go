package controller

import (
	"log"
	"net/http"

	"Go-Chat/go.mod/contracts"
	// "Go-Chat/go.mod/trace"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type Room struct {
	// forward channel that holds incomig msgs
	//that should be forwarded to other clients
	forward chan *contracts.Message

	//join is the channel for clients wishing to
	//join this channel
	join chan *Client

	//leave is the channel for
	// clients wishing to leave the room
	leave chan *Client

	//clients holds all the current clients in this room
	clients map[*Client]bool

	//tracer will recieve trace info of the activity in the room
	// tracer trace.Tracer
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.join:
			//joining clients
			r.clients[client] = true
			log.Println("New Client Joined")

		case client := <-r.leave:
			//leaving the room
			delete(r.clients, client)
			close(client.send)
			log.Println("Leave Client left")

		case msg := <-r.forward:
			log.Println("Message recieved: ", msg.Message)
			//forward message to all clients
			for client := range r.clients {
				client.send <- msg
				log.Println(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie: ", err)
		return
	}

	client := &Client{
		socket:   socket,
		send:     make(chan *contracts.Message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}

//new room makes a new room that is ready to go
func NewRoom() *Room {
	return &Room{
		forward: make(chan *contracts.Message),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}
}
