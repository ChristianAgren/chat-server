package main

import "github.com/go-playground/validator/v10"

func main() {
	validate := validator.New(validator.WithRequiredStructEnabled())
	api := NewAPIServer(":8080", validate)
	api.Serve()
}
