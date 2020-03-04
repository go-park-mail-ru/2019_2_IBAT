package service

import (
	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	// . "2019_2_IBAT/pkg/pkg/models"
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestUserService_SendNotification(t *testing.T) {
	h := Service{
		NotifChan: make(chan NotifStruct, 5),
	}

	ctx := context.Background()
	msg := notifsproto.SendNotificationMessage{
		VacancyID: uuid.New().String(),
		TagIDs: []string{
			uuid.New().String(),
			uuid.New().String(),
		},
	}
	_, err := h.SendNotification(ctx, &msg)
	if err != nil {
		t.Errorf("Unexpected error %s\n", err.Error())
		return
	}
}
