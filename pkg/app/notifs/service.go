package notifs

import (
	"context"

	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	"google.golang.org/grpc"
)

// type Service interface {
// 	SendNotification(ctx context.Context, msg *notifsproto.SendNotificationMessage) (*notifsproto.Bool, error)
// }

type ServiceClient interface {
	SendNotification(ctx context.Context, in *notifsproto.SendNotificationMessage, opts ...grpc.CallOption) (*notifsproto.Bool, error)
}
