package users

import (
	"github.com/google/uuid"

	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	"2019_2_IBAT/pkg/app/users"
	. "2019_2_IBAT/pkg/pkg/interfaces"
)

type UserService struct {
	Storage users.Repository
	// RecomService recServ.Service
	RecomService recomsproto.ServiceClient
	NotifService notifsproto.ServiceClient
	// NotifChan    chan NotifStruct
}

func (h *UserService) DeleteUser(authInfo AuthStorageValue) error {

	err := h.Storage.DeleteUser(authInfo.ID)

	return err
}

func (h *UserService) CheckUser(email string, password string) (uuid.UUID, string, bool) {
	return h.Storage.CheckUser(email, password)
}

// func (h *UserService) Notifications(connects *WsConnects) {
// 	// var notif NotifStruct
// 	for {
// 		notif := <-h.NotifChan
// 		fmt.Println("Notification accepted")
// 		fmt.Println(notif)

// 		ctx := context.Background()
// 		idsMsg, err := h.RecomService.GetUsersForTags(ctx,
// 			&recomsproto.IDsMessage{IDs: UuidsToStrings(notif.TagIDs)},
// 		)

// 		ids := StringsToUuids(idsMsg.IDs)

// 		fmt.Println("Users ids intrested in new vacancy")
// 		fmt.Println(ids)

// 		if err != nil {
// 			log.Printf("Notifications %s", err)
// 		}
// 		// connects.Mu.Lock()
// 		fmt.Println("connects.ConsMu.Lock()")

// 		connects.ConsMu.Lock()
// 		for _, id := range ids {
// 			fmt.Println("Notification ready  to be sent to user %s")
// 			fmt.Println(connects.Connects[id])
// 			if _, ok := connects.Connects[id]; ok {
// 				connects.Connects[id].Ch <- notif.VacancyId
// 				fmt.Printf("Notification was sent to user %s\n", id.String())
// 			} else {
// 				fmt.Printf("Notification can not be sent to user %s\n", id.String())
// 				delete(connects.Connects, id)
// 				fmt.Printf("Connections to user were removed\nd")
// 			}
// 		}
// 		connects.ConsMu.Unlock()
// 		fmt.Println("connects.ConsMu.Unlock()")

// 		fmt.Println(ids)
// 	}
// }

// func UuidsToStrings(ids []uuid.UUID) []string {
// 	var strIDs []string
// 	for _, id := range ids {
// 		strIDs = append(strIDs, id.String())
// 	}
// 	return strIDs
// }

// func StringsToUuids(strIDs []string) []uuid.UUID {
// 	var ids []uuid.UUID
// 	for _, id := range strIDs {
// 		ids = append(ids, uuid.MustParse(id))
// 	}
// 	return ids
// }
