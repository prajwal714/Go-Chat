package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// cookie is used to collect cookie, we  candiscardd it if are only interested if it is presnt or not.
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie || cookie.Value == "" {
		//not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		//some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//success-- call the next handler
	h.next.ServeHTTP(w, r)
}

//it accepts the http handler function checks for authentication and
//forwards towards the next http handler function
//it acts as a middleware
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

//login Handler handles third-party login process.
//format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	//we are breaking the path into segs to extract the action and provider
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]   //login action
	provider := segs[3] //login using google or Github or facebook

	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}

		loginUrl, err := provider.GetBeginAuthURL(nil, nil) //Get Begin Auth URL contains the url we must redirect the user to begin logging

		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s", provider, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback": //this is the callback which we will be rddirected after OAuth from provider
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		// we parse the Raw query and completeAuth function complete the OAuth2 provider handshake with provider
		//creds store the credentials of the user
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while trying to complete auth for %s: %s", provider, err), http.StatusInternalServerError)
			return
		}
		//we extract the user name from the provided creds
		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while trying to get user %s: %s", provider, err), http.StatusInternalServerError)
			return
		}
		chatUser := &chatUser{User: user}
		//we use hash value of user email as our user ID
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))

		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("Error while trying to get AvatarURL", "-", err)
		}

		// we base encode the name of user to store in our auth cookie
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,
			"email":      user.Email(),
		}).MustBase64()
		//after setting the aut cookie we can redirect to the /chat and it will successfully pass the AuthHandler
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})

		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}

}
