package users

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"
	mock_user_repo "2019_2_IBAT/pkg/pkg/users/service/mock_user_repo"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestUserService_DeleteUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserRepo := mock_user_repo.NewMockRepository(mockCtrl)
	h := UserService{
		Storage: mockUserRepo,
	}

	tests := []struct {
		name             string
		record           AuthStorageValue
		wantFail         bool
		wantErrorMessage string
	}{
		{
			name: "Test1",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
		},
		{
			name: "Test2",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			wantFail:         true,
			wantErrorMessage: "DeleteUser: error while deleting",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					DeleteUser(tt.record.ID).
					Return(nil)
			} else {
				mockUserRepo.
					EXPECT().
					DeleteUser(tt.record.ID).
					Return(errors.New(tt.wantErrorMessage))
			}

			err := h.DeleteUser(tt.record)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
