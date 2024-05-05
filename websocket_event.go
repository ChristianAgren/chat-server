package main

import "encoding/json"

const (
	EventSendMessage = "send_message"
)

type WebsocketEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WebsocketEventHandler func(e WebsocketEvent, c *WebsocketClient) error

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
