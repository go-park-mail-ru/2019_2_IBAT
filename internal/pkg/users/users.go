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

type UserService struct {
	Storage UserStorage
}

func (h *UserService) CreateSeeker(body io.ReadCloser) (uuid.UUID, error) {
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

func (h *UserService) CreateResume(body io.ReadCloser, cookie string, authStor AuthStorage) (uuid.UUID, error) { //should do this part by one r with if?
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		// err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, errors.New("Invalid body, transfer error")
	}

	var resumeReg Resume
	err = json.Unmarshal(bytes, &resumeReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	record, ok := authStor.Get(cookie)
	if !ok || record.Role != SeekerStr {
		// log.Printf("Invalid action: %s", err)
		return uuid.UUID{}, errors.New(ForbiddenMsg)
	}

	id, ok := h.Storage.CreateResume(resumeReg, record.ID)

	if !ok {
		// log.Printf("Error while creating resume: %s", err)
		return uuid.UUID{}, errors.New("Something went wrong. Error while creating resume")
	}

	return id, nil
}

func (h *UserService) DeleteResume(resumeId uuid.UUID, cookie string, authStor AuthStorage) error {
	record, ok := authStor.Get(cookie)
	if !ok || record.Role != SeekerStr {
		return errors.New(ForbiddenMsg)
	}

	resume, ok := h.Storage.GetResume(resumeId)

	if resume.OwnerID != record.ID || !ok {
		return errors.New(ForbiddenMsg)
	}

	ok = h.Storage.DeleteResume(resumeId)

	if !ok {
		return errors.New("Internal server error")
	}

	return nil
}

func (h *UserService) GetResume(resumeId uuid.UUID) (Resume, error) {
	resume, ok := h.Storage.GetResume(resumeId)

	if !ok {
		return resume, errors.New(InvalidIdMsg)
	}

	return resume, nil
}

func (h *UserService) PutResume(resumeId uuid.UUID, body io.ReadCloser,
	cookie string, authStor AuthStorage) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "reading body error")
	}

	var resume Resume
	err = json.Unmarshal(bytes, &resume)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return err
	}

	user, ok := authStor.Get(cookie)
	if !ok || user.Role != SeekerStr {
		// log.Printf(ForbiddenMsg, err)
		return errors.New(ForbiddenMsg)
	}

	ok = h.Storage.PutResume(resume, user.ID, resumeId)

	if !ok {
		// log.Printf("Error while creating resume: %s", err)
		return errors.New(InternalErrorMsg)
	}

	return nil
}

func (h *UserService) DeleteUser(cookie string, authStor AuthStorage) error {
	record, ok := authStor.Get(cookie)
	if !ok {
		return errors.New(ForbiddenMsg)
	}

	if record.Role == SeekerStr {
		h.Storage.DeleteSeeker(record.ID)
	} else if record.Role == EmployerStr {
		h.Storage.DeleteEmployer(record.ID)
	}

	return nil
}

func (h *UserService) PutSeeker(body io.ReadCloser, id uuid.UUID) error {
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
		// log.Printf("Error while creating seeker: %s", err)
		return errors.New("Error while changing employer")
	}

	return nil
}

func (h *UserService) CreateVacancy(body io.ReadCloser, cookie string, authStor AuthStorage) (uuid.UUID, error) { //should do this part by one r with if?
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

	record, ok := authStor.Get(cookie)
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

func (h *UserService) GetVacancy(vacancyId uuid.UUID) (Vacancy, error) {
	vacancy, ok := h.Storage.GetVacancy(vacancyId)

	if !ok {
		return vacancy, errors.New("Error while getting vacancy")
	}

	return vacancy, nil
}

func (h *UserService) DeleteVacancy(vacancyId uuid.UUID, cookie string, authStor AuthStorage) error {
	record, ok := authStor.Get(cookie)
	if !ok || record.Role != EmployerStr {
		return errors.New(ForbiddenMsg)
	}

	vacancy, ok := h.Storage.GetVacancy(vacancyId)

	if vacancy.OwnerID != record.ID || !ok {
		return errors.New(ForbiddenMsg)
	}

	ok = h.Storage.DeleteVacancy(vacancyId)

	if !ok {
		return errors.New("Error while deleting vacancy")
	}

	return nil
}

func (h *UserService) PutVacancy(vacancyId uuid.UUID, body io.ReadCloser,
	cookie string, authStor AuthStorage) error {
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

	user, ok := authStor.Get(cookie)
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
