package main

import (
	"flag"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	"log"
	"net"
	"os"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	address := flag.String("p", ":9345", "Address to listen requests")
	flag.Parse()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	// g := newGraph()
	requests := make(chan Request, 1024)
	errs := make(chan error, 1)
	go func() {
		errs <- listen(listener, requests)
	}()

	go func() {
		g := newGraph()
		for request := range requests {
			log.Printf("%#v", request)
			g.addRequest(request)
		}
	}()

	if err := <-errs; err != nil {
		log.Fatal(err)
	}
}
