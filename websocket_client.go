package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type WebsocketClientMap map[*WebsocketClient]bool

type WebsocketClient struct {
	connection *websocket.Conn
	manager    *WebsocketManager
	egress     chan WebsocketEvent
}

func NewWebsocketClient(c *websocket.Conn, m *WebsocketManager) *WebsocketClient {
	client := &WebsocketClient{
		connection: c,
		manager:    m,
		egress:     make(chan WebsocketEvent),
	}

	go client.setupHeartBeat()
	go client.readMessages()
	go client.writeMessages()

	return client
}

func (c *WebsocketClient) setupHeartBeat() {
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("was not able to reset read deadline: %v", err)
		return
	}

	c.connection.SetPongHandler(func(msg string) error {
		log.Printf("pong %v", msg)
		return c.connection.SetReadDeadline(time.Now().Add(pongWait))
	})
}

func (c *WebsocketClient) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request WebsocketEvent
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("was not able to unmarshal payload: %v into request with err: %v", payload, err)
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Printf("was not able to route event: %v", err)
		}

	}
}

func (c *WebsocketClient) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed: %v", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("was not able to marshal message: %v", err)
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)
			}

		case <-ticker.C:
			log.Println("ping")
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("failed to send ping message: %v", err)
				return
			}
		}
	}
}
