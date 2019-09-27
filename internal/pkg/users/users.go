package users

import (
	"encoding/json"

	"io"
	"io/ioutil"
	"log"

	. "hh_workspace/2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Controler struct {
	Storage UserStorage
}

func (h *Controler) HandleCreateSeeker(body io.ReadCloser) (uuid.UUID, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, err
	}

	var newSeekerReg SeekerReg
	err = json.Unmarshal(bytes, &newSeekerReg)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, err
	}

	id, ok := h.Storage.CreateSeeker(newSeekerReg)
	if !ok {
		log.Println("Here inside users")
		log.Printf("Error while creating seeker: %s", err)
		return uuid.UUID{}, errors.New("Error while creating seeker")
	}

	return id, nil
}

func (h *Controler) HandleCreateEmployer(body io.ReadCloser) (uuid.UUID, error) { //should do this part by one handler with if?
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, err
	}

	var newEmployerReg EmployerReg
	err = json.Unmarshal(bytes, &newEmployerReg)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, err
	}

	id, ok := h.Storage.CreateEmployer(newEmployerReg)

	if !ok {
		log.Printf("Error while creating employer: %s", err)
		return uuid.UUID{}, errors.New("Error while creating employer")
	}

	return id, nil
}

func (h *Controler) HandleCreateResume(body io.ReadCloser, cookie string, authStor AuthStorage) (uuid.UUID, error) { //should do this part by one handler with if?
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, err
	}

	var resumeReg Resume
	err = json.Unmarshal(bytes, &resumeReg)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, err
	}

	record, ok := authStor.Get(cookie)
	if !ok || record.Class != SeekerStr {
		log.Printf("Invalid action: %s", err)
		return uuid.UUID{}, errors.New("Invalid action")
	}

	id, ok := h.Storage.CreateResume(resumeReg, record.ID)

	if !ok {
		log.Printf("Error while creating resume: %s", err)
		return uuid.UUID{}, errors.New("Error while creating resume")
	}

	return id, nil
}

func (h *Controler) HandleDeleteResume(resumeId uuid.UUID, cookie string, authStor AuthStorage) error {
	record, ok := authStor.Get(cookie)
	if !ok || record.Class != SeekerStr {
		return errors.New("Invalid action")
	}

	resume, ok := h.Storage.GetResume(resumeId)

	if resume.OwnerID != record.ID || !ok {
		return errors.New("Error while deleting resume")
	}

	ok = h.Storage.DeleteResume(resumeId)

	if !ok {
		return errors.New("Error while deleting resume")
	}

	return nil
}

func (h *Controler) HandleGetResume(resumeId uuid.UUID, cookie string, authStor AuthStorage) (Resume, error) {
	resume, ok := h.Storage.GetResume(resumeId)

	if !ok {
		return resume, errors.New("Error while getting resume")
	}

	return resume, nil
}

func (h *Controler) HandlePutResume(resumeId uuid.UUID, body io.ReadCloser,
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
	if !ok || user.Class != SeekerStr {
		log.Printf("Invalid action: %s", err)
		return errors.New("Invalid action")
	}

	ok = h.Storage.PutResume(resume, user.ID, resumeId)

	if !ok {
		log.Printf("Error while creating resume: %s", err)
		return errors.New("Error while changing resume")
	}

	return nil
}

func (h *Controler) HandleGetSeeker(cookie string, authStor AuthStorage) (Seeker, error) {

	record, ok := authStor.Get(cookie)
	if !ok {
		return Seeker{}, errors.New("Invalid action")
	}

	res, _ := h.Storage.GetSeeker(record.ID)

	return res, nil
}

func (h *Controler) HandleGetEmployer(cookie string, authStor AuthStorage) (Employer, error) {

	record, ok := authStor.Get(cookie)
	if !ok {
		return Employer{}, errors.New("Invalid action")
	}

	res, _ := h.Storage.GetEmployer(record.ID)

	return res, nil
}

func (h *Controler) HandleDeleteUser(cookie string, authStor AuthStorage) error {
	record, ok := authStor.Get(cookie)
	if !ok {
		return errors.New("Invalid action")
	}

	if record.Class == SeekerStr {
		h.Storage.DeleteSeeker(record.ID)
	} else if record.Class == EmployerStr {
		h.Storage.DeleteEmployer(record.ID)
	}

	return nil
}

func (h *Controler) HandlePutSeeker(body io.ReadCloser, id uuid.UUID) error {
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

func (h *Controler) HandlePutEmployer(body io.ReadCloser, id uuid.UUID) error {
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

func (h *Controler) HandleCreateVacancy(body io.ReadCloser, cookie string, authStor AuthStorage) (uuid.UUID, error) { //should do this part by one handler with if?
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
	if !ok || record.Class != EmployerStr {
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

func (h *Controler) HandleGetVacancy(vacancyId uuid.UUID, cookie string, authStor AuthStorage) (Vacancy, error) {
	vacancy, ok := h.Storage.GetVacancy(vacancyId)

	if !ok {
		return vacancy, errors.New("Error while getting vacancy")
	}

	return vacancy, nil
}

func (h *Controler) HandleDeleteVacancy(vacancyId uuid.UUID, cookie string, authStor AuthStorage) error {
	record, ok := authStor.Get(cookie)
	if !ok || record.Class != EmployerStr {
		return errors.New("Invalid action")
	}

	vacancy, ok := h.Storage.GetVacancy(vacancyId)

	if vacancy.OwnerID != record.ID || !ok {
		return errors.New("Error while deleting vacancy")
	}

	ok = h.Storage.DeleteVacancy(vacancyId)

	if !ok {
		return errors.New("Error while deleting vacancy")
	}

	return nil
}

func (h *Controler) HandlePutVacancy(vacancyId uuid.UUID, body io.ReadCloser,
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
	if !ok || user.Class != EmployerStr {
		log.Printf("Invalid action: %s", err)
		return errors.New("Invalid action")
	}

	ok = h.Storage.PutVacancy(vacancy, user.ID, vacancyId)

	if !ok {
		log.Printf("Error while creating vacancy: %s", err)
		return errors.New("Error while changing vacancy")
	}

	return nil
}
