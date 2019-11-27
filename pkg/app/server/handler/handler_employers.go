package handler

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/auth/session"
	csrf "2019_2_IBAT/pkg/pkg/csrf"
	. "2019_2_IBAT/pkg/pkg/models"
)

func (h *Handler) CreateEmployer(w http.ResponseWriter, r *http.Request) { //+
	// defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("CreateEmployer Start")

	uuid, err := h.UserService.CreateEmployer(r.Body)
	if err != nil {
		SetError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("CreateEmployer Employer was created")

	sessInfo, err := h.AuthService.CreateSession(context.Background(), &session.Session{
		Id:    uuid.String(),
		Class: EmployerStr,
	})

	if err != nil {
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
		return
	}

	token, err := csrf.Tokens.Create(sessInfo.ID, sessInfo.Cookie,
		time.Now().Add(24*time.Hour).Unix())
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

func (h *Handler) GetEmployerById(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id_string := mux.Vars(r)["id"]
	log.Printf("GetEmployerById id_string: %s\n", id_string)
	emplId, err := uuid.Parse(id_string)
	log.Println("GetEmployerById Handler Start")

	if err != nil {
		log.Printf("GetEmployerById Parse id error: %s\n", err)
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	employer, err := h.UserService.GetEmployer(emplId)

	if err != nil {
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	employer.Password = "" //danger
	employerJSON, _ := employer.MarshalJSON()

	w.Write(employerJSON)
}

func (h *Handler) GetEmployers(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := h.ParseEmplQuery(r.URL.Query())

	var employers EmployerSlice
	employers, err := h.UserService.GetEmployers(params)

	if err != nil {
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
		return
	}

	employerJSON, _ := employers.MarshalJSON()

	w.Write(employerJSON)
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
