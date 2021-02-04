package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const listenAddr = "localhost:4000"

var partner = make(chan io.ReadWriteCloser)

func main() {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go match(conn)
	}
}

func match(conn io.ReadWriteCloser) {
	fmt.Fprint(conn, "Waiting for a partner...")
	// Simultaneously trying to send the connection to partner channel and Also trying to receive from the partner channel and store it in the variable p
	// Does anybody want my connection and also can somebody give me a connection?
	select {
	case partner <- conn:
		// now handled by the other goroutine
	case p := <-partner:
		// do some chat stuff
		chat(p, conn)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}
