package users

import (
	"encoding/json"

	"io"
	"io/ioutil"

	. "hh_workspace/2019_2_IBAT/internal/pkg/interfaces" //what

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// type UserStorage interface {
// 	CreateEmployer(seekerInput EmployerReg) (uuid.UUID, bool)
// 	CreateSeeker(seekerInput SeekerReg) (uuid.UUID, bool)
// 	CreateResume(resumeReg Resume) (uuid.UUID, bool)

// 	DeleteEmployer(id uuid.UUID)
// 	DeleteSeeker(id uuid.UUID)

// 	CheckUser(email string, password string) (uuid.UUID, string, bool)

// 	GetSeekers() []Seeker
// }

type Controler struct {
	Storage UserStorage
}

func (h *Controler) HandleCreateSeeker(body io.ReadCloser) (uuid.UUID, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, err
	}

	var newSeekerReg SeekerReg
	err = json.Unmarshal(bytes, &newSeekerReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, err
	}

	id, ok := h.Storage.CreateSeeker(newSeekerReg)
	if !ok {
		// log.Println("Here inside users")
		// log.Printf("Error while creating seeker: %s", err)
		return uuid.UUID{}, errors.New("Error while creating seeker")
	}

	return id, nil
}

func (h *Controler) HandleCreateEmployer(body io.ReadCloser) (uuid.UUID, error) { //should do this part by one handler with if?
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, err
	}

	var newEmployerReg EmployerReg
	err = json.Unmarshal(bytes, &newEmployerReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, err
	}

	id, ok := h.Storage.CreateEmployer(newEmployerReg)

	if !ok {
		// log.Printf("Error while creating employer: %s", err)
		return uuid.UUID{}, errors.New("Error while creating employer")
	}

	return id, nil
}

func (h *Controler) HandleCreateResume(body io.ReadCloser, cookie string, authStor AuthStorage) (uuid.UUID, error) { //should do this part by one handler with if?
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, err
	}

	var resumeReg Resume
	err = json.Unmarshal(bytes, &resumeReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, err
	}

	record, ok := authStor.Get(cookie)
	if !ok || record.Class != SeekerStr {
		// log.Printf("Invalid action: %s", err)
		return uuid.UUID{}, errors.New("Invalid action")
	}

	id, ok := h.Storage.CreateResume(resumeReg, record.ID)

	if !ok {
		// log.Printf("Error while creating resume: %s", err)
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
		// log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return err
	}

	user, ok := authStor.Get(cookie)
	if !ok || user.Class != SeekerStr {
		// log.Printf("Invalid action: %s", err)
		return errors.New("Invalid action")
	}

	ok = h.Storage.PutResume(resume, user.ID, resumeId)

	if !ok {
		// log.Printf("Error while creating resume: %s", err)
		return errors.New("Error while changing resume")
	}

	return nil
}

func (h *Controler) HandleGetSeeker(cookie string, authStor AuthStorage) (Seeker, error) {

	record, ok := authStor.Get(cookie)
	if !ok {
		return Seeker{}, errors.New("Invalid action")
	}

	res := h.Storage.GetSeeker(record.ID)

	return res, nil
}

func (h *Controler) HandleGetEmployer(cookie string, authStor AuthStorage) (Employer, error) {

	record, ok := authStor.Get(cookie)
	if !ok {
		return Employer{}, errors.New("Invalid action")
	}

	res := h.Storage.GetEmployer(record.ID)

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

// func (h *Controler) HandleCreateVacancy(body io.ReadCloser, cookie string, authStor AuthStorage) (uuid.UUID, error) { //should do this part by one handler with if?
// 	bytes, err := ioutil.ReadAll(body)
// 	if err != nil {
// 		err = errors.Wrap(err, "reading body error")
// 		return uuid.UUID{}, err
// 	}

// 	var resumeReg Resume
// 	err = json.Unmarshal(bytes, &resumeReg)
// 	if err != nil {
// 		err = errors.Wrap(err, "unmarshaling error")
// 		return uuid.UUID{}, err
// 	}

// 	record, ok := authStor.Get(cookie)
// 	if !ok || record.Class != SeekerStr {
// 		log.Printf("Invalid action: %s", err)
// 		return uuid.UUID{}, errors.New("Invalid action")
// 	}

// 	id, ok := h.Storage.CreateResume(resumeReg, record.ID)

// 	if !ok {
// 		log.Printf("Error while creating resume: %s", err)
// 		return uuid.UUID{}, errors.New("Error while creating resume")
// 	}

// 	return id, nil
// }
