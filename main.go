package main

//set the active Avatar implementation
import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"Go-Chat/go.mod/common"
	"Go-Chat/go.mod/config"
	"Go-Chat/go.mod/controller"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// HTTP handler to serve templates
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}
func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/chat", 301)
}
func main() {

	common.Initialize()

	flag.Parse() //parse the flags
	ClientID := config.GithubClientID()
	ClientSecret := config.GithubClientSecret()
	// log.Printf("Cllient ID %s, Client secret: %s", ClientID, ClientSecret)
	gomniauth.SetSecurityKey(config.SecretKey())
	gomniauth.WithProviders(
		github.New(ClientID, ClientSecret, config.CallbackURL()),
	)

	r := controller.NewRoom()
	// r.tracer = trace.New(os.Stdout)

	http.HandleFunc("/", redirect)
	http.Handle("/chat", controller.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	//handles the upload form to upload a custom avatar
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	//handles the uploader function to save the uploaded image and save it in avatars dir
	http.HandleFunc("/uploader", controller.UploaderHandler)
	//handles authentication via Github
	http.HandleFunc("/auth/", controller.LoginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	//serving the statis avatar files to the server
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))

	go r.Run()

	//start the web server

	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", config.Port()) // Default port if not specified
	}

	portInfo := ":" + port
	log.Println("Starting the web server at: ", portInfo)
	if err := http.ListenAndServe(portInfo, nil); err != nil {
		log.Fatal("Listen And serve: ", err)
	}
}
