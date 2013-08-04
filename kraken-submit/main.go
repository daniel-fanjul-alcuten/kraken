package main

import (
	"flag"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
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
	address := flag.String("p", ":9345", "Address of kraken-graph")
	flag.Parse()
	args := flag.Args()

	git := NewGit("")
	repoquest := getConfig(git, "kraken.repoquest")

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	if err := submit(git, repoquest, conn, args...); err != nil {
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
