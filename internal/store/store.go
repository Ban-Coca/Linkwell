package store

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("link not found")

var ErrAlreadyExists = errors.New("link already exists")

type Link struct {
	Code 	string
	OriginalURL string
	CreatedAt time.Time
}

type LinkStore interface {
	Create(link Link) error
	Get(code string) (Link, error)
	Exists(code string) (bool, error)
}