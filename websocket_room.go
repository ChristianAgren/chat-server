package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type WebsocketRoomMap map[uuid.UUID]*WebsocketRoom

type WebsocketRoom struct {
	Id      uuid.UUID          `json:"id"`
	Name    string             `json:"name"`
	Manager *WebsocketManager  `json:"manager"`
	Clients WebsocketClientMap `json:"clients"`

	sync.RWMutex
}

func NewWebsocketRoom(m *WebsocketManager, n string) *WebsocketRoom {
	room := &WebsocketRoom{
		Id:      uuid.New(),
		Name:    n,
		Manager: m,
		Clients: make(WebsocketClientMap),
	}

	log.Printf("Creating new room: %v", room.Name)

	return room
}

// func (r *WebsocketRoom) addClient(c *WebsocketClient) {
// 	r.Lock()
// 	defer r.Unlock()

// 	r.clients[c.id] = c
// }

func (r *WebsocketRoom) SendMessage(message *BroadcastMessageEvent) error {

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	outgoingEvent := WebsocketEvent{
		Payload: data,
		Type:    EventBroadcastMessage,
	}

	for id := range r.Clients {
		r.Clients[id].egress <- outgoingEvent
	}

	return nil
}
