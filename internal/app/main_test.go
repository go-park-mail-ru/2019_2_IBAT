package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

var expectedJSON = `[{"id":1,"name":"Afanasiy"},{"id":2,"name":"Ka"}]`

func TestGetUsers(t *testing.T) {

	h := Handlers{
		users: []User{
			{
				ID:       1,
				Name:     "Afanasiy",
				Password: "1234",
			},
			{
				ID:       2,
				Name:     "Ka",
				Password: "dfdffgdfg",
			},
		},
		mu: &sync.Mutex{},
	}

	t.Parallel()
	r := httptest.NewRequest("GET", "/users/", nil)
	w := httptest.NewRecorder()
	h.HandleListUsers(w, r)

	// t.Log(w.Code)
	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}
	bytes, _ := ioutil.ReadAll(w.Body)
	// t.Log(w.Body)
	if string(bytes) != expectedJSON {
		t.Errorf("expected: %s, got %s", expectedJSON, string(bytes))
	}
}

func TestCreateUsers(t *testing.T) {
	h := Handlers{
		users: []User{},
		mu:    &sync.Mutex{},
	}

	t.Parallel()
	body := bytes.NewReader([]byte(`{"name": "Ignat", "passord": "12"}`))
	var expectedusers = []User{
		{
			ID:       1,
			Name:     "Afanasiy",
			Password: "1234",
		},
		{
			ID:       2,
			Name:     "Ka",
			Password: "dfdffgdfg",
		},
		{
			ID:       3,
			Name:     "Ignat",
			Password: "12",
		},
	}
	r := httptest.NewRequest("POST", "/users/", body)
	w := httptest.NewRecorder()

	h.HandleCreateUser(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	reflect.DeepEqual(h.users, expectedusers)

}
