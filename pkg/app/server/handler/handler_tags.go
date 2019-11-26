package handler

import (
	"encoding/json"
	"net/http"

	. "2019_2_IBAT/pkg/pkg/interfaces"
)

func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tags, err := h.UserService.GetTags() //err handle

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write(errJSON)
		return
	}

	tagsJSON, _ := json.Marshal(tags)

	w.Write(tagsJSON)

}
