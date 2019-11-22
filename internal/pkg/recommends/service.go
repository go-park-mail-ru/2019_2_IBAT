package recommends

import (
	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
)

type Service interface {
	SetTagIDs(AuthRec AuthStorageValue, tagIDs []uuid.UUID) error
}
