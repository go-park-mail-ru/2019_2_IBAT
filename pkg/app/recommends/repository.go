package recommends

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"
)

type Repository interface {
	SetTagIDs(AuthRec AuthStorageValue, tagIDs []string) error
	GetTagIDs(AuthRec AuthStorageValue) ([]string, error)
	GetUsersForTags([]string) ([]string, error)
}
