package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Ban-Coca/Linkwell/internal/shortener"
	"github.com/Ban-Coca/Linkwell/internal/store"
)

type LinksHandler struct {
	store store.LinkStore
}

func NewLinksHandler(s store.LinkStore) *LinksHandler {
	return &LinksHandler{store: s}
}

func isValidURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed == nil {
		return false
	}
	return parsed.Scheme != "" && parsed.Host != ""
}

func pathCode(r *http.Request) string {
	if code := r.PathValue("code"); code != "" {
		return code
	}

	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) == 2 && parts[0] == "r" {
		return parts[1]
	}
	if len(parts) == 4 && parts[0] == "api" && parts[1] == "links" && parts[3] == "stats" {
		return parts[2]
	}
	return ""
}

func (h *LinksHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := pathCode(r)
	if code == "" {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	link, err := h.store.Get(code)
	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	if err := h.store.RecordClick(code, store.ClickEvent{
		Timestamp: time.Now(),
		Referrer:  r.Referer(),
		UserAgent: r.UserAgent(),
	}); err != nil {
		log.Printf("record click for %s: %v", code, err)
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func (h *LinksHandler) Stats(w http.ResponseWriter, r *http.Request) {
	code := pathCode(r)
	if code == "" {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	stats, err := h.store.GetStats(code)
	if err != nil {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *LinksHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" || !isValidURL(req.URL) {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	code, err := shortener.GenerateCode()
	if err != nil {
		http.Error(w, "could not generate link code", http.StatusInternalServerError)
		return
	}

	for {
		exists, err := h.store.Exists(code)
		if err != nil {
			http.Error(w, "could not check link code", http.StatusInternalServerError)
			return
		}
		if !exists {
			break
		}

		code, err = shortener.GenerateCode()
		if err != nil {
			http.Error(w, "could not generate link code", http.StatusInternalServerError)
			return
		}
	}

	link := store.Link{
		Code:        code,
		OriginalURL: req.URL,
		CreatedAt:   time.Now(),
	}

	if err := h.store.Create(link); err != nil {
		http.Error(w, "could not create link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}