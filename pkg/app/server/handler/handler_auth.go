package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/auth/session"
	csrf "2019_2_IBAT/pkg/pkg/csrf"
	. "2019_2_IBAT/pkg/pkg/models"
)

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("Handle CreateSession: start")
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Handle CreateSession: error while reading body")
		SetError(w, http.StatusBadRequest, BadRequestMsg)
		return
	}

	userAuthInput := new(UserAuthInput)
	err = json.Unmarshal(bytes, userAuthInput)
	if err != nil {
		log.Println("Handle CreateSession: error while unmarshaling")
		SetError(w, http.StatusBadRequest, InvalidJSONMsg)
		return
	}

	id, role, ok := h.UserService.CheckUser(userAuthInput.Email, userAuthInput.Password)
	if !ok {
		log.Println("Handle CreateSession: Check user failed")
		SetError(w, http.StatusBadRequest, InvPassOrEmailMsg)
		return
	}

	sessInfo, err := h.AuthService.CreateSession(context.Background(), &session.Session{
		Id:    id.String(),
		Class: role,
	})

	if err != nil {
		log.Println("Handle CreateSession:  Create session failed")
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
		return
	}

	token, err := csrf.Tokens.Create(id.String(), sessInfo.Cookie, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		log.Println("Handle CreateSession:  Create token failed")
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
		return
	}

	expiresAt, err := time.Parse(TimeFormat, sessInfo.Expires)
	if err != nil {
		log.Println("Handle CreateSession:  Time parsing failed")
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
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
	RoleJSON, _ := Role{Role: role}.MarshalJSON()

	w.Write(RoleJSON)
	log.Println("Handle CreateSession: end")
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	RoleJSON, _ := Role{Role: authInfo.Role}.MarshalJSON()

	w.Write(RoleJSON)
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cookie, err := r.Cookie(auth.CookieName)

	if err != nil {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	sessionBool, err := h.AuthService.DeleteSession(context.Background(), &session.Cookie{
		Cookie: cookie.Value,
	})
	if !sessionBool.Ok {
		SetError(w, http.StatusBadRequest, BadRequestMsg)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
}
