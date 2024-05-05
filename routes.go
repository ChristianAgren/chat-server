package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RoutesService struct {
	m *WebsocketManager
}

func NewRoutesService(manager *WebsocketManager) *RoutesService {
	return &RoutesService{
		m: manager,
	}
}

func (s *RoutesService) RegisterRoutes(r *mux.Router) {
	println("Setting up routes")
	r.HandleFunc("/", s.handleHomePage).Methods("GET")
	r.HandleFunc("/ws", s.m.handleWebSocketConnection).Methods("GET")
}

func (s *RoutesService) handleHomePage(w http.ResponseWriter, r *http.Request) {
	println("Welcome to websockets")
}
