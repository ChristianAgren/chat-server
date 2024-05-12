package main

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	EventSendMessage      = "send_message"
	EventBroadcastMessage = "broadcast_message"
)

type WebsocketEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WebsocketEventHandler func(e WebsocketEvent, c *WebsocketClient) error

type SendMessageEvent struct {
	Message string    `json:"message"`
	From    string    `json:"from"`
	RoomId  uuid.UUID `json:"roomId"`
}

type BroadcastMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type CreateRoomRequest struct {
	RoomName string `json:"roomName" validate:"required"`
}

type GetRoomRequest struct {
	RoomId string `json:"roomId" validate:"required"`
}
