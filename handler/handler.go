package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mvanbrummen/got-std/gotgit"
	git "github.com/mvanbrummen/got-std/gotgit"
	"github.com/mvanbrummen/got-std/model"
)

type Handler struct{}

func (h *Handler) FileHandler(c *gin.Context) {
	repoName := c.Param("repoName")
	fileName := c.Param("rest")

	// strip the leading /
	fileName = fileName[1:]

	repo, _ := git.Open(repoName)
	contents, _ := git.FileContents(repo, fileName)

	fileDetail := model.FileDetail{
		RepoName: repoName,
		Name:     fileName,
		Contents: contents,
	}

	c.HTML(http.StatusOK, "file.html", fileDetail)
}

func (h *Handler) IndexHandler(c *gin.Context) {
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

	c.HTML(http.StatusOK, "index.html", repos)
}

func (h *Handler) RepositoryHandler(c *gin.Context) {
	repo, _ := git.Open(c.Param("repoName"))
	totalCommits, _ := git.TotalCommits(repo)
	totalBranches, _ := git.TotalBranches(repo)

	files, _ := git.Files(repo)

	recentCommits, _ := git.RecentCommits(repo, 5)

	repoDetail := model.RepositoryDetail{
		Name:          c.Param("repoName"),
		TotalCommits:  totalCommits,
		TotalBranches: totalBranches,
		Files:         files,
		RecentCommits: recentCommits,
	}

	c.HTML(http.StatusOK, "repository.html", repoDetail)
}

func (h *Handler) RepositoryHandlerPost(c *gin.Context) {
	repo, _ := git.Open(c.Param("repoName"))
	totalCommits, _ := git.TotalCommits(repo)
	totalBranches, _ := git.TotalBranches(repo)

	q := c.PostForm("q")

	files, _ := git.FilesFilter(repo, q)

	recentCommits, _ := git.RecentCommits(repo, 5)

	repoDetail := model.RepositoryDetail{
		Name:          c.Param("repoName"),
		TotalCommits:  totalCommits,
		TotalBranches: totalBranches,
		Files:         files,
		RecentCommits: recentCommits,
	}

	c.HTML(http.StatusOK, "repository.html", repoDetail)
}

func (h *Handler) CommitsHandler(c *gin.Context) {
	repoName := c.Param("repoName")
	repo, _ := git.Open(repoName)

	cc, _ := gotgit.Commits(repo)

	commits := model.Commits{
		RepoName: repoName,
		Commits:  cc,
	}

	c.HTML(http.StatusOK, "commits.html", commits)
}
