package main

import (
	"io"
	"net"
)

type handler interface {
	handle(reader io.Reader, writer io.Writer) error
}

func listen(listener net.Listener, handler handler) error {

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			defer conn.Close()
			handler.handle(conn, conn)
		}()
	}
}
