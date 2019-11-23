package auth

import (
	"2019_2_IBAT/pkg/pkg/auth/session"

	"context"
	// auth_serv "2019_2_IBAT/internal/pkg/users/service"
)

const CookieName = "session-id"

type Service interface {
	CreateSession(ctx context.Context, sessInfo *session.Session) (*session.CreateSessionInfo, error)
	DeleteSession(ctx context.Context, cookie *session.Cookie) (*session.Bool, error)
	GetSession(ctx context.Context, cookie *session.Cookie) (*session.GetSessionInfo, error)
}
