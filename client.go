package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type client struct {

	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	send chan *message

	// room is the room this client is chatting in.
	room *room

	//stores info about the user
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}

		msg.Name = c.userData["name"].(string)
		msg.When = time.Now()

		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)

		c.room.forward <- msg
	}

}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}

	}
}
