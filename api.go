package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	// "github.com/gorilla/mux"
)

type APIServer struct {
	addr     string
	validate *validator.Validate
}

func NewAPIServer(addr string, v *validator.Validate) *APIServer {
	return &APIServer{addr: addr, validate: v}
}

func (s *APIServer) Serve() {
	router := http.NewServeMux()

	websocketManager := NewWebsocketManager(s.validate)

	setupRoutes(router, websocketManager)

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	log.Printf("Starting the API server on port%v", s.addr)
	log.Fatal(server.ListenAndServe())
}

func setupRoutes(r *http.ServeMux, wsManager *WebsocketManager) {
	r.Handle("GET /{$}", apiHandler(handleHomePage))
	r.Handle("GET /ws", apiHandler(wsManager.connectNewClient))
	r.Handle("GET /room/{roomId}", apiHandler(wsManager.getRoom))
	r.Handle("POST /room", apiHandler(wsManager.createNewRoom))
}

func handleHomePage(w http.ResponseWriter, r *http.Request) (*APISuccess, *APIError) {
	println("Welcome to Server ", r.RequestURI)
	return NewAPISuccess(200, ""), nil
}
