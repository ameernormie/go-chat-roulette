package main

import (
	"fmt"
	"io"
	"time"
)

var partner = make(chan io.ReadWriteCloser)

func main() {
	ticker := time.NewTicker(time.Millisecond * 250)
	boom := time.After(time.Second * 1)

	for {
		select {
		case <-ticker.C:
			fmt.Println("tick")
		case <-boom:
			fmt.Println("boom!")
			return
		}
	}
}

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	// Simultaneously trying to send the connection to partner channel and Also trying to receive from the partner channel and store it in the variable p
	// Does anybody want my connection and also can somebody give me a connection?
	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		// do some chat stuff
		chat(p, c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	go io.Copy(a, b)
	io.Copy(b, a)
}
