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
		mu:      sync.Mutex{},
	}
}

func (um *UserMap) AddUser(username string, conn net.Conn) {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.userMap[username] = conn
}

func (um *UserMap) RemoveUser(username string) {
	um.mu.Lock()
	defer um.mu.Unlock()
	delete(um.userMap, username)
}

func (um *UserMap) GetUsernames() []string {
	um.mu.Lock()
	defer um.mu.Unlock()
	var members []string
	for username := range um.userMap {
		members = append(members, username)
	}
	return members
}

func (um *UserMap) GetUserConnection(username string) (net.Conn, bool) {
	um.mu.Lock()
	defer um.mu.Unlock()
	conn, ok := um.userMap[username]
	return conn, ok
}