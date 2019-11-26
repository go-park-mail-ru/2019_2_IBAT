package users

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	. "2019_2_IBAT/pkg/pkg/interfaces"
)

func (h *UserService) CreateEmployer(body io.ReadCloser) (uuid.UUID, error) {
	bytes, err := ioutil.ReadAll(body)
	defer body.Close()

	if err != nil {
		log.Printf("error while reading body: %s", err)
		return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	var newEmployerReg Employer
	// id := uuid.New()
	// newEmployerReg.ID = id
	err = json.Unmarshal(bytes, &newEmployerReg)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	id := uuid.New()
	newEmployerReg.ID = id
	newEmployerReg.PathToImg = DefaultImg
	ok := h.Storage.CreateEmployer(newEmployerReg)

	if !ok {
		log.Printf("Error while creating employer: %s", err)
		return uuid.UUID{}, errors.New(EmailExistsMsg)
	}

	return id, nil
}

func (h *UserService) PutEmployer(body io.ReadCloser, id uuid.UUID) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		return errors.Wrap(err, BadRequestMsg)
	}

	var newEmployerReg EmployerReg
	err = json.Unmarshal(bytes, &newEmployerReg)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		return errors.New(InvalidJSONMsg)
	}

	ok := h.Storage.PutEmployer(newEmployerReg, id)
	if !ok {
		log.Printf("Error while creating employer")
		return errors.New(BadRequestMsg)
	}

	return nil
}

func (h *UserService) GetEmployer(id uuid.UUID) (Employer, error) {
	log.Println("GetEmployer Service Start")
	return h.Storage.GetEmployer(id)
}

func (h *UserService) GetEmployers(params map[string]interface{}) ([]Employer, error) {
	return h.Storage.GetEmployers(params)
}
