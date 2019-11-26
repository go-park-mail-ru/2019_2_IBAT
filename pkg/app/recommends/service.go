package recommends

import (
	"context"

	"2019_2_IBAT/pkg/app/recommends/recomsproto"
)

type Service interface {
	SetTagIDs(ctx context.Context, record *recomsproto.SetTagIDsMessage) (*recomsproto.Bool, error)
	GetTagIDs(ctx context.Context, authInfo *recomsproto.GetTagIDsMessage) (*recomsproto.IDsMessage, error)
	GetUsersForTags(ctx context.Context, userIds *recomsproto.IDsMessage) (*recomsproto.IDsMessage, error)
}
