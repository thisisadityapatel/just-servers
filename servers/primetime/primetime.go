package primetime

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
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

type Request struct {
	Method *string    `json:"method"`
	Number *BigNumber `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func (req *Request) validFields() bool {
	if req.Method == nil {
		return false
	}

	if *req.Method != "isPrime" {
		return false
	}

	if req.Number == nil {
		return false
	}

	return true
}

func handleConnection(c net.Conn, wg *sync.WaitGroup) {
	defer c.Close()
	defer wg.Done()

	reader := bufio.NewReader(c)

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			log.Print(err.Error())
			c.Write([]byte(err.Error()))
			break
		}

		log.Printf("Recieved: %v", message)

		var request Request
		err = json.Unmarshal([]byte(message), &request)

		// handle validation
		if err != nil {
			log.Print(err.Error())
			c.Write([]byte(err.Error()))
			break
		}

		if !request.validFields() {
			log.Print("Invalid Fields")
			c.Write([]byte("Invalid Fields"))
			break
		}

		// convert the string to bigInt
		if request.Number.IsFloat {
			res := Response{Method: "isPrime", Prime: false}
			handleResponse(c, res)
			break
		}

		// handle response
		res := Response{Method: "isPrime", Prime: isPrime(*request.Number.BigInt)}
		handleResponse(c, res)
	}
}

func handleResponse(c net.Conn, res Response) {
	resJson, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Response: ", string(resJson)+"\n")
	c.Write([]byte(string(resJson) + "\n"))
}

func isPrime(n big.Int) bool {
	return n.ProbablyPrime(10)
}

type BigNumber struct {
	BigInt  *big.Int
	IsFloat bool
}

func (n *BigNumber) UnmarshalJSON(data []byte) error {
	// Try to unmarshal into a string first
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
