package users

import (
	"io"
	"io/ioutil"
	"log"

	. "2019_2_IBAT/pkg/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *UserService) CreateResume(body io.ReadCloser, authInfo AuthStorageValue) (uuid.UUID, error) {
	if authInfo.Role != SeekerStr {
		return uuid.UUID{}, errors.New(ForbiddenMsg)
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	var resumeReg Resume
	err = resumeReg.UnmarshalJSON(bytes)
	if err != nil {
		return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	id := uuid.New()
	resumeReg.ID = id
	resumeReg.OwnerID = authInfo.ID
	ok := h.Storage.CreateResume(resumeReg)

	if !ok {
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
	err = resume.UnmarshalJSON(bytes)
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
	if params["id"] != nil {
		return h.Storage.GetResumesByIDs(authInfo, params)
	} else {
		return h.Storage.GetResumes(authInfo, params)
	}
}
