package main

import (
	"bufio"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	"net"
)

func listen(listener net.Listener, requests chan<- Request) error {

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func(fconn net.Conn) {
			defer fconn.Close()
			reader := bufio.NewReader(fconn)
			decoder := gob.NewDecoder(reader)
			var request Request
			for {
				if err := decoder.Decode(&request); err != nil {
					break
				}
				requests <- request
			}
		}(conn)
	}
}
