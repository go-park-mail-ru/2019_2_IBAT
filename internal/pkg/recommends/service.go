package recommends

import (
	. "2019_2_IBAT/internal/pkg/interfaces"

	"github.com/google/uuid"
)

type Service interface {
	SetTagIDs(AuthRec AuthStorageValue, tagIDs []uuid.UUID) error
	GetTagIDs(AuthRec AuthStorageValue) ([]uuid.UUID, error)

	// GetVacancyTagIDs() ([]uuid.UUID, error)

	// GetTagIDs(body io.ReadCloser) (uuid.UUID, error)

	// DeleteUser(authInfo AuthStorageValue) error
	// PutSeeker(body io.ReadCloser, id uuid.UUID) error
}
