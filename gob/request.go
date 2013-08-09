package kraken

type GoGetRequest struct {
	ImportPath string
}

type Request struct {
	Repoquest, Request    string
	Repository, Reference string
	Time                  int64
	Jobs                  []GoGetRequest
}
