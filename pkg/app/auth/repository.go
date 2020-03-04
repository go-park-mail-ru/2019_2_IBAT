package auth

import (
	"github.com/google/uuid"

	. "2019_2_IBAT/pkg/pkg/models"
)

type Repository interface {
	Get(cookie string) (AuthStorageValue, bool)
	Set(id uuid.UUID, class string) (AuthStorageValue, string, error)
	Delete(cookie string) bool
}
