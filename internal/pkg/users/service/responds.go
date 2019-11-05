package users

import (
	"encoding/json"

	"io"
	"io/ioutil"
	"log"

	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *UserService) CreateRespond(body io.ReadCloser, record AuthStorageValue) (uuid.UUID, error) { //should do this part by one r with if?
	if record.Role != SeekerStr {
		// log.Printf("Invalid action: %s", err)
		return uuid.UUID{}, errors.New("Invalid action")
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, err
	}

	var respond Respond
	err = json.Unmarshal(bytes, &respond)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, err
	}
	respond.Status = AwaitSt

	id, ok := h.Storage.CreateRespond(respond, record.ID)

	if !ok {
		log.Printf("Error while creating respond: %s", err)
		return uuid.UUID{}, errors.New("Error while creating respond")
	}

	return id, nil
}

func (h *UserService) GetResponds(authInfo AuthStorageValue, params map[string]string) ([]Respond, error) {
	responds := []Respond{}

	if params["resumeid"] != "" && params["vacancyid"] != "" {
		return responds, errors.New("Invalid message")
	}

	responds, err := h.Storage.GetResponds(authInfo, params)
	if err != nil {
		return responds, errors.New("Invalid action") ///
	}

	return responds, nil
}
