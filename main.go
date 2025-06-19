package main

import (
	"fmt"

	"github.com/thisisadityapatel/just-servers/servers/echo"
	"github.com/thisisadityapatel/just-servers/servers/means_to_an_end"
	"github.com/thisisadityapatel/just-servers/servers/primetime"
	"github.com/thisisadityapatel/just-servers/servers/budget_chat"
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
		{name: "Means To An End Server", startFunc: means_to_an_end.Means_To_An_End_Server, port: "10002"},
		{name: "Budget Chat Server", startFunc: budget_chat.Budget_Chat_Server, port: "10003"},
	}

	// initiating servers in separate goroutines
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
