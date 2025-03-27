package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	HOST = "0.0.0.0"
	PORT = "10000"
	TYPE = "tcp"
)

func main() {
	listener, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Printf("Error creating listener: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s:%s\n", HOST, PORT)

	// wait group to keep track of concurrent connection processing in goroutine
	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		// increment wait group and handle connection in goroutine
		wg.Add(1)
		go handleConnection(conn, &wg)
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	// ensure connection is closed and wait group is decremented when done
	defer conn.Close()
	defer wg.Done()

	// use io.Copy to echo data back to the connection
	_, err := io.Copy(conn, conn)
	if err != nil {
		fmt.Printf("Error during copy: %v\n", err)
		return
	}
}
