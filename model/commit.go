package model

import "time"

type Commits struct {
	RepoName string
	Commits  []Commit
}

type Commit struct {
	Hash         string
	ShortHash    string
	Message      string
	Author       Author
	ParentHashes []string
}

type Author struct {
	Name  string
	Email string
	When  time.Time
}
