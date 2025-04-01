package primetime

import (
	"encoding/json"
	"fmt"
	"math/big"
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
		numberValue, numberOk := request["number"]

		if !methodOk || method != "isPrime" || !numberOk {
			fmt.Printf("[Primetime] Received malformed request: %+v\n", request)
			sendMalformedResponse(encoder)
			return
		}

		numberFloat, isFloat := numberValue.(float64)
		if !isFloat {
			fmt.Printf("[Primetime] Received non-number type: %+v\n", request)
			sendMalformedResponse(encoder)
			return
		}

		if numberFloat != float64(int64(numberFloat)) {
			fmt.Printf("[Primetime] Received non-integer number: %+v\n", request)
			sendMalformedResponse(encoder)
			return
		}

		var number big.Int
		number.SetInt64(int64(numberFloat))

		fmt.Printf("[Primetime] Received request: %+v\n", request)

		response := map[string]interface{}{
			"method": "isPrime",
			"prime":  isPrime(&number),
		}

		fmt.Printf("[Primetime] Outgoing response: %+v\n", response)

		if err := encoder.Encode(response); err != nil {
			fmt.Printf("[Primetime] Error sending response: %v\n", err)
			return
		}
	}
}

func isPrime(n *big.Int) bool {
	return n.ProbablyPrime(20)
}

func sendMalformedResponse(encoder *json.Encoder) {
	malformedResponse := map[string]interface{}{
		"error": "malformed request",
	}
	_ = encoder.Encode(malformedResponse)
}
