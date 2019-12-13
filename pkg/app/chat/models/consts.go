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
	ChatID uuid.UUID `json:"chat_id"              db:"chat_id"`
	// OwnerInfo AuthStorageValue `json:"-"                 db:"-"`
	OwnerId   uuid.UUID `json:"owner_id"          db:"owner_id"`
	Timestamp string    `json:"created_at"         db:"created_at"`
	Text      string    `json:"content"              db:"content"`
}

type Chat struct {
	ChatID   uuid.UUID `json:"chat_id"         db:"chat_id"`
	SeekerID uuid.UUID `json:"seeker_id"       db:"seeker_id"`
	Employer uuid.UUID `json:"employer_id"     db:"employer_id"`
}
