package users

import (
	mock_user_repo "2019_2_IBAT/pkg/app/users/service/mock_user_repo"
	"fmt"
	"testing"

	. "2019_2_IBAT/pkg/pkg/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetResponds(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserRepo := mock_user_repo.NewMockRepository(mockCtrl)

	h := UserService{
		Storage: mockUserRepo,
	}

	expResponds := []Respond{
		{
			Status:    AwaitSt,
			ResumeID:  uuid.New(),
			VacancyID: uuid.New(),
		},
		{
			Status:    RejectedSt,
			ResumeID:  uuid.New(),
			VacancyID: uuid.New(),
		},
	}

	tests := []struct {
		name             string
		record           AuthStorageValue
		wantFail         bool
		wantBothArgsErr  bool
		wantErrorMessage string
		paramsMap        map[string]string
	}{
		{
			name: "Test1",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			paramsMap: make(map[string]string),
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantErrorMessage: BadRequestMsg,
			paramsMap:        make(map[string]string),
		},
		{
			name:             "Test3",
			wantFail:         true,
			wantErrorMessage: BadRequestMsg,
			paramsMap: map[string]string{
				"vacancy_id": uuid.New().String(),
				"resume_id":  uuid.New().String(),
			},
			wantBothArgsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetResponds(tt.record, gomock.Any()).
					Return(expResponds, nil)
			} else if !tt.wantBothArgsErr {
				mockUserRepo.
					EXPECT().
					GetResponds(tt.record, gomock.Any()).
					Return([]Respond{}, fmt.Errorf(tt.wantErrorMessage))
			}

			gotResponds, err := h.GetResponds(tt.record, tt.paramsMap)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, expResponds, gotResponds, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
