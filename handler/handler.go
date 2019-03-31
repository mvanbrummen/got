package handler

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"

	"github.com/mvanbrummen/got-std/gotgit"
	git "github.com/mvanbrummen/got-std/gotgit"
	"github.com/mvanbrummen/got-std/model"
	"github.com/mvanbrummen/got-std/util"
)

type Handler struct {
	templates util.Templates
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
			model.Repository{
				"vscode",
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
	recentCommits, _ := git.RecentCommits(repo, 5)

	repoDetail := model.RepositoryDetail{
		Name:          vars["repoName"],
		TotalCommits:  totalCommits,
		TotalBranches: totalBranches,
		Files:         files,
		RecentCommits: recentCommits,
	}

	h.templates["repository.html"].Execute(w, repoDetail)
}

func (h *Handler) CommitsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	repo, _ := git.Open(vars["repoName"])

	c, _ := gotgit.Commits(repo)

	commits := model.Commits{
		RepoName: vars["repoName"],
		Commits:  c,
	}

	h.templates["commits.html"].Execute(w, commits)
}
