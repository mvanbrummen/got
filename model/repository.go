package model

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
	RecentCommits []Commit
}
