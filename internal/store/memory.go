package store

import (
	"sort"
	"strings"
	"sync"
)

type MemoryStore struct {
	mu    sync.RWMutex
	links map[string]Link
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		links: make(map[string]Link),
	}
}

func (s *MemoryStore) Create(link Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.links[link.Code]; ok {
		return ErrAlreadyExists
	}

	if link.Clicks == nil {
		link.Clicks = []ClickEvent{}
	}

	s.links[link.Code] = link
	return nil
}

func (s *MemoryStore) Get(code string) (Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, ok := s.links[code]
	if !ok {
		return Link{}, ErrNotFound
	}
	return link, nil
}

func (s *MemoryStore) Exists(code string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.links[code]
	return ok, nil
}

func (s *MemoryStore) RecordClick(code string, event ClickEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	link, ok := s.links[code]
	if !ok {
		return ErrNotFound
	}

	link.Clicks = append(link.Clicks, event)
	s.links[code] = link
	return nil
}

func (s *MemoryStore) GetStats(code string) (Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, ok := s.links[code]
	if !ok {
		return Stats{}, ErrNotFound
	}

	stats := Stats{
		Code: code,
	}
	countsByDay := map[string]int{}
	referrerCounts := map[string]int{}

	for _, click := range link.Clicks {
		stats.TotalClicks++
		day := click.Timestamp.Format("2006-01-02")
		countsByDay[day]++
		ref := click.Referrer
		if ref == "" {
			ref = "direct"
		}
		referrerCounts[ref]++
	}

	for day, count := range countsByDay {
		stats.ClicksByDay = append(stats.ClicksByDay, DayCount{Date: day, Count: count})
	}
	sort.Slice(stats.ClicksByDay, func(i, j int) bool {
		return stats.ClicksByDay[i].Date < stats.ClicksByDay[j].Date
	})

	for referrer, count := range referrerCounts {
		stats.TopReferrers = append(stats.TopReferrers, ReferrerCount{Referrer: referrer, Count: count})
	}
	sort.Slice(stats.TopReferrers, func(i, j int) bool {
		if stats.TopReferrers[i].Count == stats.TopReferrers[j].Count {
			return strings.ToLower(stats.TopReferrers[i].Referrer) < strings.ToLower(stats.TopReferrers[j].Referrer)
		}
		return stats.TopReferrers[i].Count > stats.TopReferrers[j].Count
	})

	return stats, nil
}

var _ LinkStore = (*MemoryStore)(nil)