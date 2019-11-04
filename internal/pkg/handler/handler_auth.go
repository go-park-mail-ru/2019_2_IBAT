package handler

import (
	"2019_2_IBAT/internal/pkg/auth"
	csrf "2019_2_IBAT/internal/pkg/csrf"
	. "2019_2_IBAT/internal/pkg/interfaces"

	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"net/http"

	"github.com/gorilla/context"
)

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("Handle CreateSession: start")
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Handle CreateSession: error while reading body")
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	userAuthInput := new(UserAuthInput)
	err = json.Unmarshal(bytes, userAuthInput)
	if err != nil {
		log.Println("Handle CreateSession: error while unmarshaling")
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	id, role, ok := h.UserService.CheckUser(userAuthInput.Email, userAuthInput.Password)
	if !ok {
		log.Println("Handle CreateSession: Check user failed")
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: BadRequest}) //
		w.Write([]byte(errJSON))
		return
	}

	// cookie, role, err := h.AuthService.CreateSession(id, role)
	authInfo, cookieValue, err := h.AuthService.CreateSession(id, role)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Create session failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	token, err := csrf.Tokens.Create(id.String(), cookieValue, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Create token failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	expiresAt, err := time.Parse(TimeFormat, authInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Time parsing failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	cookie := http.Cookie{
		Name:    auth.CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}

	w.Header().Set("X-CSRF-Token", token)
	http.SetCookie(w, &cookie)
	RoleJSON, _ := json.Marshal(Role{role})

	w.Write([]byte(RoleJSON))
	log.Println("Handle CreateSession: end")
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := context.Get(r, AuthRec).(AuthStorageValue)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	RoleJSON, _ := json.Marshal(Role{authInfo.Role})

	w.Write([]byte(RoleJSON))
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	ok := h.AuthService.DeleteSession(cookie)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	http.SetCookie(w, cookie)
}
