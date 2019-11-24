package handler

import (
	"2019_2_IBAT/pkg/pkg/auth"
	"2019_2_IBAT/pkg/pkg/auth/session"
	csrf "2019_2_IBAT/pkg/pkg/csrf"
	. "2019_2_IBAT/pkg/pkg/interfaces"

	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"net/http"
)

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("Handle CreateSession: start")
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Handle CreateSession: error while reading body")
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: BadRequestMsg})
		w.Write(errJSON)
		return
	}

	userAuthInput := new(UserAuthInput)
	err = json.Unmarshal(bytes, userAuthInput)
	if err != nil {
		log.Println("Handle CreateSession: error while unmarshaling")
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidJSONMsg})
		w.Write(errJSON)
		return
	}

	id, role, ok := h.UserService.CheckUser(userAuthInput.Email, userAuthInput.Password)
	if !ok {
		log.Println("Handle CreateSession: Check user failed")
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvPassOrEmailMsg}) //
		w.Write(errJSON)
		return
	}

	sessInfo, err := h.AuthService.CreateSession(r.Context(), &session.Session{
		Id:    id.String(),
		Class: role,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Create session failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write(errJSON)
		return
	}

	token, err := csrf.Tokens.Create(id.String(), sessInfo.Cookie, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Create token failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write(errJSON)
		return
	}

	expiresAt, err := time.Parse(TimeFormat, sessInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Time parsing failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write(errJSON)
		return
	}

	cookie := http.Cookie{
		Name:    auth.CookieName,
		Value:   sessInfo.Cookie,
		Expires: expiresAt,
	}

	w.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
	w.Header().Set("X-Csrf-Token", token)
	http.SetCookie(w, &cookie)
	RoleJSON, _ := json.Marshal(Role{Role: role})

	w.Write(RoleJSON)
	log.Println("Handle CreateSession: end")
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write(errJSON)
		return
	}

	RoleJSON, _ := json.Marshal(Role{Role: authInfo.Role})

	w.Write(RoleJSON)
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write(errJSON)
		return
	}

	sessionBool, err := h.AuthService.DeleteSession(r.Context(), &session.Cookie{
		Cookie: cookie.Value,
	})
	if !sessionBool.Ok {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: BadRequestMsg})
		w.Write(errJSON)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
}
