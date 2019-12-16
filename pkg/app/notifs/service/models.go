package service

import (
	"sync"

	"github.com/google/uuid"
)

type NotifStruct struct {
	VacancyId uuid.UUID
	TagIDs    []uuid.UUID
}

type ConnectsPerUser struct {
	Connects []*Connect
	Mu       *sync.Mutex
}

type WsConnects struct {
	ConsMu   *sync.Mutex
	Connects map[uuid.UUID]*ConnectsPerUser
}
