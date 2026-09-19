package handlers

import (
    "net/http"

    "github.com/Ban-Coca/Linkwell/internal/store"
)

type LinksHandler struct {
	store store.LinkStore
}

func NewLinksHandler(s store.LinkStore) *LinksHandler {
    return &LinksHandler{store: s}
}

func (h *LinksHandler) Redirect(w http.ResponseWriter, r *http.Request) {
    // real logic goes here, using h.store
}