package shrty

import "github.com/gogo/protobuf/proto"

var _ ShortenedURLService = &shortenedURLService{}

type shortenedURLService struct {
	BaseURL string
	c       *DBClient
	ts      TokenService
}

// NewShortenedURLService returns a new instance of the ShortenedURLService
func NewShortenedURLService(baseURL string, c *DBClient, ts TokenService) ShortenedURLService {
	s := &shortenedURLService{baseURL, c, ts}
	return s
}

func (s *shortenedURLService) Shorten(url string) (*ShortenedURL, error) {
	tx, err := s.c.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(URLsBucket))
	id, err := b.NextSequence()
	if err != nil {
		return nil, err
	}

	u := &ShortenedURL{
		ID:          int64(id),
		OrigURL:     url,
		CreatedTime: s.c.Now().Unix(),
		Views:       0,
	}

	tk, err := s.ts.Encode(u)
	if err != nil {
		return nil, err
	}

	u.Token = tk
	u.ShrtURL = s.BaseURL + tk

	if v, err := proto.Marshal(u); err != nil {
		return nil, err
	} else if err := b.Put([]byte(u.Token), v); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *shortenedURLService) Expand(token string) (string, error) {
	tx, err := s.c.db.Begin(true)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Find it
	var u ShortenedURL
	b := tx.Bucket([]byte(URLsBucket))
	if v := b.Get([]byte(token)); v == nil {
		return "", ErrTokenNotFound
	} else if err := proto.Unmarshal(v, &u); err != nil {
		return "", err
	}

	u.Views++
	var e error
	if v, err := proto.Marshal(&u); err != nil {
	} else if err := b.Put([]byte(u.Token), v); err != nil {
	} else if err := tx.Commit(); err != nil {
	}

	return u.OrigURL, e
}
