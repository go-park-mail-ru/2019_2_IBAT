package handler

import (
	"net/http"

	. "2019_2_IBAT/pkg/pkg/models"
)

func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tags TagMap
	tags, err := h.UserService.GetTags() //err handle

	if err != nil {
		SetError(w, http.StatusInternalServerError, InternalErrorMsg)
		return
	}

	tagsJSON, _ := tags.MarshalJSON()

	w.Write(tagsJSON)

}
