package handler

import (
	"2019_2_IBAT/pkg/pkg/auth"
	"2019_2_IBAT/pkg/pkg/auth/session"
	csrf "2019_2_IBAT/pkg/pkg/csrf"
	. "2019_2_IBAT/pkg/pkg/interfaces"
	"encoding/json"
	"log"
	"time"

	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateEmployer(w http.ResponseWriter, r *http.Request) { //+
	// defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("CreateEmployer Start")

	uuid, err := h.UserService.CreateEmployer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	log.Println("CreateEmployer Employer was created")

	sessInfo, err := h.AuthService.CreateSession(r.Context(), &session.Session{
		Id:    uuid.String(),
		Class: EmployerStr,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	token, err := csrf.Tokens.Create(sessInfo.ID, sessInfo.Cookie,
		time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Create token failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	expiresAt, err := time.Parse(TimeFormat, sessInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Time parsing failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
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

	w.Write([]byte(RoleJSON))
}

func (h *Handler) GetEmployerById(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id_string := mux.Vars(r)["id"]
	log.Printf("GetEmployerById id_string: %s\n", id_string)
	emplId, err := uuid.Parse(id_string)
	log.Println("GetEmployerById Handler Start")

	if err != nil {
		log.Printf("GetEmployerById Parse id error: %s\n", err)
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
