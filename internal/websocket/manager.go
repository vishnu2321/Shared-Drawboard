package websocket

import (
	"sync"
)

type Manager struct {
	ClientList map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Mu         sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		ClientList: make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Mu:         sync.Mutex{},
	}
}

func (m *Manager) Run() {
	for {
		select {
		//register for a new client
		case client := <-m.Register:
			m.Mu.Lock()
			m.ClientList[client.ID] = client
			m.Mu.Unlock()
		//unregister client
		case client := <-m.Unregister:
			m.Mu.Lock()
			if _, ok := m.ClientList[client.ID]; ok {
				delete(m.ClientList, client.ID)
				close(client.Send)
			}
			m.Mu.Unlock()
		//manager sends message to all clients
		case message := <-m.Broadcast:
			m.Mu.Lock()
			for _, client := range m.ClientList {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.ClientList, client.ID)
				}
			}
			m.Mu.Unlock()
		}
	}
}
