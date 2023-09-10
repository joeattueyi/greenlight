package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {

	rr := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := &Application{}
	app.healthcheckHandler(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()

	type SystemInfo struct {
		environment string
		status      string
		version     string
	}

	var got struct {
		available  string
		system_inf SystemInfo
	}

	err = json.NewDecoder(rs.Body).Decode(&got)
	if err != nil {
		t.Fatal("unable to parse response")
	}
}
