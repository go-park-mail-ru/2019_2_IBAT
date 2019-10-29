package handler

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
)

type UserEmplAnswer struct {
	Role     string   `json:"role"`
	Employer Employer `json:"profile"`
}

type UserSeekAnswer struct {
	Role   string `json:"role"`
	Seeker Seeker `json:"profile"`
}

type Error struct {
	Message string            `json:"error"`
	Params  map[string]string `json:"params"`
}
