package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string
	Details string
}

type APIError struct {
	Error      error
	Message    string
	StatusCode int
}

func NewAPIError(e error, m string, c int) *APIError {
	return &APIError{
		Error:      e,
		Message:    m,
		StatusCode: c,
	}
}

type APISuccess struct {
	StatusCode int
	Value      any
}

func NewAPISuccess(c int, value any) *APISuccess {
	log.Printf("Creating success response with value: %v", value)
	return &APISuccess{
		StatusCode: c,
		Value:      value,
	}
}

type apiHandler func(http.ResponseWriter, *http.Request) (*APISuccess, *APIError)

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if v, e := fn(w, r); e != nil {
		log.Printf("caught error: %v", e.Error)

		errMessage := ErrorResponse{
			Message: e.Message,
			Details: e.Error.Error(),
		}

		message, err := json.Marshal(errMessage)

		if err != nil {
			http.Error(w, "Server error, failed to format error", 500)
		}

		http.Error(w, string(message), e.StatusCode)
	} else {
		log.Printf("Responding with status: %v and value: %v", v.StatusCode, v.Value)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(v.StatusCode)
		json.NewEncoder(w).Encode(v.Value)
	}

}
