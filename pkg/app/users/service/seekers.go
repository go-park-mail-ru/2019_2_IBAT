package users

import (
	"io"
	"io/ioutil"

	. "2019_2_IBAT/pkg/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *UserService) CreateSeeker(body io.ReadCloser) (uuid.UUID, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	var newSeekerReg Seeker
	err = newSeekerReg.UnmarshalJSON(bytes)
	if err != nil {
		return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	id := uuid.New()
	newSeekerReg.ID = id
	newSeekerReg.PathToImg = DefaultImg
	ok := h.Storage.CreateSeeker(newSeekerReg)
	if !ok {
		return uuid.UUID{}, errors.New(EmailExistsMsg)
	}

	return id, nil
}

func (h *UserService) PutSeeker(body io.ReadCloser, id uuid.UUID) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.New(BadRequestMsg)
	}

	var newSeekerReg SeekerReg
	err = newSeekerReg.UnmarshalJSON(bytes)
	if err != nil {
		return errors.New(InvalidJSONMsg)
	}

	ok := h.Storage.PutSeeker(newSeekerReg, id)
	if !ok {
		// log.Println("Here inside users")
		// log.Printf("Error while creating seeker: %s", err)
		return errors.New(BadRequestMsg)
	}

	return nil
}

func (h *UserService) GetSeeker(id uuid.UUID) (Seeker, error) {
	// log.Println("service.GetSeeker")
	return h.Storage.GetSeeker(id)
}

func (h *UserService) GetSeekers() ([]Seeker, error) {
	return h.Storage.GetSeekers()
}
