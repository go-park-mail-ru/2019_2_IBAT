package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "2019_2_IBAT/pkg/pkg/interfaces"
)

func (h *Handler) GetResponds(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write(errJSON)
	}

	v := r.URL.Query()
	params := make(map[string]string)
	params["vacancy_id"] = v.Get("vacancy_id")
	params["resume_id"] = v.Get("resume_id")
	fmt.Printf("vacancyid = %s, resumeid = %s", params["vacancyid"], params["resumeid"])

	responds, _ := h.UserService.GetResponds(authInfo, params) //error handling

	respondsJSON, _ := json.Marshal(responds)

	w.Write(respondsJSON)

}

func (h *Handler) CreateRespond(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write(errJSON)
		return
	}

	err := h.UserService.CreateRespond(r.Body, authInfo)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write(errJSON)
		return
	}
}
