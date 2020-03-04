package service

import (
	"sync"

	"github.com/google/uuid"
)

type ConnectsPerUser struct {
	Connects []*Connect
	Mu       *sync.Mutex
	// Ch       chan uuid.UUID
}

type WsConnects struct {
	ConsMu   *sync.Mutex
	Connects map[uuid.UUID]*ConnectsPerUser
}
