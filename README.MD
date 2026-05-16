# URL Shortener

A lightweight URL shortener service written in Go. Submit a long URL, get a short one back. Visit the short link and get redirected — 302, no caching, click-tracking-ready.

## Requirements

- Go 1.21+

## Installation

```bash
git clone https://github.com/amerucandior/go-url-shortener.git
cd go-url-shortener
go run main.go
```

Server starts on `http://localhost:8080`.

## API

### Shorten a URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.google.com"}'
```

**Response**

```json
{"short_url":"http://localhost:8080/gVKMvG"}
```

### Visit a short link

```bash
curl -v http://localhost:8080/gVKMvG
```

Returns `302 Found` with a `Location` header pointing to the original URL.

### Error cases

| Scenario | Status |
|---|---|
| Wrong method on `/shorten` | `405 Method Not Allowed` |
| Code not found | `404 Not Found` |
| Malformed JSON body | `400 Bad Request` |

## How it works

1. `POST /shorten` accepts a JSON body, generates a random 6-character base62 code, stores the mapping in memory, and returns the short URL.
2. `GET /{code}` looks up the code and redirects to the original URL.

> **Note:** Storage is in-memory. All links are lost on restart. Persistence on the roadmap.

## License

MIT