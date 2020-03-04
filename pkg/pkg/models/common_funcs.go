package models

import (
	"net/http"

	"github.com/google/uuid"
)

func UuidsToStrings(ids []uuid.UUID) []string {
	var strIDs []string
	if ids == nil {
		return strIDs
	}

	for _, id := range ids {
		strIDs = append(strIDs, id.String())
	}
	return strIDs
}

func StringsToUuids(strIDs []string) []uuid.UUID {
	var ids []uuid.UUID
	if strIDs == nil {
		return ids
	}

	for _, id := range strIDs {
		ids = append(ids, uuid.MustParse(id))
	}
	return ids
}

func SetError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errJSON, _ := Error{Message: msg}.MarshalJSON()
	w.Write(errJSON)
	return
}
