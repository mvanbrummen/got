package main

import (
	"net/http"

	"github.com/spf13/viper"

	"github.com/mvanbrummen/got-std/handler"
	"github.com/mvanbrummen/got-std/util"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

const (
	templatesDir = "templates/"
	staticDir    = "static/"
)

var templates util.Templates

func init() {
	// init config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// init templates
	var err error
	if templates, err = util.InitTemplates(templatesDir); err != nil {
		panic(err)
	}

	// create the app directory
	util.InitGotDir(viper.GetString("got.dir"))
}

func main() {
	r := mux.NewRouter()

	handler := handler.NewHandler(templates)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	r.HandleFunc("/repository/{repoName}/blob/{rest:.*}", handler.FileHandler)
	r.HandleFunc("/repository/{repoName}", handler.RepositoryHandler)
	r.HandleFunc("/repository/{repoName}/commits", handler.CommitsHandler)
	r.HandleFunc("/", handler.IndexHandler)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(viper.GetString("application.port"), r))
}
