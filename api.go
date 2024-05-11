package main

import (
	"log"
	"net/http"
	// "github.com/gorilla/mux"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{addr: addr}
}

func (s *APIServer) Serve() {
	router := http.NewServeMux()

	websocketManager := NewWebsocketManager()

	setupRoutes(router, websocketManager)

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	log.Printf("Starting the API server on port%v", s.addr)
	log.Fatal(server.ListenAndServe())
}

func setupRoutes(r *http.ServeMux, wsManager *WebsocketManager) {
	r.HandleFunc("GET /{$}", handleHomePage)
	r.HandleFunc("GET /ws", wsManager.connect)
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	println("Welcome to Server ", r.RequestURI)
}
