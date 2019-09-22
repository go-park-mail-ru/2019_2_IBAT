package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AuthStorage interface {
	// Check(cookie string)
	Get(cookie string) (StorageValue, bool)
	Set(id uint64) string
	Delete(cookie string) string
}

type StorageValue struct {
	ID      uuid.UUID
	Expires string
}

type UserAuthInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Handlers struct {
	Storage AuthStorage
	Mu      *sync.Mutex
}

const CookieName = "session-id"

func (h *Handlers) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return
	}

	userAuthInput := new(UserAuthInput)
	err = json.Unmarshal(bytes, userAuthInput)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return
	}

	h.Mu.Lock()
	var id uuid.UUID = 1             //should get id from users struct
	cookieValue := h.Storage.Set(id) //possible return authInfo
	h.Mu.Unlock()

	authInfo, ok := h.Storage.Get(cookieValue) //impossible error, should use only Set method
	if !ok {
		log.Printf("Error: %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return
	}

	expiresAt, err := time.Parse(TimeFormat, authInfo.Expires)
	if err != nil {
		log.Printf("Error while time conversing: %s", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return
	} //impossible error

	cookie := http.Cookie{
		Name:    CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("{}"))
}

func (h *Handlers) DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No correct session cookie detected"))
		return
	}

	fmt.Printf("cookie value: %s\n", cookie.Value)
	_, ok := h.Storage.Get(cookie.Value) //is there need in this
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No session detected"))
		return
	}

	h.Storage.Delete(cookie.Value)
	cookie.Expires = time.Now().AddDate(0, 0, -1)

	http.SetCookie(w, cookie)
	w.Write([]byte("Cookie deleted"))
}
