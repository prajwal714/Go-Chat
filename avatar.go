package main

import (
	"errors"
)

//ErrNoAvatarURL is returned when the avatar instace is unable to provide an avatar from the url
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

type Avatar interface {
	//GetAvatarURL gets the avatar URL from the specified client,
	//returns an error if something goes wrong

	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}