package service

import (
	"2019_2_IBAT/pkg/app/auth/session"
	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	. "2019_2_IBAT/pkg/pkg/models"
	"context"

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

// h.NotifChan <- NotifStruct{
// 	VacancyId: id,
// 	TagIDs:    tagIDs,
// }
