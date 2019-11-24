package users

import (
	"2019_2_IBAT/pkg/pkg/users"
	"fmt"
	"log"

	. "2019_2_IBAT/pkg/pkg/interfaces"
	recServ "2019_2_IBAT/pkg/pkg/recommends/service"

	"github.com/google/uuid"
)

type UserService struct {
	Storage      users.Repository
	RecomService recServ.Service
	NotifChan    chan NotifStruct
}

func (h *UserService) DeleteUser(authInfo AuthStorageValue) error {

	err := h.Storage.DeleteUser(authInfo.ID)

	return err
}

func (h *UserService) CheckUser(email string, password string) (uuid.UUID, string, bool) {
	return h.Storage.CheckUser(email, password)
}

func (h *UserService) Notifications(connects *WsConnects) {
	// var notif NotifStruct
	for {
		notif := <-h.NotifChan
		fmt.Println("Notification accepted")
		fmt.Println(notif)

		ids, err := h.RecomService.GetUsersForTags(notif.TagIDs)
		fmt.Println("Users ids intrested in new vacancy")
		fmt.Println(ids)

		if err != nil {
			log.Printf("Notifications %s", err)
		}
		// connects.Mu.Lock()
		fmt.Println("connects.ConsMu.Lock()")

		connects.ConsMu.Lock()
		for _, id := range ids {
			fmt.Println("Notification ready  to be sent to user %s")
			fmt.Println(connects.Connects[id])
			if _, ok := connects.Connects[id]; ok {
				connects.Connects[id].Ch <- notif.VacancyId
				fmt.Printf("Notification was sent to user %s\n", id.String())
			}
			fmt.Printf("Notification can not be sent to user %s\n", id.String())
		}
		connects.ConsMu.Unlock()
		fmt.Println("connects.ConsMu.Unlock()")

		fmt.Println(ids)
	}
}
