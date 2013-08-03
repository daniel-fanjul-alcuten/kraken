package main

import (
	"flag"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken"
	"log"
	"os"
	"strings"
)

func main() {

	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: <local ref>+\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	git := NewGit("")
	url := getConfig(git, "remote.kraken.url")
	repository := getConfig(git, "kraken.repository")

	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := push(git, url, host, repository, args...); err != nil {
		log.Fatal(err)
	}
}

func getConfig(git *Git, name string) string {
	output, err := git.Cmd("config", name).Output()
	if err != nil {
		log.Fatalf("git config %s: not found", name)
	}
	return strings.TrimSpace(string(output))
}
