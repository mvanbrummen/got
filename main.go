package main

import (
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/mvanbrummen/got-std/handler"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

const (
	templatesDir = "templates/"
	staticDir    = "static/"
)

var templates map[string]*template.Template

func init() {
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		panic(err)
	}

	templates = make(map[string]*template.Template, 0)

	for _, f := range files {
		templates[f.Name()] = template.Must(template.ParseFiles(templatesDir + f.Name()))
	}
}

func main() {
	r := mux.NewRouter()

	handler := handler.NewHandler(templates)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	r.HandleFunc("/repository/{repoName}/blob/{rest:.*}", handler.FileHandler)
	r.HandleFunc("/repository/{repoName}", handler.RepositoryHandler)
	r.HandleFunc("/", handler.IndexHandler)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
