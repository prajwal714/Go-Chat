package controller

import (
	"Go-Chat/go.mod/contracts"
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type Client struct {

	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	send chan *contracts.Message

	// room is the room this client is chatting in.
	room *Room

	//stores info about the user
	userData map[string]interface{}
}

func (c *Client) read() {
	defer c.socket.Close()

	for {
		var msg *contracts.Message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}

		msg.Name = c.userData["name"].(string)
		msg.When = time.Now()

		if avatarUrl, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarUrl.(string)
		}

		c.room.forward <- msg
	}

}

func (c *Client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}

	}
}
