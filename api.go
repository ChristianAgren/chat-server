package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{addr: addr}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	websocketManager := NewWebsocketManager()
	routesService := NewRoutesService(websocketManager)

	routesService.RegisterRoutes(subrouter)

	log.Printf("Starting the API server on port%v", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, subrouter))
}
