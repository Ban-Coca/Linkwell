package main

import (
	"log"
	"net/http"

	"github.com/Ban-Coca/Linkwell/internal/handlers"
	"github.com/Ban-Coca/Linkwell/internal/store"
)

func main() {
	store := store.NewMemoryStore()
	linksHandler := handlers.NewLinksHandler(store)

	http.HandleFunc("/r/", linksHandler.Redirect)
	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
