package main

import (
	"Pdf-Management/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	mux := chi.NewRouter()
	mux.Get("/",handlers.Home)
	mux.Get("/login",handlers.Login)
	mux.Post("/signup",handlers.Signup)
	mux.Post("/login",handlers.Authenticate)


	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("There is error in port number")
	}

}
