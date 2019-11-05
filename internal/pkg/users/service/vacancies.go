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

func (h *UserService) CreateVacancy(body io.ReadCloser, authInfo AuthStorageValue) (uuid.UUID, error) { //should do this part by one r with if?
	if authInfo.Role != EmployerStr {
		// log.Printf("Invalid action: %s", err)
		return uuid.UUID{}, errors.New("Invalid action")
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, err
	}

	var vacancyReg Vacancy
	err = json.Unmarshal(bytes, &vacancyReg)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, err
	}

	id := uuid.New()
	vacancyReg.ID = id
	vacancyReg.OwnerID = authInfo.ID
	ok := h.Storage.CreateVacancy(vacancyReg)

	if !ok {
		log.Printf("Error while creating vacancy: %s", err)
		return uuid.UUID{}, errors.New("Error while creating vacancy")
	}

	return id, nil
}

func (h *UserService) GetVacancy(vacancyId uuid.UUID) (Vacancy, error) {
	vacancy, err := h.Storage.GetVacancy(vacancyId)

	if err != nil { //error wrap
		return vacancy, errors.New("Error while getting vacancy")
	}

	return vacancy, nil
}

func (h *UserService) DeleteVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) error {
	if authInfo.Role != EmployerStr {
		return errors.New(ForbiddenMsg)
	}

	vacancy, err := h.Storage.GetVacancy(vacancyId)

	if vacancy.OwnerID != authInfo.ID || err != nil { //error wrap
		return errors.New(ForbiddenMsg)
	}

	err = h.Storage.DeleteVacancy(vacancyId)

	if err != nil {
		return errors.New("Error while deleting vacancy")
	}

	return nil
}

func (h *UserService) PutVacancy(vacancyId uuid.UUID, body io.ReadCloser, authInfo AuthStorageValue) error {
	if authInfo.Role != EmployerStr {
		// log.Printf("Invalid action: %s", err)
		return errors.New("Invalid action")
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "reading body error")
	}

	var vacancy Vacancy
	err = json.Unmarshal(bytes, &vacancy)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return err
	}

	ok := h.Storage.PutVacancy(vacancy, authInfo.ID, vacancyId)

	if !ok {
		log.Printf("Error while changing vacancy")
		return errors.New("Error while changing vacancy")
	}

	return nil
}

func (h *UserService) GetVacancies(params map[string]interface{}) ([]Vacancy, error) {
	return h.Storage.GetVacancies(params)
}
