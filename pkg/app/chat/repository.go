package chat

import (
	. "2019_2_IBAT/pkg/app/chat/models"
	. "2019_2_IBAT/pkg/pkg/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateChat(seekerId uuid.UUID, employerId uuid.UUID) (uuid.UUID, error)
	CreateMessage(msg InChatMessage) error

	GetChats(authInfo AuthStorageValue) ([]Chat, error)
	GetCompanionId(msg InChatMessage) (uuid.UUID, error)
	GetChatHistory(authInfo AuthStorageValue, chatId uuid.UUID) ([]OutChatMessage, error)
}