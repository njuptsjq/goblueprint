package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// temp reprents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP requests
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	http.Handle("/", &templateHandler{filename: "chat.html"})

	//start server
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
