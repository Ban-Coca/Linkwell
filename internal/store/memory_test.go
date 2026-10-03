package store

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestMemoryStoreCreateGetAndExists(t *testing.T) {
	s := NewMemoryStore()
	link := Link{Code: "abc123", OriginalURL: "https://example.com", CreatedAt: time.Now()}

	if exists, err := s.Exists(link.Code); err != nil || exists {
		t.Fatalf("initial Exists() = (%t, %v), want (false, nil)", exists, err)
	}
	if err := s.Create(link); err != nil {
		t.Fatalf("Create() returned an error: %v", err)
	}
	if err := s.Create(link); !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("duplicate Create() error = %v, want ErrAlreadyExists", err)
	}
	if exists, err := s.Exists(link.Code); err != nil || !exists {
		t.Fatalf("Exists() = (%t, %v), want (true, nil)", exists, err)
	}

	got, err := s.Get(link.Code)
	if err != nil {
		t.Fatalf("Get() returned an error: %v", err)
	}
	if !reflect.DeepEqual(got, Link{Code: link.Code, OriginalURL: link.OriginalURL, CreatedAt: link.CreatedAt, Clicks: []ClickEvent{}}) {
		t.Fatalf("Get() = %#v, want initialized link", got)
	}
}

func TestMemoryStoreMissingLinks(t *testing.T) {
	s := NewMemoryStore()

	if _, err := s.Get("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() error = %v, want ErrNotFound", err)
	}
	if err := s.RecordClick("missing", ClickEvent{}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("RecordClick() error = %v, want ErrNotFound", err)
	}
	if _, err := s.GetStats("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetStats() error = %v, want ErrNotFound", err)
	}
}

func TestMemoryStoreStatsAggregatesAndSorts(t *testing.T) {
	s := NewMemoryStore()
	if err := s.Create(Link{Code: "abc123", OriginalURL: "https://example.com"}); err != nil {
		t.Fatal(err)
	}

	clicks := []ClickEvent{
		{Timestamp: time.Date(2026, 10, 3, 12, 0, 0, 0, time.UTC), Referrer: "https://search.example"},
		{Timestamp: time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC), Referrer: ""},
		{Timestamp: time.Date(2026, 10, 3, 13, 0, 0, 0, time.UTC), Referrer: "https://search.example"},
		{Timestamp: time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC), Referrer: "https://other.example"},
	}
	for _, click := range clicks {
		if err := s.RecordClick("abc123", click); err != nil {
			t.Fatal(err)
		}
	}

	stats, err := s.GetStats("abc123")
	if err != nil {
		t.Fatalf("GetStats() returned an error: %v", err)
	}
	wantDays := []DayCount{
		{Date: "2026-10-01", Count: 1},
		{Date: "2026-10-02", Count: 1},
		{Date: "2026-10-03", Count: 2},
	}
	wantReferrers := []ReferrerCount{
		{Referrer: "https://search.example", Count: 2},
		{Referrer: "direct", Count: 1},
		{Referrer: "https://other.example", Count: 1},
	}
	if stats.Code != "abc123" || stats.TotalClicks != len(clicks) {
		t.Fatalf("stats summary = %#v, want code and %d clicks", stats, len(clicks))
	}
	if !reflect.DeepEqual(stats.ClicksByDay, wantDays) {
		t.Fatalf("ClicksByDay = %#v, want %#v", stats.ClicksByDay, wantDays)
	}
	if !reflect.DeepEqual(stats.TopReferrers, wantReferrers) {
		t.Fatalf("TopReferrers = %#v, want %#v", stats.TopReferrers, wantReferrers)
	}
}
