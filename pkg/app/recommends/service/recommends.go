package service

import (
	"context"

	"github.com/google/uuid"

	"2019_2_IBAT/pkg/app/recommends"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	. "2019_2_IBAT/pkg/pkg/models"
)

type Service struct {
	Storage recommends.Repository
}

func (serv Service) SetTagIDs(ctx context.Context, record *recomsproto.SetTagIDsMessage) (*recomsproto.Bool, error) {
	return &recomsproto.Bool{
			Ok: true,
		},
		serv.Storage.SetTagIDs(
			AuthStorageValue{
				ID:      uuid.MustParse(record.ID),
				Role:    record.Role,
				Expires: record.Expires,
			},
			record.IDs,
		)
}

func (serv Service) GetTagIDs(ctx context.Context, authInfo *recomsproto.GetTagIDsMessage) (*recomsproto.IDsMessage, error) {
	strIDs, err := serv.Storage.GetTagIDs(
		AuthStorageValue{
			ID:      uuid.MustParse(authInfo.ID),
			Role:    authInfo.Role,
			Expires: authInfo.Expires,
		})
	return &recomsproto.IDsMessage{
			IDs: strIDs,
		},
		err

}

func (serv Service) GetUsersForTags(ctx context.Context, userIds *recomsproto.IDsMessage) (*recomsproto.IDsMessage, error) {
	strIDs, err := serv.Storage.GetUsersForTags(userIds.IDs)

	return &recomsproto.IDsMessage{
			IDs: strIDs,
		},
		err
}
