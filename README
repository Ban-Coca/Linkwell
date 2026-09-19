# Linkwell
 
A link shortener with real click analytics — built in Go to learn Go's concurrency model (goroutines, channels) rather than just port a Spring Boot app.
 
## What it does
 
- Shortens long URLs into short codes (base62, 6–7 chars)
- Redirects short links instantly, without waiting on analytics writes
- Tracks clicks (timestamp, referrer, user-agent) in the background via a buffered channel + worker goroutine
- Serves click stats: totals, clicks per day, top referrers
## Why
 
Redirect handling and click tracking are decoupled on purpose — the redirect response returns immediately while a separate goroutine drains a channel and persists the click event. That's the core learning point of the project: reaching for a channel instead of a thread pool.
 
## Stack
 
- Go, standard library-first (`net/http`)
- SQLite for persistence (in-memory map used in early development)
- No external framework
## Getting started
 
```bash
go mod init github.com/yourname/linkwell
go run .
```
 
Server starts on `:8080`.
 
## API
 
| Method | Route                       | Description                        |
|--------|------------------------------|-------------------------------------|
| POST   | `/api/links`                 | Create a short link from a long URL |
| GET    | `/:code`                     | Redirect to the original URL        |
| GET    | `/api/links/:code/stats`     | Get click analytics for a link      |
 
### Create a link
 
```bash
curl -X POST localhost:8080/api/links \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/some/long/path"}'
```
 
### Follow it
 
```bash
curl -L localhost:8080/abc123
```
 
### Check stats
 
```bash
curl localhost:8080/api/links/abc123/stats
```
 
## Project structure
 
```
main.go                  # wiring: store, handlers, server startup
internal/
  store/                 # LinkStore interface + in-memory / SQLite implementations
  shortener/              # base62 code generation
  clicks/                # ClickEvent + buffered channel + worker
  ratelimit/              # hand-rolled token-bucket limiter
  handlers/               # HTTP handlers
migrations/               # SQL schema
```
 
## Status
 
Work in progress — built phase by phase:
 
- [ ] Phase 1 — basic shortening + redirect
- [ ] Phase 2 — concurrent click tracking
- [ ] Phase 3 — analytics API
- [ ] Phase 4 — rate limiting