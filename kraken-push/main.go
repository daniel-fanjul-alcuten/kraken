package main

import (
	"flag"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/push"
	. "github.com/daniel-fanjul-alcuten/kraken/version"
	"log"
	"os"
)

func main() {

	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: <local ref>+\n", os.Args[0])
		flag.PrintDefaults()
	}
	version := flag.Bool("version", false, "Shows version")
	flag.Parse()

	if *version {
		ShowVersion()
	}

	git := NewGit("")
	url := getConfig(git, "remote.kraken.url")
	repository := getConfig(git, "kraken.repository")

	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()
	if _, err := Push(git, url, host, repository, args...); err != nil {
		log.Fatal(err)
	}
}

func getConfig(git *Git, name string) string {
	value, err := git.String(nil, "config", name)
	if err != nil {
		log.Fatal(err)
	}
	return value
}
