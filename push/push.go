package kraken

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/json"
	"strings"
	"time"
)

func Push(git *Git, url, host, name string, refs ...string) ([]string, error) {

	now := time.Now()

	requests := make([]string, 0, len(refs))
	pushArgs := append(make([]string, 0, 2+len(refs)), "push", url)

	for _, ref := range refs {

		var hash string
		if output, err := git.Cmd("rev-parse", ref).Output(); err != nil {
			return nil, fmt.Errorf("rev-parse %s: %s", ref, err)
		} else {
			hash = strings.TrimSpace(string(output))
		}

		var fullref string
		if ref == "HEAD" {
			fullref = ref
		} else if output, err := git.Cmd("rev-parse", "--verify", "--symbolic-full-name", ref).Output(); err != nil {
			return nil, fmt.Errorf("rev-parse %s: %s", ref, err)
		} else {
			fullref = strings.TrimSpace(string(output))
			if fullref == "" {
				fullref = hash
			}
		}

		requestId, err := randomString()
		if err != nil {
			return nil, fmt.Errorf("random string: %s", err)
		}

		buffer := &bytes.Buffer{}
		buffer.WriteString("object ")
		buffer.WriteString(hash)
		buffer.WriteString("\ntype commit\ntag kraken-request-")
		buffer.WriteString(requestId)
		buffer.WriteString("\ntagger kraken-push <kraken-push@")
		buffer.WriteString(host)
		buffer.WriteString("> ")
		fmt.Fprint(buffer, now.Unix())
		_, zone := now.Zone()
		fmt.Fprintf(buffer, " %+03d00\n\n", zone/3600)
		encoder := json.NewEncoder(buffer)
		if err := encoder.Encode(RequestRef{name, fullref}); err != nil {
			return nil, fmt.Errorf("json encoding: %s", err)
		}

		cmd := git.Cmd("mktag")
		cmd.Stdin = buffer
		if output, err := cmd.Output(); err != nil {
			return nil, fmt.Errorf("git mktag: %s", err)
		} else {
			hash := strings.TrimSpace(string(output))
			request := "refs/requests/" + requestId
			requests = append(requests, request)
			pushArgs = append(pushArgs, hash+":"+request)
		}
	}

	if err := git.Cmd(pushArgs...).Run(); err != nil {
		return nil, fmt.Errorf("git push: %s", err)
	}

	return requests, nil
}

func randomString() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}
