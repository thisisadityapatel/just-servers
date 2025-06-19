package budget_chat

import (
	"net"
	"sync"
)

type UserMap struct {
	userMap map[string]net.Conn
	mu      sync.Mutex
}

func NewUserMap() *UserMap {
	return &UserMap{
		userMap: make(map[string]net.Conn),
	}
}

func (um *UserMap) AddUser(username string, conn net.Conn) {
	um.mu.lock()
	defer um.mu.unlock()
	um.userMap[username] = conn
}

func (um *UserMap) RemoveUser(username string) {
	um.mu.lock()
	defer um.mu.unlock()
	delete(um.userMap, username)
}

func (um *UserMap) GetUsernames (net.Conn) ([]string) {
	um.mu.lock()
	defer um.mu.unlock()
	var members []string
	for username := range um.userMap {
		members = append(members, username)
	}
	return members
}

func (um *UserMap) GetUserConnection(username string) (net.Conn, bool) {
	um.mu.lock()
	defer um.mu.unlock()
	conn, ok := um.userMap[username]
	return conn, ok
}