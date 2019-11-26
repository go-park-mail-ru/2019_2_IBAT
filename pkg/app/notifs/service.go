package notifs

import (
	"context"

	"2019_2_IBAT/pkg/app/notifs/notifsproto"
)

type Service interface {
	SendNotification(ctx context.Context, msg *notifsproto.SendNotificationMessage) (*notifsproto.Bool, error)
}
