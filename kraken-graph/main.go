package main

import (
	"flag"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/queue"
	"log"
	"net"
	"os"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	version := flag.Bool("version", false, "Shows version")
	address := flag.String("p", ":9345", "Address to listen requests")
	flag.Parse()

	if *version {
		ShowVersion()
	}

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

	queue := NewQueue(1024)
	input := queue.GetInput()
	go func() {
		g := newGraph()
		for request := range requests {
			log.Printf("%#v", request)
			jobs := g.addRequest(request)
			for _, job := range jobs {
				log.Printf("%#v", job)
				input <- job
			}
		}
	}()

	if err := <-errs; err != nil {
		log.Fatal(err)
	}
}
