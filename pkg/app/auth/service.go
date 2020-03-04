package auth

import (
	"context"

	"2019_2_IBAT/pkg/app/auth/session"

	"google.golang.org/grpc"
)

const CookieName = "session-id"

// type Service interface {
// 	CreateSession(ctx context.Context, sessInfo *session.Session) (*session.CreateSessionInfo, error)
// 	DeleteSession(ctx context.Context, cookie *session.Cookie) (*session.Bool, error)
// 	GetSession(ctx context.Context, cookie *session.Cookie) (*session.GetSessionInfo, error)
// }

type ServiceClient interface {
	CreateSession(ctx context.Context, in *session.Session, opts ...grpc.CallOption) (*session.CreateSessionInfo, error)
	DeleteSession(ctx context.Context, in *session.Cookie, opts ...grpc.CallOption) (*session.Bool, error)
	GetSession(ctx context.Context, in *session.Cookie, opts ...grpc.CallOption) (*session.GetSessionInfo, error)
}
