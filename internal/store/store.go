package store

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("link not found")

var ErrAlreadyExists = errors.New("link already exists")

type ClickEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"userAgent"`
}

type Stats struct {
	Code          string `json:"code"`
	TotalClicks   int    `json:"totalClicks"`
	ClicksByDay   []DayCount `json:"clicksByDay"`
	TopReferrers  []ReferrerCount `json:"topReferrers"`
}

type DayCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ReferrerCount struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

type Link struct {
	Code        string       `json:"code"`
	OriginalURL string       `json:"originalUrl"`
	CreatedAt   time.Time    `json:"createdAt"`
	Clicks      []ClickEvent `json:"-"`
}

type LinkStore interface {
	Create(link Link) error
	Get(code string) (Link, error)
	Exists(code string) (bool, error)
	RecordClick(code string, event ClickEvent) error
	GetStats(code string) (Stats, error)
}