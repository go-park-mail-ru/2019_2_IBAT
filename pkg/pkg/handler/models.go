package handler

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"
)

type UserEmplAnswer struct {
	Role     string   `json:"role"`
	Employer Employer `json:"profile"`
}

type UserSeekAnswer struct {
	Role   string `json:"role"`
	Seeker Seeker `json:"profile"`
}
