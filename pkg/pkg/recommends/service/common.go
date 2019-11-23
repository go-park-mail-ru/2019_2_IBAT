package service

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"
	"2019_2_IBAT/pkg/pkg/recommends"

	"github.com/google/uuid"
)

type Service struct {
	Storage recommends.Repository
}

func (serv Service) SetTagIDs(AuthRec AuthStorageValue, tagIDs []uuid.UUID) error {
	return serv.Storage.SetTagIDs(AuthRec, tagIDs)
}

func (serv Service) GetTagIDs(AuthRec AuthStorageValue) ([]uuid.UUID, error) {
	return serv.Storage.GetTagIDs(AuthRec)
}
