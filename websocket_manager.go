package main

import (
	"errors"
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
	clients  WebsocketClientMap
	handlers map[string]WebsocketEventHandler

	sync.RWMutex
}

func NewWebsocketManager() *WebsocketManager {
	m := &WebsocketManager{
		clients:  make(WebsocketClientMap),
		handlers: make(map[string]WebsocketEventHandler),
	}

	func() {
		m.handlers[EventSendMessage] = SendMessage
	}()

	return m
}

func (m *WebsocketManager) routeEvent(e WebsocketEvent, c *WebsocketClient) error {
	if handler, ok := m.handlers[e.Type]; ok {
		if err := handler(e, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
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

func SendMessage(e WebsocketEvent, c *WebsocketClient) error {
	log.Println(e)
	c.egress <- e
	return nil
}
