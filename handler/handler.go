package handler

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"

	git "github.com/mvanbrummen/got-std/gotgit"
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

	repo, _ := git.Open(repoName)
	contents, _ := git.FileContents(repo, fileName)

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

	repo, _ := git.Open(vars["repoName"])
	totalCommits, _ := git.TotalCommits(repo)
	totalBranches, _ := git.TotalBranches(repo)
	files, _ := git.Files(repo)

	repoDetail := model.RepositoryDetail{
		Name:          vars["repoName"],
		TotalCommits:  totalCommits,
		TotalBranches: totalBranches,
		Files:         files,
	}

	h.templates["repository.html"].Execute(w, repoDetail)
}