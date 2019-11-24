package recommends

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"

	"github.com/google/uuid"
)

type Repository interface {
	SetTagIDs(AuthRec AuthStorageValue, tagIDs []uuid.UUID) error
	GetTagIDs(AuthRec AuthStorageValue) ([]uuid.UUID, error)
}