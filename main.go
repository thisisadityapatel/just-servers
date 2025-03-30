package main

import (
	"fmt"

	"github.com/thisisadityapatel/just-servers/servers/echo"
	"github.com/thisisadityapatel/just-servers/servers/primetime"
)

func main() {
	// server configurations with unique ports
	servers := []struct {
		name      string
		startFunc func(string) error
		port      string
	}{
		{name: "Echo Server", startFunc: echo.EchoServer, port: "10000"},
		{name: "Prime Server", startFunc: primetime.PrimeServer, port: "10001"},
	}

	// starting each server in a goroutine
	for _, server := range servers {
		go func(name, port string, startFunc func(string) error) {
			fmt.Printf("Starting %s on port %s...\n", name, port)
			if err := startFunc(port); err != nil {
				fmt.Printf("Error starting %s: %v\n", name, err)
			}
		}(server.name, server.port, server.startFunc)
	}

	// keeping the main run indefinitely
	select {}
}
