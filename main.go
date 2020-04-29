package main

//set the active Avatar implementation
import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
)

//we add all 3 implementations of the Avatars, it implements whichever returns the Avatar URL
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	// UseGravatar,
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
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
func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/chat", 301)
}
func main() {

	godotenv.Load(".env")

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() //parse the flags
	ClientID := goDotEnvVariable("GITHUB_CLIENT_ID")
	ClientSecret := goDotEnvVariable("GITHUB_CLIENT_SECRET")
	log.Printf("Cllient ID %s, Client secret: %s", ClientID, ClientSecret)
	gomniauth.SetSecurityKey("admin1234")
	gomniauth.WithProviders(
		github.New(ClientID, ClientSecret, "http://localhost:8080/auth/callback/github"),
	)

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)

	http.HandleFunc("/", redirect)
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
