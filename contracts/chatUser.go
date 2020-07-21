package contracts

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}
