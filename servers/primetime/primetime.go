package primetime

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"sync"

	"github.com/thisisadityapatel/just-servers/utilities"
)

func PrimeServer(Port string) error {
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
			fmt.Printf("[Primetime] Error accepting connection: %v\n", err)
			continue
		}
		fmt.Printf("[Primetime] New connection established from %s\n", conn.RemoteAddr())
		wg.Add(1)
		go handleConnection(conn, &wg)
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var request map[string]interface{}
		if err := decoder.Decode(&request); err != nil {
			sendMalformedResponse(encoder)
			return
		}

		method, methodOk := request["method"].(string)
		number, numberOk := request["number"].(float64)

		if !methodOk || method != "isPrime" || !numberOk {
			fmt.Printf("[Primetime] Received malformed request: %+v\n", request)
			sendMalformedResponse(encoder)
			return
		}

		fmt.Printf("[Primetime] Received request: %+v\n", request)

		response := map[string]interface{}{
			"method": "isPrime",
			"prime":  isPrime(number),
		}

		if err := encoder.Encode(response); err != nil {
			fmt.Printf("[Primetime] Error sending response: %v\n", err)
			return
		}
	}
}

func sendMalformedResponse(encoder *json.Encoder) {
	malformedResponse := map[string]interface{}{
		"error": "malformed request",
	}
	_ = encoder.Encode(malformedResponse)
}

func isPrime(number float64) bool {
	if number <= 1 || number != math.Floor(number) {
		return false
	}
	n := int(number)
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
