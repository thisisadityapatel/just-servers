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
		type Request struct {
			Method string    `json:"method"`
			Number BigNumber `json:"number"`
		}

		var request Request
		if err := decoder.Decode(&request); err != nil {
			sendMalformedResponse(encoder)
			return
		}

		if request.Method != "isPrime" {
			fmt.Printf("[Primetime] Received invalid method: %s\n", request.Method)
			sendMalformedResponse(encoder)
			return
		}

		if request.Number.IsFloat {
			fmt.Printf("[Primetime] Received float number, expected integer: %+v\n", request.Number)
			sendMalformedResponse(encoder)
			return
		}

		if request.Number.BigInt.Sign() < 0 {
			fmt.Printf("[Primetime] Received negative number: %+v\n", request.Number)
			sendMalformedResponse(encoder)
			return
		}

		fmt.Printf("[Primetime] Received request: method=%s, number=%s\n",
			request.Method, request.Number.BigInt.String())

		response := map[string]interface{}{
			"method": "isPrime",
			"prime":  isPrime(request.Number.BigInt),
		}

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

type BigNumber struct {
	BigInt  *big.Int
	IsFloat bool
}

func (n *BigNumber) UnmarshalJSON(data []byte) error {
	numStr := string(data)

	floatValue := new(big.Float)
	if _, ok := floatValue.SetString(numStr); ok {
		if floatValue.IsInt() {
			n.BigInt = new(big.Int)
			floatValue.Int(n.BigInt)
			n.IsFloat = false
		} else {
			n.IsFloat = true
		}
	} else {
		return fmt.Errorf("invalid number format: %s", numStr)
	}
	return nil
}
