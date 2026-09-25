package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ban-Coca/Linkwell/internal/store"
)

func TestCreateAndStatsRoundTrip(t *testing.T) {
	mem := store.NewMemoryStore()
	h := NewLinksHandler(mem)

	createReq := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(`{"url":"https://example.com/long"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	h.Create(createResp, createReq)

	if createResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", createResp.Code, createResp.Body.String())
	}

	var created struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}

	r := httptest.NewRequest(http.MethodGet, "/r/"+created.Code, nil)
	r.Header.Set("Referer", "https://example.com/ref")
	r.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	h.Redirect(w, r)

	if w.Code != http.StatusFound {
		t.Fatalf("expected redirect status 302, got %d", w.Code)
	}

	statsReq := httptest.NewRequest(http.MethodGet, "/api/links/"+created.Code+"/stats", nil)
	statsResp := httptest.NewRecorder()
	h.Stats(statsResp, statsReq)

	if statsResp.Code != http.StatusOK {
		t.Fatalf("expected stats 200, got %d: %s", statsResp.Code, statsResp.Body.String())
	}

	body := statsResp.Body.String()
	if !strings.Contains(body, "\"totalClicks\":1") {
		t.Fatalf("expected one click in stats, got %s", body)
	}
}
