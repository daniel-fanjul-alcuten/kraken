package main

import (
	"flag"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/version"
	. "github.com/daniel-fanjul-alcuten/kraken/work"
	"log"
	"net"
	"os"
)

func main() {

	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	version := flag.Bool("version", false, "Shows version")
	address := flag.String("j", ":9346", "Address where kraken-graph sends jobs")
	flag.Parse()

	if *version {
		ShowVersion()
	}

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stderr, "", 0)
	if err := Work(conn, conn, logger); err != nil {
		log.Fatal(err)
	}
}
