package main

import (
	"flag"
	"fmt"
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
	args := flag.Args()

	if *version {
		ShowVersion()
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := krak(wd, args...); err != nil {
		log.Fatal(err)
	}
}
