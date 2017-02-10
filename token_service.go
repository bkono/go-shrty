package shrty

import hashids "github.com/speps/go-hashids"

type tokenService struct {
	salt string
	h    *hashids.HashID
}

// NewTokenService generates a new TokenService using the provided salt
func NewTokenService(salt string) TokenService {
	ts := &tokenService{salt: salt}
	hd := hashids.NewData()
	hd.Salt = salt
	h := hashids.NewWithData(hd)
	ts.h = h
	return ts
}

func (ts *tokenService) Encode(su *ShortenedURL) (string, error) {
	e, err := ts.h.EncodeInt64([]int64{su.ID})
	if err != nil {
		return "", err
	}

	return e, nil
}

func (ts *tokenService) Decode(token string) (int64, error) {
	d, err := ts.h.DecodeInt64WithError(token)
	if err != nil {
		return 0, err
	}

	return d[0], nil
}
