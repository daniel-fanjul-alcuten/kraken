[![Build Status](https://drone.io/github.com/daniel-fanjul-alcuten/kraken/status.png)](https://drone.io/github.com/daniel-fanjul-alcuten/kraken/latest)

Description
===========

build system with the following features:

* no one at the moment

that will have these others in the future:

* compile go sources
* compile java sources
* compile c/c++ sources
* run locally and in the cloud at the same time
* distribute steps across several machines
* cache the intermediate and final artifacts

Compilation
===========

Use 'go get':
<pre>
go get github.com/daniel-fanjul-alcuten/kraken
</pre>

Deploy
======

These are the steps to deploy:

kraken-push
-----------

kraken-push pushes the refs from the current git repository that contains the code to another repository.

Every ref is not pushed as in a usual 'git push'. An annotated tag object that points to the commit referenced by the ref is created first. A tag reference is not created. The body of the message contains some metadata in JSON format: the repository and the reference the commit comes from. The tag object is pushed to a remote ref refs/requests/&lt;random hex string&gt;.

Mandatory git config entries:

* _remote.kraken.url_: the url of the git repository to push to.
* _kraken.repository_: the name that identifies this local repository.

Optional git config entries:

* _remote.kraken.skipDefaultUpdate_: some users can find this property useful.

Usage:
<pre>
Usage of kraken-push: &lt;local ref&gt;+
</pre>

Typical command line:

<pre>
git remote add kraken &lt;url&gt;
git config remote.kraken.skipDefaultUpdate true
git config kraken.repository $(hostname):$(pwd)

kraken-push master dfanjul/topic1 dfanjul/topic2
</pre>

Typical post-commit hook:

<pre>
#!/bin/bash
head=$(git symbolic-ref HEAD 2>/dev/null)
kraken-push HEAD $head
</pre>

Typical post-checkout hook:

<pre>
#!/bin/bash
old="$1" && shift || exit 1
new="$1" && shift || exit 1
flag="$1" && shift || exit 1
if [ "$flag" = 1 ]; then
  kraken-push HEAD
fi
</pre>

Typical post-rewrite hook:

<pre>
#!/bin/bash
command="$1" && shift || exit 1
while read old new; do
  :
done
if [ "$command" = rebase ]; then
  head=$(git symbolic-ref HEAD 2>/dev/null)
  kraken-push HEAD $head
fi
</pre>

Typical post-receive hook:

<pre>
#!/bin/bash
refs=""
while read old new ref; do
  refs="$refs $ref"
done
kraken-push $refs
</pre>

kraken-graph
------------

kraken-graph keeps a graph with all requests in memory and listens a port to accept new requests.

Usage:
<pre>
Usage of kraken-graph:
  -p=":9345": Address to listen requests
</pre>

Typical command line:

<pre>
kraken-graph -p :12345
</pre>

kraken-submit
------------

kraken-submit sends the metadata of the requests from a local repository to a kraken-graph.

It transfers the request encoded as [gobs](http://golang.org/pkg/encoding/gob/) through plain sockets.

Mandatory git config entries:

* _kraken.repoquest_: the public url of the repository.

Usage:
<pre>
Usage of kraken-submit: &lt;request ref&gt;+
  -p=":9345": Address of kraken-graph
</pre>

Typical command line:

<pre>
git config kraken.repoquest $(hostname):$(pwd)

kraken-submit refs/requests/<string>
git for-each-ref refs/requests/ --format='%(refname)' | xargs kraken-submit
</pre>
