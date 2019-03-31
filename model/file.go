package model

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
