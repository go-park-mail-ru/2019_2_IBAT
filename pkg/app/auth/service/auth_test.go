package service

import (
	mock_auth_rep "2019_2_IBAT/pkg/app/auth/service/mock_auth"
	"2019_2_IBAT/pkg/app/auth/session"
	. "2019_2_IBAT/pkg/pkg/models"

	"time"

	"context"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateSession(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthRepo := mock_auth_rep.NewMockRepository(mockCtrl)

	h := AuthService{
		Storage: mockAuthRepo,
	}

	tests := []struct {
		name             string
		class            string
		userId           uuid.UUID
		wantFail         bool
		wantInvJSON      bool
		wantErrorMessage string
		invJSON          string
		ctx              context.Context
	}{
		{
			name:   "Test1",
			userId: uuid.New(),
			ctx:    context.Background(),
			class:  SeekerStr,
		},
		{
			name:             "Test2",
			userId:           uuid.New(),
			ctx:              context.Background(),
			class:            EmployerStr,
			wantFail:         true,
			wantErrorMessage: BadRequestMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockAuthRepo.
					EXPECT().
					Set(tt.userId, tt.class).
					Return(AuthStorageValue{
						ID:      tt.userId,
						Role:    tt.class,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					}, "cookie", nil)
			} else {
				mockAuthRepo.
					EXPECT().
					Set(tt.userId, tt.class).
					Return(AuthStorageValue{
						ID:      tt.userId,
						Role:    tt.class,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					}, "cookie", errors.New(BadRequestMsg))
			}

			sess := session.Session{
				Id:    tt.userId.String(),
				Class: tt.class,
			}
			_, err := h.CreateSession(tt.ctx, &sess)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_GetSession(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthRepo := mock_auth_rep.NewMockRepository(mockCtrl)

	h := AuthService{
		Storage: mockAuthRepo,
	}

	tests := []struct {
		name             string
		class            string
		userId           uuid.UUID
		wantFail         bool
		wantInvJSON      bool
		wantErrorMessage string
		invJSON          string
		ctx              context.Context
		cookie           string
		ok               bool
	}{
		{
			name:   "Test1",
			userId: uuid.New(),
			ctx:    context.Background(),
			class:  SeekerStr,
			cookie: "Cookie",
			ok:     true,
		},
		{
			name:   "Test1",
			userId: uuid.New(),
			ctx:    context.Background(),
			class:  SeekerStr,
			cookie: "Cookie",
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockAuthRepo.
					EXPECT().
					Get(tt.cookie).
					Return(AuthStorageValue{
						ID:      uuid.New(),
						Role:    tt.class,
						Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
					}, tt.ok)
			}

			sess := session.Cookie{
				Cookie: tt.cookie,
			}
			_, err := h.GetSession(tt.ctx, &sess)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

// func (h AuthService) DeleteSession(ctx context.Context, cookie *session.Cookie) (*session.Bool, error) {
// 	_, ok := h.Storage.Get(cookie.Cookie)
// 	if !ok {
// 		log.Printf("No such session")
// 		return &session.Bool{Ok: false}, nil
// 	}

// 	ok = h.Storage.Delete(cookie.Cookie)
// 	if !ok {
// 		return &session.Bool{Ok: false}, nil
// 	}

// 	return &session.Bool{Ok: true}, nil
// }

func TestUserService_DeleteSession(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthRepo := mock_auth_rep.NewMockRepository(mockCtrl)

	h := AuthService{
		Storage: mockAuthRepo,
	}

	tests := []struct {
		name             string
		class            string
		userId           uuid.UUID
		wantFail         bool
		wantInvJSON      bool
		wantErrorMessage string
		invJSON          string
		ctx              context.Context
		cookie           string
		ok               bool
	}{
		{
			name:   "Test1",
			userId: uuid.New(),
			ctx:    context.Background(),
			class:  SeekerStr,
			cookie: "cookie",
			ok:     true,
		},
		{
			name:   "Test1",
			userId: uuid.New(),
			ctx:    context.Background(),
			class:  SeekerStr,
			cookie: "cookie",
			ok:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockAuthRepo.
					EXPECT().
					Delete(tt.cookie).
					Return(tt.ok)
			}

			cookie := session.Cookie{
				Cookie: tt.cookie,
			}

			gotState, err := h.DeleteSession(tt.ctx, &cookie)

			if !tt.wantFail {
				require.Equal(t, tt.ok, gotState.Ok)
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
