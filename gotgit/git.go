package gotgit

import (
	"strings"

	"github.com/mvanbrummen/got-std/model"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Open opens a git repository
func Open(repoName string) (*git.Repository, error) {
	return git.PlainOpen(".got/" + repoName + "/.git")
}

// TotalBranches returns the total number of branches
func TotalBranches(repo *git.Repository) (int, error) {
	bIter, err := repo.Branches()
	if err != nil {
		return 0, err
	}

	totalBranches := 0
	bIter.ForEach(func(p *plumbing.Reference) error {
		totalBranches++
		return nil
	})

	return totalBranches, nil
}

// Commits returns all commits for a repo limited by limit
func RecentCommits(repo *git.Repository, limit int) ([]model.Commit, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	i := 0
	commits := make([]model.Commit, 0)
	cIter.ForEach(func(c *object.Commit) error {
		if i >= limit {
			return nil
		}
		parentHashes := make([]string, 0)
		for _, ph := range c.ParentHashes {
			parentHashes = append(parentHashes, ph.String())
		}

		commits = append(commits, model.Commit{
			Hash:         c.Hash.String(),
			ShortHash:    c.Hash.String()[:8],
			Message:      c.Message,
			ParentHashes: parentHashes,
			Author: model.Author{
				Name:  c.Author.Name,
				Email: c.Author.Email,
				When:  c.Author.When,
			},
		})

		i++
		return nil
	})

	return commits, nil
}

// Commits returns all commits for a repo
func Commits(repo *git.Repository) ([]model.Commit, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	commits := make([]model.Commit, 0)
	cIter.ForEach(func(c *object.Commit) error {
		parentHashes := make([]string, 0)
		for _, ph := range c.ParentHashes {
			parentHashes = append(parentHashes, ph.String())
		}

		commits = append(commits, model.Commit{
			Hash:         c.Hash.String(),
			ShortHash:    c.Hash.String()[:8],
			Message:      c.Message,
			ParentHashes: parentHashes,
			Author: model.Author{
				Name:  c.Author.Name,
				Email: c.Author.Email,
				When:  c.Author.When,
			},
		})

		return nil
	})

	return commits, nil
}

// TotalCommits returns the total number of commits
func TotalCommits(repo *git.Repository) (int, error) {
	ref, err := repo.Head()
	if err != nil {
		return 0, err
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return 0, err
	}

	var cCount int
	cIter.ForEach(func(c *object.Commit) error {
		cCount++

		return nil
	})

	return cCount, nil
}

// Files returns all the files of a repo
func Files(repo *git.Repository) ([]model.File, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	c, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	iter, err := c.Files()
	if err != nil {
		return nil, err
	}

	files := make([]model.File, 0)
	iter.ForEach(func(f *object.File) error {
		files = append(files, model.File{f.Name, f.Hash.String()})
		return nil
	})

	return files, nil
}

// Files returns all the files of a repo that match the query
func FilesFilter(repo *git.Repository, q string) ([]model.File, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	c, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	iter, err := c.Files()
	if err != nil {
		return nil, err
	}

	files := make([]model.File, 0)
	iter.ForEach(func(f *object.File) error {
		if strings.Contains(f.Name, q) {
			files = append(files, model.File{f.Name, f.Hash.String()})
		}

		return nil
	})

	return files, nil
}

// FileContents returns the contents of the file
func FileContents(repo *git.Repository, fileName string) (string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	c, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	iter, err := c.Files()
	if err != nil {
		return "", err
	}

	var contents string
	iter.ForEach(func(f *object.File) error {
		if f.Name == fileName {
			cont, _ := f.Contents()
			contents = cont
		}
		return nil
	})

	return contents, nil
}
