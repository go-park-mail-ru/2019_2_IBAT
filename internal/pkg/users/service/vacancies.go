package users

import (
	"encoding/json"

	"io"
	"io/ioutil"
	"log"

	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h UserService) CreateVacancy(body io.ReadCloser, cookie string, authStor auth.Service) (uuid.UUID, error) { //should do this part by one r with if?
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

	record, ok := authStor.GetSession(cookie)
	if !ok || record.Role != EmployerStr {
		log.Printf("Invalid action: %s", err)
		return uuid.UUID{}, errors.New("Invalid action")
	}

	id, ok := h.Storage.CreateVacancy(vacancyReg, record.ID)

	if !ok {
		log.Printf("Error while creating vacancy: %s", err)
		return uuid.UUID{}, errors.New("Error while creating vacancy")
	}

	return id, nil
}

func (h UserService) GetVacancy(vacancyId uuid.UUID) (Vacancy, error) {
	vacancy, err := h.Storage.GetVacancy(vacancyId)

	if err != nil { //error wrap
		return vacancy, errors.New("Error while getting vacancy")
	}

	return vacancy, nil
}

func (h UserService) DeleteVacancy(vacancyId uuid.UUID, cookie string, authStor auth.Service) error {
	record, ok := authStor.GetSession(cookie)
	if !ok || record.Role != EmployerStr {
		return errors.New(ForbiddenMsg)
	}

	vacancy, err := h.Storage.GetVacancy(vacancyId)

	if vacancy.OwnerID != record.ID || err != nil { //error wrap
		return errors.New(ForbiddenMsg)
	}

	err = h.Storage.DeleteVacancy(vacancyId)

	if err != nil {
		return errors.New("Error while deleting vacancy")
	}

	return nil
}

func (h UserService) PutVacancy(vacancyId uuid.UUID, body io.ReadCloser,
	cookie string, authStor auth.Service) error {
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

	user, ok := authStor.GetSession(cookie)
	if !ok || user.Role != EmployerStr {
		log.Printf("Invalid action: %s", err)
		return errors.New("Invalid action")
	}

	ok = h.Storage.PutVacancy(vacancy, user.ID, vacancyId)

	if !ok {
		log.Printf("Error while changing vacancy")
		return errors.New("Error while changing vacancy")
	}

	return nil
}

func (h UserService) GetVacancies() ([]Vacancy, error) {
	return h.Storage.GetVacancies()
}
