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

func (h *UserService) CreateResume(body io.ReadCloser, authInfo AuthStorageValue) (uuid.UUID, error) { //should do this part by one r with if?
	if authInfo.Role != SeekerStr {
		// log.Printf("Invalid action: %s", err)
		return uuid.UUID{}, errors.New(ForbiddenMsg)
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		// log.Printf("error while reading body: %s", err)
		// err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	var resumeReg Resume
	// id := uuid.New()
	// resumeReg.ID = id
	// resumeReg.OwnerID = authInfo.ID
	err = json.Unmarshal(bytes, &resumeReg)
	if err != nil {
		// log.Printf("Error while unmarshaling: %s", err)
		// err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	id := uuid.New()
	resumeReg.ID = id
	resumeReg.OwnerID = authInfo.ID
	ok := h.Storage.CreateResume(resumeReg)

	if !ok {
		// log.Printf("Error while creating resume: %s", err)
		return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	return id, nil
}

func (h *UserService) DeleteResume(resumeId uuid.UUID, authInfo AuthStorageValue) error {
	if authInfo.Role != SeekerStr {
		return errors.New(ForbiddenMsg)
	}

	resume, err := h.Storage.GetResume(resumeId)

	if resume.OwnerID != authInfo.ID || err != nil {
		return errors.New(ForbiddenMsg)
	}

	err = h.Storage.DeleteResume(resumeId)

	if err != nil {
		return errors.New(InternalErrorMsg)
	}

	return nil
}

func (h *UserService) GetResume(resumeId uuid.UUID) (Resume, error) {
	resume, err := h.Storage.GetResume(resumeId)

	if err != nil {
		log.Println("Service GetResume: failed to get resume")
		return resume, errors.New(InvalidIdMsg)
	}

	return resume, nil
}

func (h *UserService) PutResume(resumeId uuid.UUID, body io.ReadCloser, authInfo AuthStorageValue) error {
	if authInfo.Role != SeekerStr {
		// log.Printf(ForbiddenMsg, err)
		return errors.New(ForbiddenMsg)
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.New(BadRequestMsg)
	}

	var resume Resume
	err = json.Unmarshal(bytes, &resume)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		return errors.New(InvalidJSONMsg)
	}

	ok := h.Storage.PutResume(resume, authInfo.ID, resumeId)

	if !ok {
		// log.Printf("Error while creating resume: %s", err)
		return errors.New(InternalErrorMsg)
	}

	return nil
}

func (h *UserService) GetResumes(authInfo AuthStorageValue, params map[string]interface{}) ([]Resume, error) {
	return h.Storage.GetResumes(authInfo, params)
}
