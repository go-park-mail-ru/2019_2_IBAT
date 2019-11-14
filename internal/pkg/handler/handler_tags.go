package handler

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"
	"io/ioutil"
	"log"

	"net/http"

)

func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tags, err := h.UserService.GetTags() //err handle

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	tagsJSON, _ := json.Marshal(tags)

	w.Write([]byte(tagsJSON))

}

func (h *Handler) TestUnmar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")


	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		// err = errors.Wrap(err, "reading body error")
		// return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	var resumeReg Resume
	// id := uuid.New()
	// resumeReg.ID = id
	// resumeReg.OwnerID = authInfo.ID
	err = json.Unmarshal(bytes, &resumeReg)

	log.Println(resumeReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "unmarshaling error")
		// return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	// if err != nil {
	// 	// w.WriteHeader(StatusBadRequest)

	// 	errJSON, _ := json.Marshal(Error{Message: err.Error()})
	// 	w.Write([]byte(errJSON))
	// 	return
	// }

	// idJSON, err := json.Marshal(Id{Id: id.String()})

	// if err != nil {
	// 	errJSON, _ := json.Marshal(Error{Message: err.Error()})
	// 	w.Write([]byte(errJSON))
	// 	return
	// }
	// log.Printf("Returning id: %s", id.String())

	// w.Write([]byte(idJSON))
}

// "spheres": [{"first": "string1", "second": "string2"},  {"first": "string2", "second": "string3"}]
