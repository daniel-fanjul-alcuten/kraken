package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
)

//// to compile a go get-able repository
type goget struct {
	importPath string
	request    *request
}

// a request, refs/requests/*
type request struct {
	ref  string
	time int64
	// jobs in undefined order
	jobs      []goget
	repoquest *repoquest
	reference *reference
}

// a repository with requests
type repoquest struct {
	name     string
	requests map[string]*request
	graph    *Graph
}

// the ref in the original repository
type reference struct {
	name       string
	srequests  []map[*request]*request
	mrequests  map[int64]map[*request]*request
	repository *repository
}

// a repository with references
type repository struct {
	name       string
	references map[string]*reference
	graph      *Graph
}

type Graph struct {
	repoquests   map[string]*repoquest
	repositories map[string]*repository
}

func NewGraph() *Graph {
	repoquests := make(map[string]*repoquest)
	repositories := make(map[string]*repository)
	return &Graph{repoquests, repositories}
}

func (j *goget) goGetJob() GoGetJob {
	repoquest := j.request.repoquest.name
	request := j.request.ref
	return GoGetJob{j.importPath, repoquest, request}
}
