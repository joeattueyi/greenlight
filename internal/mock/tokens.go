package mock

import (
	"time"

	"greenlight.jattueyi.com/internal/data"
)

var mockToken = data.Token{
	Plaintext: "token",
	Hash:      []byte("token"),
	UserID:    1,
	Expiry:    time.Now().AddDate(1, 0, 0),
	Scope:     "authentication",
}

type TokenModel struct{}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	return &mockToken, nil
}

func (m TokenModel) Insert(token *data.Token) error {
	if token.Plaintext != mockToken.Plaintext {
		return data.ErrRecordNotFound
	}
	return nil
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	return nil
}
