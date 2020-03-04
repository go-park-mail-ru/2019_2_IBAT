package recommends

import (
	"context"

	"2019_2_IBAT/pkg/app/recommends/recomsproto"

	"google.golang.org/grpc"
)

// type Service interface {
// 	SetTagIDs(ctx context.Context, record *recomsproto.SetTagIDsMessage) (*recomsproto.Bool, error)
// 	GetTagIDs(ctx context.Context, authInfo *recomsproto.GetTagIDsMessage) (*recomsproto.IDsMessage, error)
// 	GetUsersForTags(ctx context.Context, userIds *recomsproto.IDsMessage) (*recomsproto.IDsMessage, error)
// }

type ServiceClient interface {
	SetTagIDs(ctx context.Context, in *recomsproto.SetTagIDsMessage, opts ...grpc.CallOption) (*recomsproto.Bool, error)
	GetTagIDs(ctx context.Context, in *recomsproto.GetTagIDsMessage, opts ...grpc.CallOption) (*recomsproto.IDsMessage, error)
	GetUsersForTags(ctx context.Context, in *recomsproto.IDsMessage, opts ...grpc.CallOption) (*recomsproto.IDsMessage, error)
}
