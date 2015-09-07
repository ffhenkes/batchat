package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/ffhenkes/battrack"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, r)
}

func main() {

	var addr = flag.String("addr", ":8000", "Batcave entry.")
	flag.Parse()

	bcave := newCave()
	bcave.battracker = battrack.New(os.Stdout)

	http.Handle("/", &templateHandler{filename: "batchat.html"})
	http.Handle("/batcave", bcave)

	// get the batcave going
	go bcave.run()
	log.Println("The Batcave Chat is rolling on port", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
