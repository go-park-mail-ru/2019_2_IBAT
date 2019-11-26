package auth

import (
	"context"

	"2019_2_IBAT/pkg/app/auth/session"
)

const CookieName = "session-id"

type Service interface {
	CreateSession(ctx context.Context, sessInfo *session.Session) (*session.CreateSessionInfo, error)
	DeleteSession(ctx context.Context, cookie *session.Cookie) (*session.Bool, error)
	GetSession(ctx context.Context, cookie *session.Cookie) (*session.GetSessionInfo, error)
}
