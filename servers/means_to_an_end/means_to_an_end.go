package means_to_an_end

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/thisisadityapatel/just-servers/utilities"
)

func Means_To_An_End_Server(Port string) error {
	echoServer := utilities.NewTcpServer(Port)
	listener, err := utilities.GetListener(*echoServer)
	if err != nil {
		return err
	}
	defer listener.Close()

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		// increment wait group and handle connection in goroutine
		fmt.Printf("[MeansToAnEnd] New connection established from %s\n", conn.RemoteAddr())
		wg.Add(1)
		go handleConnection(conn, &wg)
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	// Setup in-memory store for this client connection
	store := NewInMemoryStore()
	request := make(Request, 9)

	for {
		err := binary.Read(conn, binary.BigEndian, request)
		fmt.Printf("[MeansToAnEnd] Received binary data: %v\n", request)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("[MeansToAnEnd] Connection closed by client: %v\n", conn.RemoteAddr())
				return
			} else {
				fmt.Printf("[MeansToAnEnd] Error reading data: %v\n", err)
				break
			}
		}

		operation, input1, input2 := request.Decode()
		fmt.Printf("[MeansToAnEnd] Decoded operation: %c, input1: %d, input2: %d\n", operation, input1, input2)

		handleRequest(conn, operation, input1, input2, store)
	}
}

func handleRequest(conn net.Conn, operation rune, input1 int32, input2 int32, store *InMemoryStore) {
	switch operation {
	case 'I':
		// Handle insert operation
		handleInsert(store, input1, input2)
	case 'Q':
		// Handle query operation
		handleQuery(conn, store, input1, input2)
	default:
		fmt.Printf("[MeansToAnEnd] Unknown operation: %c\n", operation)
		conn.Write([]byte("Unknown operation\n"))
	}
}

func handleInsert(store *InMemoryStore, timestamp int32, price int32) {
	store.Insert(timestamp, price)
	fmt.Printf("[MeansToAnEnd] Insert operation: %d, %d\n", timestamp, price)
}

func handleQuery(conn net.Conn, store *InMemoryStore, minTime int32, maxTime int32) {
	average := store.Query(minTime, maxTime)
	response := make([]byte, 4)
	binary.BigEndian.PutUint32(response, uint32(average))
	fmt.Printf("[MeansToAnEnd] Query operation: minTime: %d, maxTime: %d, average: %d\n", minTime, maxTime, average)
	conn.Write(response)
}
