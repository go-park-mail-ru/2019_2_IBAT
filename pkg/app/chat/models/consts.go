package models

import (
	. "2019_2_IBAT/pkg/pkg/models"

	"github.com/google/uuid"
)

type InChatMessage struct {
	ChatID    uuid.UUID        `json:"chat_id"                 db:"id"`
	OwnerInfo AuthStorageValue `json:"-"                 db:"-"`
	Timestamp string           `json:"timestamp"                 db:"id"`
	Text      string           `json:"text"                 db:"id"`
}

type OutChatMessage struct {
	ChatID     uuid.UUID `json:"chat_id"`
	OwnerId    uuid.UUID `json:"owner_id"`
	OwnerName  string    `json:"owner_name"`
	Timestamp  string    `json:"created_at"`
	Text       string    `json:"content"`
	IsNotYours bool      `json:"is_not_yours"`
}

type Chat struct {
	ChatID        uuid.UUID `json:"chat_id"`
	CompanionName string    `json:"companion_name"`
	CompanionID   uuid.UUID `json:"companion_id"`
}
