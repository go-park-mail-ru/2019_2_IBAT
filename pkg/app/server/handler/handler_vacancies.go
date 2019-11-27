package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	. "2019_2_IBAT/pkg/pkg/models"
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

	var vacancies VacancySlice
	vacancies, err := h.UserService.GetVacancies(authInfo, params, tags) //err handle

	if err != nil {
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
		return
	}

	vacanciesJSON, _ := vacancies.MarshalJSON()
	w.Write(vacanciesJSON)
}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
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

		SetError(w, code, err.Error())

		return
	}

	idJSON, _ := Id{Id: id.String()}.MarshalJSON()

	w.Write(idJSON)
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())

	if !ok {
		authInfo = AuthStorageValue{}
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	vacancy, err := h.UserService.GetVacancy(vacId, authInfo)

	if err != nil {
		// var code int
		// switch err.Error() {
		// case ForbiddenMsg:
		// 	code = http.StatusForbidden
		// case UnauthorizedMsg:
		// 	code = http.StatusUnauthorized
		// case InternalErrorMsg:
		// 	code = http.StatusInternalServerError
		// default:
		// 	code = http.StatusBadRequest
		// }
		SetError(w, http.StatusBadRequest, InvalidIdMsg)

		return
	}

	vacancyJSON, _ := vacancy.MarshalJSON()

	w.Write(vacancyJSON)
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
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

		SetError(w, code, err.Error())

		return
	}
}

func (h *Handler) PutVacancy(w http.ResponseWriter, r *http.Request) { //+
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	err = h.UserService.PutVacancy(vacId, r.Body, authInfo)

	if err != nil {
		SetError(w, http.StatusForbidden, ForbiddenMsg)
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
