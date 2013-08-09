package main

import (
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
)

//// to compile a go get-able repository
type gojob struct {
	importPath string
	request    *request
}

// a request, refs/requests/*
type request struct {
	ref  string
	time int64
	// jobs in undefined order
	jobs      []gojob
	repoquest *repoquest
	reference *reference
}

// a repository with requests
type repoquest struct {
	name     string
	requests map[string]*request
	graph    *graph
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
	graph      *graph
}

type graph struct {
	repoquests   map[string]*repoquest
	repositories map[string]*repository
}

func newGraph() *graph {
	repoquests := make(map[string]*repoquest)
	repositories := make(map[string]*repository)
	return &graph{repoquests, repositories}
}

func (j *gojob) goGetJob() GoGetJob {
	repoquest := j.request.repoquest.name
	request := j.request.ref
	return GoGetJob{j.importPath, repoquest, request}
}
