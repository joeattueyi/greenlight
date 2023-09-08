package mock

import (
	"time"

	"greenlight.jattueyi.com/internal/data"
)

var passString = "pa55word"
var mockPassword = data.Password{
	Plaintext: &passString,
	Hash:      []byte("ABC"),
}
var mockUser = data.User{
	ID:        1,
	CreatedAt: time.Now(),
	Name:      "alice freeman",
	Email:     "alice@foo.com",
	Activated: true,
	Version:   0,
	Password:  mockPassword,
}

type UserModel struct{}

func (m UserModel) Insert(user *data.User) error {
	switch user.Email {
	case "dupe@foo.com":
		return data.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m UserModel) GetByEmail(email string) (*data.User, error) {

	switch email {
	case "alice@gfoo.com":
		return &mockUser, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (m UserModel) Update(user *data.User) error {
	switch user.Email {
	case "dupe@foo.com":
		return data.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*data.User, error) {

	return &mockUser, nil
}
