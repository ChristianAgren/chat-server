package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WebsocketManager struct {
	clients WebsocketClientMap

	sync.RWMutex
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		clients: make(WebsocketClientMap),
	}
}

func (m *WebsocketManager) handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := NewWebsocketClient(conn, m)

	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *WebsocketManager) addClient(c *WebsocketClient) {
	m.Lock()
	defer m.Unlock()

	m.clients[c] = true
}

func (m *WebsocketManager) removeClient(c *WebsocketClient) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[c]; ok {
		c.connection.Close()
		delete(m.clients, c)
	}
}
