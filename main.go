package main

import (
	"io/ioutil"
	"net/http"
	"text/template"

	"gopkg.in/src-d/go-git.v4/plumbing"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/gorilla/mux"
)

const (
	templatesDir = "templates/"
	staticDir    = "static/"
)

type RepositoryList struct {
	Repositories []Repository
}

type Repository struct {
	Name string
}

type RepositoryDetail struct {
	Name          string
	Files         []File
	TotalCommits  int
	TotalBranches int
}

type File struct {
	Name string
	Hash string
}

type FileDetail struct {
	RepoName string
	Name     string
	Contents string
	Hash     string
}

func repositoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	repo, _ := git.PlainOpen(".got/" + vars["repoName"] + "/.git")
	ref, _ := repo.Head()

	bIter, _ := repo.Branches()

	totalBranches := 0
	bIter.ForEach(func(p *plumbing.Reference) error {
		totalBranches++
		return nil
	})

	c, _ := repo.CommitObject(ref.Hash())

	iter, _ := c.Files()

	files := make([]File, 0)
	iter.ForEach(func(f *object.File) error {
		files = append(files, File{f.Name, f.Hash.String()})
		return nil
	})

	cIter, _ := repo.Log(&git.LogOptions{From: ref.Hash()})

	var cCount int
	cIter.ForEach(func(c *object.Commit) error {
		cCount++

		return nil
	})

	repoDetail := RepositoryDetail{
		Name:          vars["repoName"],
		TotalCommits:  cCount,
		TotalBranches: totalBranches,
		Files:         files,
	}

	templates["repository.html"].Execute(w, repoDetail)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	repos := RepositoryList{
		[]Repository{
			Repository{
				"go-bencoder2",
			},
		},
	}

	templates["index.html"].Execute(w, repos)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repoName"]
	fileName := vars["rest"]

	fileDetail := FileDetail{
		RepoName: repoName,
		Name:     fileName,
		Contents: `
			println("hello world!")
			exit(0)
		`,
	}

	templates["file.html"].Execute(w, fileDetail)
}

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
	handler := mux.NewRouter()

	handler.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	handler.HandleFunc("/repository/{repoName}/blob/{rest:.*}", fileHandler)
	handler.HandleFunc("/repository/{repoName}", repositoryHandler)
	handler.HandleFunc("/", indexHandler)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
