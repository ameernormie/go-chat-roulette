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
	go io.Copy(a, b)
	io.Copy(b, a)
}
