package service

import (
	"2019_2_IBAT/pkg/app/auth/session"
	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	. "2019_2_IBAT/pkg/pkg/models"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Service struct {
	NotifChan    chan NotifStruct
	ConnectsPool WsConnects
	AuthService  session.ServiceClient
	RecomService recomsproto.ServiceClient
}

func (h Service) SendNotification(ctx context.Context,
	msg *notifsproto.SendNotificationMessage) (*notifsproto.Bool, error) {
	notif := NotifStruct{
		VacancyId: uuid.MustParse(msg.VacancyID),
		TagIDs:    StringsToUuids(msg.TagIDs),
	}

	h.NotifChan <- notif
	return &notifsproto.Bool{}, nil
}

func (h Service) Notifications() {
	for {
		notif := <-h.NotifChan
		fmt.Println("Notification accepted")
		fmt.Println(notif)

		ctx := context.Background()
		idsMsg, err := h.RecomService.GetUsersForTags(ctx,
			&recomsproto.IDsMessage{IDs: UuidsToStrings(notif.TagIDs)},
		)

		ids := StringsToUuids(idsMsg.IDs)

		fmt.Println("Users ids intrested in new vacancy")
		fmt.Println(ids)

		if err != nil {
			log.Printf("Notifications %s", err)
		}
		fmt.Println("connects.ConsMu.Lock()")

		h.ConnectsPool.ConsMu.Lock()
		for _, id := range ids {
			fmt.Println("Notification ready  to be sent to user")
			fmt.Println(h.ConnectsPool.Connects[id])
			if cons, ok := h.ConnectsPool.Connects[id]; ok {
				fmt.Printf("Notification was sent to user %s\n", id.String())
				for _, con := range cons.Connects {
					con.Ch <- notif.VacancyId
				}
			}
		}
		h.ConnectsPool.ConsMu.Unlock()
		fmt.Println("connects.ConsMu.Unlock()")

		fmt.Println(ids)
	}
}

// func (s Service) ProcessMessage() {
// 	for {
// 		msg := <-s.MainChan
// 		fmt.Printf("ProcessMessage msg %s was read from main channel\n", msg.Text)
// 		id, name, err := s.Storage.GetCompanionIdAndName(msg)
// 		fmt.Printf("ProcessMessage companion id %s was accepted\n", id.String())
// 		fmt.Println(id)
// 		fmt.Println(err)

// 		// msg.ChatID =
// 		msg.Timestamp = time.Now().In(Loc).Format(TimeFormat)

// 		outMsg := OutChatMessage{
// 			ChatID:     msg.ChatID,
// 			OwnerId:    id,
// 			OwnerName:  name,
// 			Text:       msg.Text,
// 			IsNotYours: true,
// 			Timestamp:  msg.Timestamp,
// 		}

// 		s.ConnectsPool.ConsMu.Lock()
// 		if cons, ok := s.ConnectsPool.Connects[id]; ok {
// 			s.ConnectsPool.Connects[id].Mu.Lock()
// 			// cons := s.ConnectsPool.Connects[id]
// 			for _, con := range cons.Connects {
// 				fmt.Printf("ProcessMessage msg %s was sent to output channel\n", msg.Text)
// 				con.Ch <- outMsg
// 			}
// 			s.ConnectsPool.Connects[id].Mu.Unlock()
// 		}
// 		s.ConnectsPool.ConsMu.Unlock()

// 		// if not ok set unread
// 		go s.Storage.CreateMessage(msg)
// 	}
// }
