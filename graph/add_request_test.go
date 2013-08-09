package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	"testing"
)

func TestAddRequest(t *testing.T) {

	repository := "server1:path1"
	reference := "refs/heads/master"
	repoquest := "server2:path2"
	request := "refs/requests/reference1"
	time := int64(5)
	importPath := "github.com/daniel-fanjul-alcuten/kraken"

	g := NewGraph()
	jobs := []GoGetRequest{GoGetRequest{importPath}}
	rjobs := g.AddRequest(Request{repoquest, request, repository, reference, time, jobs})

	rq := g.repoquests[repoquest]
	if rq == nil {
		t.Fatal(rq)
	}
	req := rq.requests[request]
	if req == nil {
		t.Fatal(req)
	}
	repo := g.repositories[repository]
	if repo == nil {
		t.Fatal(repo)
	}
	ref := repo.references[reference]
	if ref == nil {
		t.Fatal(ref)
	}
	if len(req.jobs) != 1 {
		t.Fatal(len(req.jobs))
	}
	j := req.jobs[0]

	if repo.name != repository {
		t.Error(repo.name)
	}
	if repo.graph != g {
		t.Error(repo.graph)
	}

	if ref.name != reference {
		t.Error(ref.name)
	}
	if ref.repository != repo {
		t.Error(ref.repository)
	}
	if len(ref.srequests) != 1 {
		t.Error(len(ref.srequests))
	} else if ref.srequests[0][req] != req {
		t.Error(ref.srequests[0][req])
	}
	if len(ref.mrequests) != 1 {
		t.Error(len(ref.mrequests))
	} else if ref.mrequests[time][req] != req {
		t.Error(ref.mrequests[time][req])
	}

	if rq.name != repoquest {
		t.Error(rq.name)
	}
	if rq.graph != g {
		t.Error(rq.graph)
	}

	if req.ref != request {
		t.Error(req.ref)
	}
	if req.time != time {
		t.Error(req.time)
	}
	if req.repoquest != rq {
		t.Error(req.repoquest)
	}
	if req.reference != ref {
		t.Error(req.reference)
	}

	if j.importPath != importPath {
		t.Error(j.importPath)
	}

	if len(rjobs) != 1 {
		t.Error(len(rjobs))
	} else {
		if rjobs[0].ImportPath != importPath {
			t.Error(rjobs[0].ImportPath)
		}
		if rjobs[0].Repoquest != repoquest {
			t.Error(rjobs[0].Repoquest)
		}
		if rjobs[0].Request != request {
			t.Error(rjobs[0].Request)
		}
	}
}
