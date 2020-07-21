package controller

//this is the test file for avatar.go , here we are creating an empty client with no avatar URL

//after returning an error we then set the value of avatar url to a gravater url then test again that now it isnt returning any error
import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

// func TestAuthAvatar(t *testing.T) {
// 	var authAvatar AuthAvatar

// 	//creating an empty client
// 	client := new(client)

// 	//fetching url from empty client returns err
// 	url, err := authAvatar.GetAvatarURL(client)

// 	//we check if err == ErrNoAvatarURL, since it should be returned
// 	if err != ErrNoAvatarURL {
// 		t.Error("AuthAVatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
// 	}
// 	//now we set our testURL to client avatar_url
// 	//set a vlaue
// 	testURL := "http://url-to-gravatar/"
// 	//we insert in the client userdata map
// 	client.userData = map[string]interface{}{
// 		"avatar_url": testURL,
// 	}

// 	//now we check for error again, since now the url of avatar is not empty
// 	url, err = authAvatar.GetAvatarURL(client)

// 	//this time err returned will not be empty
// 	if err != nil {
// 		t.Error("AuthAvatar.GetAVatarURL should return no error when value present")
// 	}

// 	if url != testURL {
// 		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
// 	}

// }
func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)

	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURLwhen no value present")
	}

	testUrl := "http://url-to-gravatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)

	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no errorwhen value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

//this is the test code for gravatar same as test code for auth avatar
func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}

	url, err := gravatarAvatar.GetAvatarURL(user)

	if err != nil {
		t.Error("GravatarAVatar.GetAvatarURL should not return error")
	}

	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

//this is the test code for our uploaded avatar file
func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer os.Remove(filename)

	var fileSystemAvatar FileSystemAvatar

	user := &chatUser{uniqueID: "abc"}

	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return error")
	}

	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
