package kraken

type GoGetJob struct {
	ImportPath string
	Repoquest  string
	Request    string
}

type JobResult struct {
	Success bool
	Output  string
}
