package main

import (
	"flag"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/graph"
	. "github.com/daniel-fanjul-alcuten/kraken/queue"
	. "github.com/daniel-fanjul-alcuten/kraken/version"
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
	reqAddress := flag.String("r", ":9345", "Address to listen requests")
	jobAddress := flag.String("j", ":9346", "Address to listen to send jobs")
	flag.Parse()

	if *version {
		ShowVersion()
	}

	reqListener, err := net.Listen("tcp", *reqAddress)
	if err != nil {
		log.Fatal(err)
	}

	jobListener, err := net.Listen("tcp", *jobAddress)
	if err != nil {
		log.Fatal(err)
	}

	requests := make(chan Request, 1024)
	errs := make(chan error, 1)
	go func() {
		errs <- listen(reqListener, &requestHandler{requests})
	}()

	queue := NewQueue(10 * 1024)
	results := make(chan result, 64)
	go func() {
		errs <- listen(jobListener, &jobHandler{queue, results})
	}()

	input := queue.GetInput()
	go func() {
		g := NewGraph()
		for request := range requests {
			log.Printf("New Request: %#v", request)
			jobs := g.AddRequest(request)
			for _, job := range jobs {
				log.Printf("New Job: %#v", job)
				input <- job
			}
		}
	}()

	go func() {
		for result := range results {
			if result.ack {
				log.Printf("New Result: %v for %#v", result.result.Success, result.job)
			} else {
				log.Printf("Lost Result: %#v", result.job)
			}
		}
	}()

	if err := <-errs; err != nil {
		log.Fatal(err)
	}
}
