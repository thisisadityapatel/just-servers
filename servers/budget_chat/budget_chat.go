package budget_chat

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"unicode"

	"github.com/thisisadityapatel/just-servers/utilities"
)

func Budget_Chat_Server(Port string) error {
	budgetChatServer := utilities.NewTcpServer(Port)
	listener, err := utilities.GetListener(*budgetChatServer)
	if err != nil {
		return err
	}
	defer listener.Close()

	userMap := NewUserMap()

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		fmt.Printf("[BudgetChat] New connection established from %s\n", conn.RemoteAddr())
		wg.Add(1)
		go handleConnection(conn, &wg, userMap)
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup, userMap *UserMap) {
	defer conn.Close()
	defer wg.Done()

	// send welcome message and prompt for username
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	scanner := bufio.NewScanner(conn)

	nickname, err := joinRequest(scanner, userMap)

	if err != nil || nickname == "" {
		log.Printf("validation error: %v", err)
		return
	}

	getServerUserPresentce(conn, userMap)

	userMap.AddUser(nickname, conn)

	defer leaveRequest(nickname, userMap)

	Broadcast(nickname, "* "+nickname+" joined the room", userMap)
	log.Printf("%s joined the room", nickname)

	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())

		if len(message) > 1000 {
			conn.Write([]byte("message is too long. Re-send the message\n"))
			continue
		}

		Broadcast(nickname, "["+nickname+"] "+message, userMap)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %v", err)
	}

}

func joinRequest(scanner *bufio.Scanner, userMap *UserMap) (string, error) {
	ok := scanner.Scan()

	if !ok {
		return "", scanner.Err()
	}

	inputUsername := strings.TrimSpace(scanner.Text())
	log.Print("[BudgetChat] Received name: ", inputUsername)

	if len(inputUsername) < 1 || len(inputUsername) > 18 {
		return "", errors.New("Length of the username is less than 2 or greater than 19")
	}

	for _, character := range inputUsername {
		if !unicode.IsLetter(character) && !unicode.IsDigit(character) {
			return "", errors.New("Invalid characters")
		}
	}

	if _, ok := userMap.GetUserConnection(inputUsername); ok {
		return "", errors.New("Username already taken in the chat server")
	}

	return inputUsername, nil
}

func getServerUserPresentce(conn net.Conn, userMap *UserMap) {
	roomMembers := userMap.GetUsernames()
	nicknames := strings.Join(roomMembers, ", ")
	conn.Write([]byte("* The room contains: " + nicknames + "\n"))
}

func Broadcast(sender string, message string, userMap *UserMap) {
	for _, username := range userMap.GetUsernames() {
		if username != sender {
			conn, _ := userMap.GetUserConnection(username)
			conn.Write([]byte(message + "\n"))
		}
	}
}

func leaveRequest(username string, userMap *UserMap) {
	userMap.RemoveUser(username)
	Broadcast(username, "* "+username+" has left the room", userMap)
	log.Printf("%s left the room", username)
}