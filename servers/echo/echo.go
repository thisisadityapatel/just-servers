package echo

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/thisisadityapatel/just-servers/utilities"
)

func EchoServer(Port string) error {
	echoServer := utilities.NewTcpServer(Port)
	listener, err := utilities.GetListener(*echoServer)
	if err != nil {
		return err
	}
	defer listener.Close()

	// wait group for tracking concurrent goroutines
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
	defer conn.Close()
	defer wg.Done()

	_, err := io.Copy(conn, conn)
	if err != nil {
		fmt.Printf("Error during copy: %v\n", err)
		return
	}
}
