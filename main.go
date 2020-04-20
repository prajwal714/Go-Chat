package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

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
	t.templ.Execute(w, r)
}

func main() {

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() //parse the flags

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	//start the web server
	log.Println("Starting the web server at: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listen And serve: ", err)
	}
}
