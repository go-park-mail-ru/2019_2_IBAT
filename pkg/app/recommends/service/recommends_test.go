package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"2019_2_IBAT/pkg/app/recommends/service/mock_recommends"

	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	. "2019_2_IBAT/pkg/pkg/models"
)

func TestUserService_SetTagIDs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthRepo := mock_recommends.NewMockRepository(mockCtrl)

	h := Service{
		Storage: mockAuthRepo,
	}

	Ids := []string{
		uuid.New().String(),
		uuid.New().String(),
	}
	tests := []struct {
		name             string
		role             string
		userId           uuid.UUID
		expires          string
		wantFail         bool
		wantErrorMessage string
		ctx              context.Context
	}{
		{
			name:   "Test1",
			userId: uuid.New(),
			ctx:    context.Background(),
			role:   SeekerStr,
		},
		{
			name:             "Test2",
			userId:           uuid.New(),
			ctx:              context.Background(),
			role:             EmployerStr,
			wantFail:         true,
			wantErrorMessage: InternalErrorMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockAuthRepo.
					EXPECT().
					SetTagIDs(AuthStorageValue{
						ID:      tt.userId,
						Role:    tt.role,
						Expires: tt.expires,
					},
						Ids).
					Return(nil)
			} else {
				mockAuthRepo.
					EXPECT().
					SetTagIDs(AuthStorageValue{
						ID:      tt.userId,
						Role:    tt.role,
						Expires: tt.expires,
					},
						Ids).
					Return(fmt.Errorf(tt.wantErrorMessage))
			}

			msg := recomsproto.SetTagIDsMessage{
				ID:      tt.userId.String(),
				Role:    tt.role,
				Expires: tt.expires,
				IDs:     Ids,
			}

			_, err := h.SetTagIDs(tt.ctx, &msg)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_GetTagIDs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthRepo := mock_recommends.NewMockRepository(mockCtrl)

	h := Service{
		Storage: mockAuthRepo,
	}

	Ids := []string{
		uuid.New().String(),
		uuid.New().String(),
	}
	tests := []struct {
		name             string
		role             string
		userId           uuid.UUID
		expires          string
		wantFail         bool
		wantErrorMessage string
		ctx              context.Context
	}{
		{
			name:   "Test1",
			userId: uuid.New(),
			ctx:    context.Background(),
			role:   SeekerStr,
		},
		{
			name:             "Test2",
			userId:           uuid.New(),
			ctx:              context.Background(),
			role:             EmployerStr,
			wantFail:         true,
			wantErrorMessage: BadRequestMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockAuthRepo.
					EXPECT().
					GetTagIDs(AuthStorageValue{
						ID:      tt.userId,
						Role:    tt.role,
						Expires: tt.expires,
					}).
					Return(Ids, nil)
			} else {
				mockAuthRepo.
					EXPECT().
					GetTagIDs(AuthStorageValue{
						ID:      tt.userId,
						Role:    tt.role,
						Expires: tt.expires,
					}).
					Return([]string{}, fmt.Errorf(BadRequestMsg))
			}

			msg := recomsproto.GetTagIDsMessage{
				ID:      tt.userId.String(),
				Role:    tt.role,
				Expires: tt.expires,
			}

			_, err := h.GetTagIDs(tt.ctx, &msg)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_GetUsersForTags(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAuthRepo := mock_recommends.NewMockRepository(mockCtrl)

	h := Service{
		Storage: mockAuthRepo,
	}

	userIds := []string{
		uuid.New().String(),
		uuid.New().String(),
	}

	tagIds := []string{
		uuid.New().String(),
		uuid.New().String(),
	}
	tests := []struct {
		name             string
		wantFail         bool
		wantErrorMessage string
		invJSON          string
		ctx              context.Context
	}{
		{
			name: "Test1",
			ctx:  context.Background(),
		},
		{
			name:             "Test2",
			ctx:              context.Background(),
			wantFail:         true,
			wantErrorMessage: InternalErrorMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockAuthRepo.
					EXPECT().
					GetUsersForTags(userIds).
					Return(tagIds, nil)
			} else {
				mockAuthRepo.
					EXPECT().
					GetUsersForTags(userIds).
					Return([]string{}, fmt.Errorf(InternalErrorMsg))
			}

			msg := recomsproto.IDsMessage{
				IDs: userIds,
			}

			_, err := h.GetUsersForTags(tt.ctx, &msg)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
