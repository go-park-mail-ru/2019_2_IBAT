package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/auth/session"
	csrf "2019_2_IBAT/pkg/pkg/csrf"
	. "2019_2_IBAT/pkg/pkg/models"
)

func (h *Handler) CreateSeeker(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uuid, err := h.UserService.CreateSeeker(r.Body)
	if err != nil {
		SetError(w, http.StatusBadRequest, err.Error())
		return
	}

	sessInfo, err := h.AuthService.CreateSession(r.Context(), &session.Session{
		Id:    uuid.String(),
		Class: SeekerStr,
	})

	if err != nil {
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
		return
	}

	token, err := csrf.Tokens.Create(sessInfo.ID, sessInfo.Cookie, time.Now().Add(24*time.Hour).Unix())
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
	RoleJSON, _ := Role{Role: sessInfo.Role}.MarshalJSON()

	w.Write(RoleJSON)
}

func (h *Handler) GetSeekerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	seekId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	seeker, err := h.UserService.GetSeeker(seekId)

	if err != nil {
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	seeker.Password = "" //danger
	seekerJSON, _ := seeker.MarshalJSON()

	w.Write(seekerJSON)
}
