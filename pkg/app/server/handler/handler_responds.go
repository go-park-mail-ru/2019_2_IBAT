package handler

import (
	"fmt"
	"net/http"

	. "2019_2_IBAT/pkg/pkg/models"
)

func (h *Handler) GetResponds(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	v := r.URL.Query()
	params := make(map[string]string)
	params["vacancy_id"] = v.Get("vacancy_id")
	params["resume_id"] = v.Get("resume_id")
	fmt.Printf("vacancyid = %s, resumeid = %s", params["vacancyid"], params["resumeid"])

	var responds RespondSlice
	responds, _ = h.UserService.GetResponds(authInfo, params) //error handling

	respondsJSON, _ := responds.MarshalJSON()

	w.Write(respondsJSON)

}

func (h *Handler) CreateRespond(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
		return
	}

	err := h.UserService.CreateRespond(r.Body, authInfo)
	if err != nil {
		SetError(w, http.StatusForbidden, ForbiddenMsg)
		return
	}
}
