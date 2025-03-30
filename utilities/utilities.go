package utilities

import (
	"fmt"
	"net"
	"sync"
)

const Host string = "0.0.0.0"
const TcpType string = "tcp"

type ServerConfig struct {
	Host string
	Port string
	Type string
}

type Server interface {
	Start(wg *sync.WaitGroup) error
}

type TcpServer struct {
	config ServerConfig
}

// creating a server listner for a given server config
func GetListener(server TcpServer) (net.Listener, error) {
	listener, err := net.Listen(server.config.Type, server.config.Host+":"+server.config.Port)
	if err != nil {
		return nil, fmt.Errorf("error creating listener: %v", err)
	}
	fmt.Printf("Server listening on %s:%s\n", server.config.Host, server.config.Port)
	return listener, nil
}

func NewTcpServer(port string) *TcpServer {
	return &TcpServer{
		config: ServerConfig{
			Host: Host,
			Port: port,
			Type: TcpType,
		},
	}
}
