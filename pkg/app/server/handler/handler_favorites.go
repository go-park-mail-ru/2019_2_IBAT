package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	. "2019_2_IBAT/pkg/pkg/models"
)

func (h *Handler) CreateFavorite(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	vacId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("Handle CreateFavorite: invalid id - %s", err)
		SetError(w, http.StatusBadRequest, InvalidIdMsg)
		return
	}

	err = h.UserService.CreateFavorite(vacId, authInfo)

	if err != nil {
		SetError(w, http.StatusForbidden, ForbiddenMsg)
		return
	}

}

func (h *Handler) GetFavoriteVacancies(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	var vacancies VacancySlice
	vacancies, _ = h.UserService.GetFavoriteVacancies(authInfo) //error handling

	respondsJSON, _ := vacancies.MarshalJSON()

	w.Write(respondsJSON)

}

func (h *Handler) DeleteFavoriteVacancy(w http.ResponseWriter, r *http.Request) {
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

	err = h.UserService.DeleteFavoriteVacancy(vacId, authInfo)

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
