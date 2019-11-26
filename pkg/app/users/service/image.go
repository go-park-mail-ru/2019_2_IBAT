package users

import (
	"github.com/google/uuid"
)

func (h UserService) SetImage(id uuid.UUID, class string, imageName string) bool {
	return h.Storage.SetImage(id, class, imageName)
}
