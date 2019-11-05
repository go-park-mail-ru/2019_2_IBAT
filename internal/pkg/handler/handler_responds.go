package handler

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"
	"fmt"

	"net/http"
)

func (h *Handler) GetResponds(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
	}

	v := r.URL.Query()
	params := make(map[string]string)
	params["vacancyid"] = v.Get("vacancyid")
	params["resumeid"] = v.Get("resumeid")
	fmt.Printf("vacancyid = %s, resumeid = %s", params["vacancyid"], params["resumeid"])

	responds, _ := h.UserService.GetResponds(authInfo, params) //error handling

	respondsJSON, _ := json.Marshal(responds)

	w.Write([]byte(respondsJSON))

}

func (h *Handler) CreateRespond(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authInfo, ok := FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		errJSON, _ := json.Marshal(Error{Message: UnauthorizedMsg})
		w.Write([]byte(errJSON))
		return
	}

	id, err := h.UserService.CreateRespond(r.Body, authInfo)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errJSON, _ := json.Marshal(Error{Message: ForbiddenMsg})
		w.Write([]byte(errJSON))
		return
	}

	idJSON, _ := json.Marshal(Id{Id: id.String()})

	w.Write([]byte(idJSON))
}
