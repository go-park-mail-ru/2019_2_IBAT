package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/auth/session"
	. "2019_2_IBAT/pkg/pkg/models"
)

type AuthService struct {
	Storage auth.Repository
}

const CookieName = "session-id"

func (h AuthService) CreateSession(ctx context.Context, sessInfo *session.Session) (*session.CreateSessionInfo, error) {

	authInfo, cookieValue, err := h.Storage.Set(uuid.MustParse(sessInfo.Id), sessInfo.Class)

	if err != nil {
		fmt.Printf("Error while unmarshaling: %s\n", err)
		err = errors.Wrap(err, "error while unmarshaling")
		return &session.CreateSessionInfo{}, errors.New(BadRequestMsg)
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
	ok := h.Storage.Delete(cookie.Cookie)
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
