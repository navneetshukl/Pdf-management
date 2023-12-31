package main

import (
	"Pdf-Management/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	mux := chi.NewRouter()
	mux.Get("/", handlers.Home)
	mux.Get("/login", handlers.Login)
	mux.Post("/signup", handlers.Signup)
	mux.Post("/login", handlers.Authenticate)
	mux.Get("/upload-pdf",handlers.Pdf)
	mux.Post("/upload-pdf",handlers.StorePDF)
	mux.Get("/get-pdf",handlers.GetAllPdf)
	mux.Post("/pdf",handlers.HandlePDF)
	mux.Get("/options",handlers.Options)
	mux.Get("/single-pdf",handlers.GetLink)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("There is error in port number ",err)
	}

}
