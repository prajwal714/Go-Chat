# Go-Chat
Web-Based Chat Application using web sockets that allows multiple users to have realtime conversation.
Implemented Authentication using Github OAuth, we can also add Authentication via Google and Facebook using gomniauth package of golang.
Seperate Go Routines for Server and Chat Rooms. Careful Modelling of chatroom, users and avatar structs.
It uses github avatars, gravatar Avatars or upload ther own avatars as chatting icons.

## Packages Used
```bash

"flag"
"log"
"sync"
"net/http"  
"io/ioutil"
"path/filepath"
"text/template"

"github.com/gorilla/websocket"
"github.com/stretchr/gomniauth"
"github.com/stretchr/gomniauth/providers/github"
"github.com/stretchr/objx"
```
## Structures Implemented

Wrapper for handling HTML templates and serving them to API endpoints
```bash
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}
```
Chat Room structure which handles the joining and leaving clients, and the forwarded messages in the Chat Room
```bash
type room struct {
	// forward channel that holds incomig msgs
	//that should be forwarded to other clients
	forward chan *message

	//join is the channel for clients wishing to
	//join this channel
	join chan *client

	//leave is the channel for
	// clients wishing to leave the room
	leave chan *client

	//clients holds all the current clients in this room
	clients map[*client]bool

	//tracer will recieve trace info of the activity in the room
	tracer trace.Tracer
}
```
Client Struct holds the reference of websocket, stores info about the User in a map, 
also the channel for sending  messages and the room reference client is associated with. 
```bash
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

```
Chat User stuct holds the unique ID of user and User profile name using gomniauth
```bash
type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}
```
## Author
[Prajwal Singh](https://prj-prajwal.netlify.app/)

## License
[MIT](https://choosealicense.com/licenses/mit/)
