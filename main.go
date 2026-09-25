package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/Ban-Coca/Linkwell/internal/handlers"
	"github.com/Ban-Coca/Linkwell/internal/store"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	store := store.NewMemoryStore()
	linksHandler := handlers.NewLinksHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/links", linksHandler.Create)
	mux.HandleFunc("GET /api/links/{code}/stats", linksHandler.Stats)
	mux.HandleFunc("GET /{code}", linksHandler.Redirect)
	mux.HandleFunc("GET /r/{code}", linksHandler.Redirect)

	staticRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticRoot)))

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
