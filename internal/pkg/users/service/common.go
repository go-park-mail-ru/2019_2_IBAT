package users

import (
	"2019_2_IBAT/internal/pkg/users"

	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
)

type UserService struct {
	Storage users.Repository
}

func (h *UserService) DeleteUser(authInfo AuthStorageValue) error {

	err := h.Storage.DeleteUser(authInfo.ID)

	return err
}

func (h *UserService) CheckUser(email string, password string) (uuid.UUID, string, bool) {
	return h.Storage.CheckUser(email, password)
}
