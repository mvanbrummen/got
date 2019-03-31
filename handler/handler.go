package handler

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/mvanbrummen/got-std/model"
)

type Handler struct {
	templates map[string]*template.Template
}

func NewHandler(templates map[string]*template.Template) *Handler {
	return &Handler{
		templates: templates,
	}
}

func (h *Handler) FileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repoName"]
	fileName := vars["rest"]

	repo, _ := git.PlainOpen(".got/" + repoName + "/.git")
	ref, _ := repo.Head()

	c, _ := repo.CommitObject(ref.Hash())

	iter, _ := c.Files()

	var contents string
	iter.ForEach(func(f *object.File) error {
		if f.Name == fileName {
			cont, _ := f.Contents()
			contents = cont
		}
		return nil
	})

	fileDetail := model.FileDetail{
		RepoName: repoName,
		Name:     fileName,
		Contents: contents,
	}

	h.templates["file.html"].Execute(w, fileDetail)
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	repos := model.RepositoryList{
		[]model.Repository{
			model.Repository{
				"go-bencoder2",
			},
			model.Repository{
				"go-git",
			},
		},
	}

	h.templates["index.html"].Execute(w, repos)
}

func (h *Handler) RepositoryHandler(w http.ResponseWriter, r *http.Request) {
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

	files := make([]model.File, 0)
	iter.ForEach(func(f *object.File) error {
		files = append(files, model.File{f.Name, f.Hash.String()})
		return nil
	})

	cIter, _ := repo.Log(&git.LogOptions{From: ref.Hash()})

	var cCount int
	cIter.ForEach(func(c *object.Commit) error {
		cCount++

		return nil
	})

	repoDetail := model.RepositoryDetail{
		Name:          vars["repoName"],
		TotalCommits:  cCount,
		TotalBranches: totalBranches,
		Files:         files,
	}

	h.templates["repository.html"].Execute(w, repoDetail)
}
