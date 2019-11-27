package users

import (
	"github.com/google/uuid"

	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	"2019_2_IBAT/pkg/app/users"
	. "2019_2_IBAT/pkg/pkg/models"
)

type UserService struct {
	Storage      users.Repository
	RecomService recomsproto.ServiceClient
	NotifService notifsproto.ServiceClient
}

func (h *UserService) DeleteUser(authInfo AuthStorageValue) error {

	err := h.Storage.DeleteUser(authInfo.ID)

	return err
}

func (h *UserService) CheckUser(email string, password string) (uuid.UUID, string, bool) {
	return h.Storage.CheckUser(email, password)
}
