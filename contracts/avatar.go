package contracts

import (
	"Go-Chat/go.mod/constants"
	"io/ioutil"
	"path"
)

type Avatar interface {
	//GetAvatarURL gets the avatar URL from the specified client,
	//returns an error if something goes wrong

	GetAvatarURL(ChatUser) (string, error)
}

type GravatarAvatar struct {
}

type AuthAvatar struct{}

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", constants.ErrNoAvatarURL
	}
	return url, nil
}

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", constants.ErrNoAvatarURL
}
