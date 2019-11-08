package handler

import (
	"2019_2_IBAT/internal/pkg/auth"
	csrf "2019_2_IBAT/internal/pkg/csrf"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"
	"log"
	"time"

	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateEmployer(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	uuid, err := h.UserService.CreateEmployer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	authInfo, cookieValue, err := h.AuthService.CreateSession(uuid, EmployerStr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	token, err := csrf.Tokens.Create(authInfo.ID.String(), cookieValue, time.Now().Add(24*time.Hour).Unix())
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

	w.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
	w.Header().Set("X-Csrf-Token", token)
	http.SetCookie(w, &cookie)
	RoleJSON, _ := json.Marshal(Role{Role: authInfo.Role})

	w.Write([]byte(RoleJSON))
}

func (h *Handler) GetEmployerById(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	emplId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	employer, err := h.UserService.GetEmployer(emplId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	employer.Password = "" //danger
	employerJSON, _ := json.Marshal(employer)

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetEmployers(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := h.ParseEmplQuery(r.URL.Query())
	employers, err := h.UserService.GetEmployers(params)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	employerJSON, _ := json.Marshal(employers)

	w.Write([]byte(employerJSON))
}

func (h *Handler) ParseEmplQuery(query url.Values) map[string]interface{} {
	params := make(map[string]interface{})

	if query.Get("company_name") != "" {
		params["company_name"] = query.Get("company_name")
	} else {
		if query.Get("empl_num") != "" {
			params["empl_num"] = query.Get("empl_num")
		}
		if query.Get("region") != "" {
			params["region"] = query.Get("region")
		}
	}

	return params
}
