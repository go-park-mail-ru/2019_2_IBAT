package handler

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"

	"net/http"
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

	err := h.UserService.CreateFavorite(r.Body, authInfo)
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
