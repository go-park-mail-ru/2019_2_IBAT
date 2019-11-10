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

// write tcp 127.0.0.1:37786->127.0.0.1:6379: use of closed network connection
// Can not get auth info: redigo: unexpected response line
// (possible server error or unsupported concurrent read by application)

// fail
// &{GET /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f HTTP/2.0 2 0
// 	map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Accept-Language:
// 	[ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7] Cookie:[session-id=uoa7If53qWDgQdFR3DrGSeZXL9nLnBZT]
// 	Origin:[http://localhost:8080] Referer:[http://localhost:8080/]
// 	User-Agent:[Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36
// 	(KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36]]
// 	0xc0000751a0 <nil> 0 [] false 82.146.43.113:8080 map[]
// 	map[] <nil> map[] 93.171.198.4:37498 /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f
// 	0xc000268c60 <nil> <nil> 0xc0000260d0}
