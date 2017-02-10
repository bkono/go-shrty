package shrty

import "errors"

// ShortenedURLService for shortening and expanding urls and tokens
type ShortenedURLService interface {
	Shorten(url string) (*ShortenedURL, error)
	Expand(token string) (string, error)
}

// TokenService generates tokens and decodes them to ids
type TokenService interface {
	Encode(su *ShortenedURL) (string, error)
	Decode(token string) (int64, error)
}

// Client creates a connection to the services
type Client interface {
	ShortenedURLService() ShortenedURLService
	TokenService() TokenService
}

// Shrty Errors
var (
	ErrTokenNotFound = errors.New("Token not found")
)
