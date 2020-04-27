package main

//set the active Avatar implementation
import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
)

var avatars Avatar = UseFileSystemAvatar

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//HTTP handler to serve templates
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

func main() {

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() //parse the flags

	gomniauth.SetSecurityKey("admin1234")
	gomniauth.WithProviders(
		github.New("22f0dcb22b1b50033d6d", "0d08156c1c332d6bda74be24803e715818de015c", "http://localhost:8080/auth/callback/github"),
	)

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	//handles the upload form to upload a custom avatar
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	//handles the uploader function to save the uploaded image and save it in avatars dir
	http.HandleFunc("/uploader", uploaderHandler)
	//handles authentication via Github
	http.HandleFunc("/auth/", loginHandler)
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

	go r.run()

	//start the web server
	log.Println("Starting the web server at: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listen And serve: ", err)
	}
}
