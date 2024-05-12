package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

const (
	FailedConnectionUpgrade       = "Was not able to upgrade connection"
	FailedDecodeCreateRoomRequest = "Was not able to decode create new room request"
	InvalidRequestProperties      = "Request has invalid properties"
	InvalidPathParameters         = "Request has invalid path parameters"
	FailedParsingUUID             = "Was not able to parse UUID"
	NotFoundWebsocketRoom         = "Could not find chat room"
)

type WebsocketManager struct {
	clients  WebsocketClientMap
	rooms    WebsocketRoomMap
	handlers map[string]WebsocketEventHandler

	validate *validator.Validate

	sync.RWMutex
}

func NewWebsocketManager(v *validator.Validate) *WebsocketManager {
	m := &WebsocketManager{
		clients:  make(WebsocketClientMap),
		rooms:    make(WebsocketRoomMap),
		handlers: make(map[string]WebsocketEventHandler),
		validate: v,
	}

	func() {
		m.handlers[EventSendMessage] = SendRoomMessage
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

func (m *WebsocketManager) connectNewClient(w http.ResponseWriter, r *http.Request) (*APISuccess, *APIError) {
	log.Println("Connecting new client")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		return nil, NewAPIError(err, FailedConnectionUpgrade, 500)
	}

	client := NewWebsocketClient(m, conn)

	m.addClient(client)
	return NewAPISuccess(200, ""), nil
}

func (m *WebsocketManager) createNewRoom(_ http.ResponseWriter, r *http.Request) (*APISuccess, *APIError) {
	var body CreateRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, NewAPIError(err, FailedDecodeCreateRoomRequest, 500)
	}

	if err := m.validate.Struct(body); err != nil {
		return nil, NewAPIError(err, InvalidRequestProperties, 400)
	}

	log.Printf("Creating new room %v", body.RoomName)

	room := NewWebsocketRoom(m, body.RoomName)
	m.addRoom(room)

	return NewAPISuccess(201, room), nil
}

func (m *WebsocketManager) getRoom(w http.ResponseWriter, r *http.Request) (*APISuccess, *APIError) {
	roomId := r.PathValue("roomId")
	log.Printf("Getting room with id: %v", roomId)

	if roomId == "" {
		return nil, NewAPIError(errors.New("path is missing value 'roomId'"), InvalidPathParameters, 400)
	}

	log.Printf("Parsing id: %v", roomId)

	if err := uuid.Validate(roomId); err != nil {
		return nil, NewAPIError(err, InvalidPathParameters, 400)
	}

	roomUUID, err := uuid.Parse(roomId)

	if err != nil {
		return nil, NewAPIError(err, InvalidPathParameters, 400)
	}

	log.Printf("Getting room with uuid: %v", roomUUID)

	if room, ok := m.rooms[roomUUID]; !ok {
		log.Printf("Did not find room with id: %v", roomUUID)
		return nil, NewAPIError(errors.New("no room found with passed id"), NotFoundWebsocketRoom, 404)
	} else {
		log.Printf("Found room: %v", room)
		return NewAPISuccess(200, room), nil
	}
}

func (m *WebsocketManager) addClient(c *WebsocketClient) {
	m.Lock()
	defer m.Unlock()

	m.clients[c.id] = c
}

func (m *WebsocketManager) removeClient(c *WebsocketClient) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[c.id]; ok {
		c.connection.Close()
		delete(m.clients, c.id)
	}
}

func (m *WebsocketManager) addRoom(r *WebsocketRoom) {
	m.Lock()
	defer m.Unlock()

	log.Printf("Adding room: %v to manager room map", r.Id)

	m.rooms[r.Id] = r
}

func SendRoomMessage(e WebsocketEvent, c *WebsocketClient) error {
	var event SendMessageEvent

	if err := json.Unmarshal(e.Payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	var broadcastMessage BroadcastMessageEvent
	broadcastMessage.Sent = time.Now()
	broadcastMessage.From = event.From
	broadcastMessage.Message = event.Message

	// data, err := json.Marshal(broadcastMessage)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal broadcast message: %v", err)
	// }

	// outgoingEvent := WebsocketEvent{
	// 	Payload: data,
	// 	Type:    EventBroadcastMessage,
	// }

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
