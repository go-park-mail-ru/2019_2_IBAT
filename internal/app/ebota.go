package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type UserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

type Handlers struct {
	users []User
	mu    *sync.Mutex
}

func (h *Handlers) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		w.Write([]byte("{}"))
		return
	}

	newUserInput := new(UserInput)
	err = json.Unmarshal(bytes, newUserInput)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		w.Write([]byte("{}"))
		return
	}

	h.mu.Lock()
	var id uint64 = 0
	if len(h.users) != 0 {
		id = uint64(h.users[len(h.users)-1].ID + 1)
	}

	h.users = append(h.users, User{
		ID:       id,
		Name:     newUserInput.Name,
		Password: newUserInput.Password,
	})
	h.mu.Unlock()

}

func (h *Handlers) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	usersJSON, err := json.Marshal(h.users)
	h.mu.Unlock()
	if err != nil {
		log.Printf("Error while marshalong: %s", err)
		w.Write([]byte("{}"))
		return
	}

	w.Write(usersJSON)
}

func main() {
	handlers := Handlers{
		users: make([]User, 0),
		mu:    &sync.Mutex{},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "applcation/json")

		if r.Method == http.MethodPost {
			handlers.HandleCreateUser(w, r)
		}

		handlers.HandleListUsers(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
