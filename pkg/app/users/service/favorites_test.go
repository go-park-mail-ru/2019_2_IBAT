package users

import (
	mock_user_repo "2019_2_IBAT/pkg/app/users/service/mock_user_repo"
	"fmt"
	"testing"
	"time"

	. "2019_2_IBAT/pkg/pkg/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateFavorite(t *testing.T) {
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
		vacancyId        uuid.UUID
	}{
		{
			name: "Test1",
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			vacancyId: uuid.New(),
		},
		{
			name: "Test2",
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			vacancyId:        uuid.New(),
			wantFail:         true,
			wantErrorMessage: "Error while creating favorite_vacancy",
		},
		{
			name:             "Test3",
			wantFail:         true,
			wantErrorMessage: "Invalid action",
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					CreateFavorite(gomock.Any()).
					Return(true)
			} else if tt.wantErrorMessage != "Invalid action" {
				mockUserRepo.
					EXPECT().
					CreateFavorite(gomock.Any()).
					Return(false)
			}
			err := h.CreateFavorite(tt.vacancyId, tt.record)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_GetFavoriteVacancies(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserRepo := mock_user_repo.NewMockRepository(mockCtrl)

	h := UserService{
		Storage: mockUserRepo,
	}

	expVacancies := []Vacancy{
		{
			ID:           uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd430c8"),
			OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
			CompanyName:  "MCDonalds",
			Experience:   "None",
			Position:     "",
			Tasks:        "bring food to costumers",
			Requirements: "middle school education",
			WageFrom:     "1000 USD",
			Conditions:   "nice team",
			About:        "nice job",
		},

		{
			ID:           uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"),
			OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
			CompanyName:  "PETUH",
			Experience:   "None",
			Position:     "driver",
			Tasks:        "drive",
			Requirements: "middle school education",
			WageFrom:     "50000 RUB",
			Conditions:   "nice team",
			About:        "nice job",
		},
	}

	tests := []struct {
		name             string
		record           AuthStorageValue
		wantFail         bool
		wantRecomms      bool
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
			name:             "Test2",
			wantFail:         true,
			wantErrorMessage: "Invalid action",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetFavoriteVacancies(gomock.Any()).
					Return(expVacancies, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetFavoriteVacancies(tt.record).
					Return([]Vacancy{}, fmt.Errorf("Invalid action"))
			}

			gotVacs, err := h.GetFavoriteVacancies(tt.record)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, expVacancies, gotVacs, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

// func (h *UserService) DeleteFavoriteVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) error {

// 	err := h.Storage.DeleteFavoriteVacancy(vacancyId, authInfo)

// 	if err != nil {
// 		return errors.New(InternalErrorMsg)
// 	}

// 	return nil
// }

func TestUserService_DeleteFavoriteVacancy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserRepo := mock_user_repo.NewMockRepository(mockCtrl)
	h := UserService{
		Storage: mockUserRepo,
	}

	tests := []struct {
		name             string
		vacancyId        uuid.UUID
		record           AuthStorageValue
		vacancy          Vacancy
		wantFail         bool
		wantUnauth       bool
		wantErrorMessage string
	}{
		{
			name:      "Test1",
			vacancyId: uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantErrorMessage: InternalErrorMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					DeleteFavoriteVacancy(tt.vacancyId, tt.record).
					Return(nil)
			} else {
				mockUserRepo.
					EXPECT().
					DeleteFavoriteVacancy(tt.vacancyId, tt.record).
					Return(errors.New(InternalErrorMsg))
			}

			err := h.DeleteFavoriteVacancy(tt.vacancyId, tt.record)

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
