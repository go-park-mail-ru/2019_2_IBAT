package users

import (
	"encoding/json"

	"io"
	"io/ioutil"

	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *UserService) CreateEmployer(body io.ReadCloser) (uuid.UUID, error) { //should do this part by one r with if?
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		// err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, errors.New("Invalid body, transfer error")
	}

	var newEmployerReg EmployerReg
	err = json.Unmarshal(bytes, &newEmployerReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, errors.New("Invalid JSON")
	}

	id, ok := h.Storage.CreateEmployer(newEmployerReg)

	if !ok {
		// log.Printf("Error while creating employer: %s", err)
		return uuid.UUID{}, errors.New("Email already exists")
	}

	return id, nil
}

func (h *UserService) PutEmployer(body io.ReadCloser, id uuid.UUID) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		return errors.Wrap(err, "reading body error")
	}

	var newEmployerReg EmployerReg
	err = json.Unmarshal(bytes, &newEmployerReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		return errors.Wrap(err, "unmarshaling error")
	}

	ok := h.Storage.PutEmployer(newEmployerReg, id)
	if !ok {
		// log.Println("Here inside users")
		// log.Printf("Error while creating employer: %s", err)
		return errors.New("Error while changing employer")
	}

	return nil
}

func (h *UserService) GetEmployer(id uuid.UUID) (Employer, error) {
	return h.Storage.GetEmployer(id)
}

func (h *UserService) GetEmployers() ([]Employer, error) {
	return h.Storage.GetEmployers()
}
