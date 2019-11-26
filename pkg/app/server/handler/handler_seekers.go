package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/auth/session"
	csrf "2019_2_IBAT/pkg/pkg/csrf"
	. "2019_2_IBAT/pkg/pkg/interfaces"
)

func (h *Handler) CreateSeeker(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uuid, err := h.UserService.CreateSeeker(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write(errJSON)
		return
	}

	sessInfo, err := h.AuthService.CreateSession(r.Context(), &session.Session{
		Id:    uuid.String(),
		Class: SeekerStr,
	})

	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "error while unmarshaling")
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write(errJSON)
		return
	}

	token, err := csrf.Tokens.Create(sessInfo.ID, sessInfo.Cookie, time.Now().Add(24*time.Hour).Unix())
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
	RoleJSON, _ := json.Marshal(Role{Role: sessInfo.Role})

	w.Write(RoleJSON)
}

//should test method
func (h *Handler) GetSeekerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	seekId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write(errJSON)
		return
	}

	seeker, err := h.UserService.GetSeeker(seekId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write(errJSON)
		return
	}

	seeker.Password = "" //danger
	seekerJSON, _ := json.Marshal(seeker)

	w.Write(seekerJSON)
}
