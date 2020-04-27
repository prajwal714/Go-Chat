package main

//this is the test file for avatar.go , here we are creating an empty client with no avatar URL

//after returning an error we then set the value of avatar url to a gravater url then test again that now it isnt returning any error
import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar

	//creating an empty client
	client := new(client)

	//fetching url from empty client returns err
	url, err := authAvatar.GetAvatarURL(client)

	//we check if err == ErrNoAvatarURL, since it should be returned
	if err != ErrNoAvatarURL {
		t.Error("AuthAVatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}
	//now we set our testURL to client avatar_url
	//set a vlaue
	testURL := "http://url-to-gravatar/"
	//we insert in the client userdata map
	client.userData = map[string]interface{}{
		"avatar_url": testURL,
	}

	//now we check for error again, since now the url of avatar is not empty
	url, err = authAvatar.GetAvatarURL(client)

	//this time err returned will not be empty
	if err != nil {
		t.Error("AuthAvatar.GetAVatarURL should return no error when value present")
	}

	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}

}
