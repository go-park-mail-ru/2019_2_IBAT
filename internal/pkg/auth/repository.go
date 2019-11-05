package auth

import (
	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
)

type Repository interface {
	Get(cookie string) (AuthStorageValue, bool)
	Set(id uuid.UUID, class string) (AuthStorageValue, string, error)
	Delete(cookie string) bool
}
