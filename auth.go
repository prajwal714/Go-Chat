package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// _ is used to collect cookie, we discardd it coz we are only interested if it is presnt or not.
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
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
	action := segs[2]
	provider := segs[3]

	switch action {
	case "login":
		log.Println("Todo Handle loging for: ", provider)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}

}
