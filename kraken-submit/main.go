package main

import (
	"flag"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/submit"
	. "github.com/daniel-fanjul-alcuten/kraken/version"
	"log"
	"net"
	"os"
)

func main() {

	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: <request ref>+\n", os.Args[0])
		flag.PrintDefaults()
	}
	version := flag.Bool("version", false, "Shows version")
	address := flag.String("r", ":9345", "Address where kraken-graph listens requests")
	flag.Parse()

	if *version {
		ShowVersion()
	}

	git := NewGit("")
	repoquest := getConfig(git, "kraken.repoquest")

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()
	if err := Submit(git, repoquest, conn, args...); err != nil {
		log.Fatal(err)
	}
}

func getConfig(git *Git, name string) string {
	value, err := git.Config(name)
	if err != nil {
		log.Fatal(err)
	}
	return value
}
