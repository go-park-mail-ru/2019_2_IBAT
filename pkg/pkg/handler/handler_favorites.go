package handler

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"
	"encoding/json"
	"log"

	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateFavorite(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	// vacId_string := r.URL.Query().Get("vacancy_id")
	// vacId, err := uuid.Parse(vacId_string)

	vacId, err := uuid.Parse(mux.Vars(r)["id"])

	// vacId, err := uuid.Parse(vacId_string)
	if err != nil {
		log.Printf("Handle CreateFavorite: invalid id - %s", err)
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	err = h.UserService.CreateFavorite(vacId, authInfo)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}

}

func (h *Handler) GetFavoriteVacancies(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	vacancies, _ := h.UserService.GetFavoriteVacancies(authInfo) //error handling

	respondsJSON, _ := json.Marshal(vacancies)

	w.Write([]byte(respondsJSON))

}

func (h *Handler) DeleteFavoriteVacancy(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(code)

		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}
}
