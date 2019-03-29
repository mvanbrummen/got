package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

const (
	templatesDir = "templates/"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates["index.html"].Execute(w, nil)
}

var templates map[string]*template.Template

func init() {
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		panic(err)
	}

	templates = make(map[string]*template.Template, 0)

	for _, f := range files {
		templates["index.html"] = template.Must(template.ParseFiles(templatesDir + f.Name()))
	}
}

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
