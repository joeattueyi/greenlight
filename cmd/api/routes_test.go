package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"greenlight.jattueyi.com/internal/data"
	"greenlight.jattueyi.com/internal/mock"
)

var app = &Application{
	models: data.Models{
		Users:       mock.UserModel{},
		Movies:      mock.MovieModel{},
		Permissions: mock.PermissionModel{},
		Tokens:      mock.TokenModel{},
	},
	mailer: mock.MockMailer{},
}

func TestRegisterUserHandler(t *testing.T) {

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	rs, err := ts.Client().Post(ts.URL+"/v1/users", "application/json", strings.NewReader(`{
		"name": "alice freeman",
		"email": "alice@foo.com",
		"password": "pa55word"
	}`))
	if err != nil {
		t.Fatal(err)
	}

	if rs.StatusCode != http.StatusAccepted {
		t.Errorf("got %d", rs.StatusCode)
	}
	defer rs.Body.Close()

	type AuthToken struct {
		token  string
		expiry string
	}
	var out struct {
		authentication_token AuthToken
	}

	err = json.NewDecoder(rs.Body).Decode(&out)
	if err != nil {
		t.Fatal("unable to parse result")
	}

}
