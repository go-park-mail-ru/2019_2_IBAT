package service

import (
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"context"

	"log"

	"2019_2_IBAT/internal/pkg/auth/session"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AuthService struct {
	Storage auth.Repository
}

const CookieName = "session-id"

type Service interface {
	CreateSession(ctx context.Context, sessInfo *session.Session) (*session.CreateSessionInfo, error)
	DeleteSession(ctx context.Context, cookie *session.Cookie) (*session.Bool, error)
	GetSession(ctx context.Context, cookie *session.Cookie) (*session.GetSessionInfo, error)
}

func (h AuthService) CreateSession(ctx context.Context, sessInfo *session.Session) (*session.CreateSessionInfo, error) {

	authInfo, cookieValue, err := h.Storage.Set(uuid.MustParse(sessInfo.Id), sessInfo.Class)

	if err != nil {
		log.Printf("Error while unmarshaling: %s\n", err)
		err = errors.Wrap(err, "error while unmarshaling")
		return &session.CreateSessionInfo{}, errors.New("Creating session error")
	}

	return &session.CreateSessionInfo{
			ID:      authInfo.ID.String(),
			Role:    authInfo.Role,
			Expires: authInfo.Expires,
			Cookie:  cookieValue,
		},
		nil
}

func (h AuthService) DeleteSession(ctx context.Context, cookie *session.Cookie) (*session.Bool, error) {
	_, ok := h.Storage.Get(cookie.Cookie)
	if !ok {
		log.Printf("No such session")
		return &session.Bool{Ok: false}, nil
	}

	ok = h.Storage.Delete(cookie.Cookie)
	if !ok {
		return &session.Bool{Ok: false}, nil
	}

	return &session.Bool{Ok: true}, nil
}

func (auth AuthService) GetSession(ctx context.Context, cookie *session.Cookie) (*session.GetSessionInfo, error) {
	authInfo, ok := auth.Storage.Get(cookie.Cookie)
	return &session.GetSessionInfo{
		ID:      authInfo.ID.String(),
		Role:    authInfo.Role,
		Expires: authInfo.Expires,
		Ok:      ok,
	}, nil
}
