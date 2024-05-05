package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type WebsocketClientMap map[*WebsocketClient]bool

type WebsocketClient struct {
	connection *websocket.Conn
	manager    *WebsocketManager
}

func NewWebsocketClient(c *websocket.Conn, m *WebsocketManager) *WebsocketClient {
	return &WebsocketClient{
		connection: c,
		manager:    m,
	}
}

func (c *WebsocketClient) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		log.Println(messageType)
		log.Println(string(payload))
	}
}

func (c *WebsocketClient) writeMessages() {

}
