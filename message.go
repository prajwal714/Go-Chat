package main

import (
	"time"
)

//we difine a struct for our message
type message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
}
