package users

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"

	. "2019_2_IBAT/pkg/pkg/models"
)

func (h *UserService) CreateRespond(body io.ReadCloser, record AuthStorageValue) error { //should do this part by one r with if?
	if record.Role != SeekerStr {
		// log.Printf("Invalid action: %s", err)
		return errors.New("Invalid action")
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return err
	}

	var respond Respond
	err = respond.UnmarshalJSON(bytes)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return err
	}
	respond.Status = AwaitSt

	ok := h.Storage.CreateRespond(respond, record.ID)

	if !ok {
		log.Printf("Error while creating respond: %s", err)
		return errors.New("Error while creating respond")
	}

	return nil
}

func (h *UserService) GetResponds(authInfo AuthStorageValue, params map[string]string) ([]Respond, error) {
	responds := []Respond{}

	if params["resume_id"] != "" && params["vacancy_id"] != "" {
		return responds, errors.New(BadRequestMsg)
	}

	responds, err := h.Storage.GetResponds(authInfo, params)
	if err != nil {
		return responds, errors.New("Invalid action") ///
	}

	return responds, nil
}
