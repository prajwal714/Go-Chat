package controller

import (
	"Go-Chat/go.mod/constants"
	"Go-Chat/go.mod/contracts"
)

type TryAvatars []contracts.Avatar

var avatars contracts.Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	// UseGravatar,
}

//ErrNoAvatarURL is returned when the avatar instace is unable to provide an avatar from the url

func (a TryAvatars) GetAvatarURL(u contracts.ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", constants.ErrNoAvatarURL
}

var UseAuthAvatar contracts.AuthAvatar

var UseGravatar contracts.GravatarAvatar

var UseFileSystemAvatar contracts.FileSystemAvatar
