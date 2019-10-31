package users

import (
	"encoding/json"

	"io"
	"io/ioutil"

	"2019_2_IBAT/internal/pkg/auth"
	"2019_2_IBAT/internal/pkg/users"

	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserService struct {
	Storage users.Repository
}

func (h UserService) CreateSeeker(body io.ReadCloser) (uuid.UUID, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		// err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, errors.New("Invalid body, transfer error")
	}

	var newSeekerReg SeekerReg
	err = json.Unmarshal(bytes, &newSeekerReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, errors.New("Invalid JSON")
	}

	id, ok := h.Storage.CreateSeeker(newSeekerReg)
	if !ok {
		// log.Println("Here inside users")
		// log.Printf("Error while creating seeker: %s", err)
		return uuid.UUID{}, errors.New("Email already exists")
	}

	return id, nil
}

func (h UserService) CreateEmployer(body io.ReadCloser) (uuid.UUID, error) { //should do this part by one r with if?
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

func (h UserService) DeleteUser(cookie string, authStor auth.Service) error {
	record, ok := authStor.GetSession(cookie)
	if !ok {
		return errors.New(ForbiddenMsg)
	}

	h.Storage.DeleteUser(record.ID)

	return nil
}

func (h UserService) PutSeeker(body io.ReadCloser, id uuid.UUID) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		return errors.Wrap(err, "reading body error")
	}

	var newSeekerReg SeekerReg
	err = json.Unmarshal(bytes, &newSeekerReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		return errors.Wrap(err, "unmarshaling error")
	}

	ok := h.Storage.PutSeeker(newSeekerReg, id)
	if !ok {
		// log.Println("Here inside users")
		// log.Printf("Error while creating seeker: %s", err)
		return errors.New("Error while changing seeker")
	}

	return nil
}

func (h UserService) PutEmployer(body io.ReadCloser, id uuid.UUID) error {
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

func (h UserService) GetSeeker(id uuid.UUID) (Seeker, error) {
	return h.Storage.GetSeeker(id)
}

func (h UserService) GetEmployer(id uuid.UUID) (Employer, error) {
	return h.Storage.GetEmployer(id)
}

func (h UserService) GetEmployers() ([]Employer, error) {
	return h.Storage.GetEmployers()
}

func (h UserService) GetSeekers() ([]Seeker, error) {
	return h.Storage.GetSeekers()
}

func (h UserService) CheckUser(email string, password string) (uuid.UUID, string, bool) {
	return h.Storage.CheckUser(email, password)
}
