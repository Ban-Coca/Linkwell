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

func TestCreateRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{"url":`},
		{name: "missing URL", body: `{}`},
		{name: "empty URL", body: `{"url":""}`},
		{name: "relative URL", body: `{"url":"/relative"}`},
		{name: "missing host", body: `{"url":"https://"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewLinksHandler(store.NewMemoryStore())
			req := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(tt.body))
			resp := httptest.NewRecorder()

			h.Create(resp, req)

			if resp.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", resp.Code, resp.Body.String())
			}
		})
	}
}

func TestRedirectAndStatsReturnNotFound(t *testing.T) {
	h := NewLinksHandler(store.NewMemoryStore())

	for _, test := range []struct {
		name    string
		request *http.Request
		handle  func(http.ResponseWriter, *http.Request)
	}{
		{name: "redirect missing code", request: httptest.NewRequest(http.MethodGet, "/", nil), handle: h.Redirect},
		{name: "redirect unknown code", request: httptest.NewRequest(http.MethodGet, "/r/missing", nil), handle: h.Redirect},
		{name: "stats missing code", request: httptest.NewRequest(http.MethodGet, "/api/links//stats", nil), handle: h.Stats},
		{name: "stats unknown code", request: httptest.NewRequest(http.MethodGet, "/api/links/missing/stats", nil), handle: h.Stats},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			test.handle(resp, test.request)

			if resp.Code != http.StatusNotFound {
				t.Fatalf("expected 404, got %d: %s", resp.Code, resp.Body.String())
			}
		})
	}
}

func TestPathCodeSupportsRegisteredAndFallbackPaths(t *testing.T) {
	for _, test := range []struct {
		path string
		want string
	}{
		{path: "/abc123", want: "abc123"},
		{path: "/r/abc123", want: "abc123"},
		{path: "/api/links/abc123/stats", want: "abc123"},
		{path: "/", want: ""},
		{path: "/other/path", want: ""},
	} {
		t.Run(test.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			if got := pathCode(req); got != test.want {
				t.Fatalf("pathCode(%q) = %q, want %q", test.path, got, test.want)
			}
		})
	}
}
