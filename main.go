package main

import (
	"net/http"

	"github.com/mvanbrummen/got-std/handler"
	"github.com/mvanbrummen/got-std/util"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

const (
	templatesDir = "templates/"
	staticDir    = "static/"
	gotDir       = ".got"
)

var templates util.Templates

func init() {
	var err error
	if templates, err = util.InitTemplates(templatesDir); err != nil {
		panic(err)
	}

	util.InitGotDir(gotDir)
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
