package handler

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	vacancyParams = [...]string{
		"position",
		"region",
		"wage",
		"experience",
		"type_of_employment",
		"work_schedule",
	}
)

func (h *Handler) GetVacancies(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

	if !ok {
		authInfo = AuthStorageValue{} //check if nil possible
	}

	params := h.ParseVacanciesQuery(r.URL.Query())
	tags := h.ParseTags(r.URL.Query())

	log.Printf("Params map length: %d\n", len(params))

	vacancies, err := h.UserService.GetVacancies(authInfo, params, tags) //err handle

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacanciesJSON, _ := json.Marshal(vacancies)

	w.Write([]byte(vacanciesJSON))

}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserService.CreateVacancy(r.Body, authInfo)
	if err != nil {
		var code int
		switch err.Error() {
		case ForbiddenMsg:
			code = http.StatusForbidden
		case UnauthorizedMsg:
			code = http.StatusUnauthorized
		case InternalErrorMsg:
			code = http.StatusInternalServerError
		default:
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)

		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Id{Id: id.String()})

	w.Write([]byte(idJSON))
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

	if !ok {
		authInfo = AuthStorageValue{} //check if nil possible
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		var code int
		switch err.Error() {
		case ForbiddenMsg:
			code = http.StatusForbidden
		case UnauthorizedMsg:
			code = http.StatusUnauthorized
		case InternalErrorMsg:
			code = http.StatusInternalServerError
		default:
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)

		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	vacancy, err := h.UserService.GetVacancy(vacId, authInfo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacancyJSON, _ := json.Marshal(vacancy)

	w.Write([]byte(vacancyJSON))
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.DeleteVacancy(vacId, authInfo)

	if err != nil {
		var code int
		switch err.Error() {
		case ForbiddenMsg:
			code = http.StatusForbidden
		case UnauthorizedMsg:
			code = http.StatusUnauthorized
		case InternalErrorMsg:
			code = http.StatusInternalServerError
		default:
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)

		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) PutVacancy(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.PutVacancy(vacId, r.Body, authInfo)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}
}

func (h *Handler) ParseVacanciesQuery(query url.Values) map[string]interface{} {
	params := make(map[string]interface{})

	if query.Get("recommended") != "" {
		params["recommended"] = query.Get("recommended")
		return params //no sense to continue
	}

	if query.Get("position") != "" {
		params["position"] = query.Get("position")
	}

	if query.Get("region") != "" {
		params["region"] = query.Get("region")
	}
	if query.Get("wage") != "" {
		params["wage_from"] = query.Get("wage")
	}
	if query.Get("experience") != "" {
		params["experience"] = query.Get("experience")
	}
	if query.Get("type_of_employment") != "" {
		params["type_of_employment"] = query.Get("type_of_employment")
	}
	if query.Get("work_schedule") != "" {
		params["work_schedule"] = query.Get("work_schedule")
	}

	return params
}

func (h *Handler) ParseTags(query url.Values) map[string]interface{} {
	params := make(map[string]interface{})

	for i, item := range query {
		fmt.Printf("%s  %s\n", i, item)
		flag := false

		for _, param := range vacancyParams {
			if i == param {
				flag = true
				break
			}
		}

		if !flag {
			params[i] = item
			fmt.Printf("%s  %s\n", i, params[i])
		}

	}

	log.Println("Tag array")
	log.Println(params)
	return params
}
