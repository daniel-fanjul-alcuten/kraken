package main

import (
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
)

func (g *graph) addRequest(wreq Request) []gojob {

	rq, ok := g.repoquests[wreq.Repoquest]
	if !ok {
		requests := make(map[string]*request)
		rq = &repoquest{wreq.Repoquest, requests, g}
		g.repoquests[wreq.Repoquest] = rq
	}

	repo, ok := g.repositories[wreq.Repository]
	if !ok {
		references := make(map[string]*reference)
		repo = &repository{wreq.Repository, references, g}
		g.repositories[wreq.Repository] = repo
	}

	ref, ok := repo.references[wreq.Reference]
	if !ok {
		srequests := make([]map[*request]*request, 0)
		mrequests := make(map[int64]map[*request]*request, 0)
		ref = &reference{wreq.Reference, srequests, mrequests, repo}
		repo.references[wreq.Reference] = ref
	}

	req, ok := rq.requests[wreq.Request]
	if !ok {
		jobs := make([]gojob, len(wreq.Jobs))
		req = &request{wreq.Request, wreq.Time, jobs, rq, ref}
		rq.requests[wreq.Request] = req
		for i, job := range wreq.Jobs {
			jobs[i] = gojob{job.ImportPath, req}
		}
		m, ok := ref.mrequests[wreq.Time]
		if !ok {
			m = make(map[*request]*request)
			ref.srequests = append(ref.srequests, m)
			ref.mrequests[wreq.Time] = m
		}
		m[req] = req
		return jobs
	}
	return nil
}
