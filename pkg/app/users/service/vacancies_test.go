package users

import (
	mock_user_repo "2019_2_IBAT/pkg/app/users/service/mock_user_repo"
	. "2019_2_IBAT/pkg/pkg/interfaces"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateVacancy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserRepo := mock_user_repo.NewMockRepository(mockCtrl)
	h := UserService{
		Storage: mockUserRepo,
	}

	tests := []struct {
		name string
		// fields  fields
		// args    args
		record           AuthStorageValue
		wantFail         bool
		wantInvJSON      bool
		wantErrorMessage string
		vacancy          Vacancy
		invJSON          string
	}{
		{
			name: "Test1",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				CompanyName:  "PETUH",
				Experience:   "None",
				Position:     "",
				Tasks:        "drive",
				Requirements: "middle school education",
				WageFrom:     "50000 RUB",
				Conditions:   "nice team",
				About:        "nice job",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name: "Test2",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				CompanyName:  "PETUH",
				Experience:   "None",
				Position:     "",
				Tasks:        "drive",
				Requirements: "middle school education",
				WageFrom:     "50000 RUB",
				Conditions:   "nice team",
				About:        "nice job",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			wantFail:         true,
			wantErrorMessage: BadRequestMsg,
		},
		{
			name:             "Test3",
			wantFail:         true,
			wantErrorMessage: InvalidJSONMsg,
			wantInvJSON:      true,
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			invJSON: "{fsdfsd,cvvlxcfp}|}><P@#@:W:ED?SAD<FAS:DL |||",
		},
		{
			name:             "Test4",
			wantFail:         true,
			wantErrorMessage: ForbiddenMsg,
			wantInvJSON:      true,
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			invJSON: "{fsdfsd,cvvlxcfp}|}><P@#@:W:ED?SAD<FAS:DL |||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var str string
			if !tt.wantInvJSON {
				wantJSON, _ := json.Marshal(tt.vacancy)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					CreateVacancy(tt.vacancy).
					Return(true)
			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					CreateVacancy(tt.vacancy).
					Return(false)
			}
			_, err := h.CreateVacancy(r, tt.record)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_DeleteVacancy(t *testing.T) {
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
			vacancy: Vacancy{
				ID:           uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				CompanyName:  "PETUH",
				Experience:   "None",
				Position:     "",
				Tasks:        "drive",
				Requirements: "middle school education",
				WageFrom:     "50000 RUB",
				Conditions:   "nice team",
				About:        "nice job",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name:             "Test2",
			wantUnauth:       true,
			wantFail:         true,
			wantErrorMessage: ForbiddenMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetVacancy(tt.vacancyId).
					Return(tt.vacancy, nil)
				mockUserRepo.
					EXPECT().
					DeleteVacancy(tt.vacancyId).
					Return(nil)
			} else if !tt.wantUnauth {
				mockUserRepo.
					EXPECT().
					GetVacancy(tt.vacancyId).
					Return(tt.vacancy, nil)
				mockUserRepo.
					EXPECT().
					DeleteVacancy(tt.vacancyId).
					Return(errors.New(tt.wantErrorMessage))
			}

			err := h.DeleteVacancy(tt.vacancyId, tt.record)

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

func TestUserService_GetVacancy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserRepo := mock_user_repo.NewMockRepository(mockCtrl)
	h := UserService{
		Storage: mockUserRepo,
	}

	tests := []struct {
		name             string
		vacancy          Vacancy
		wantFail         bool
		wantErrorMessage string
	}{
		{
			name: "Test1",
			vacancy: Vacancy{
				ID:           uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				CompanyName:  "PETUH",
				Experience:   "None",
				Position:     "",
				Tasks:        "drive",
				Requirements: "middle school education",
				WageFrom:     "50000 RUB",
				Conditions:   "nice team",
				About:        "nice job",
			},
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantErrorMessage: InvalidIdMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetVacancy(tt.vacancy.ID).
					Return(tt.vacancy, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetVacancy(tt.vacancy.ID).
					Return(Vacancy{}, errors.New(tt.wantErrorMessage))
			}

			gotVacancy, err := h.GetVacancy(tt.vacancy.ID)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, tt.vacancy, gotVacancy, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_PutVacancy(t *testing.T) {
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
		wantInvJSON      bool
		wantErrorMessage string
		vacancy          Vacancy
		invJSON          string
	}{
		{
			name: "Test1",
			vacancy: Vacancy{
				ID:           uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				CompanyName:  "PETUH",
				Experience:   "None",
				Position:     "",
				Tasks:        "drive",
				Requirements: "middle school education",
				WageFrom:     "50000 RUB",
				Conditions:   "nice team",
				About:        "nice job",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name: "Test2",
			vacancy: Vacancy{
				ID:           uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				CompanyName:  "PETUH",
				Experience:   "None",
				Position:     "",
				Tasks:        "drive",
				Requirements: "middle school education",
				WageFrom:     "50000 RUB",
				Conditions:   "nice team",
				About:        "nice job",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			wantFail:         true,
			wantErrorMessage: BadRequestMsg,
		},
		{
			name:             "Test3",
			wantFail:         true,
			wantErrorMessage: InvalidJSONMsg,
			wantInvJSON:      true,
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			invJSON: "{fsdfsd,cvvlxcfp}|}><P@#@:W:ED?SAD<FAS:DL |||",
		},
		{
			name:             "Test4",
			wantFail:         true,
			wantErrorMessage: ForbiddenMsg,
			wantInvJSON:      true,
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			invJSON: "{fsdfsd,cvvlxcfp}|}><P@#@:W:ED?SAD<FAS:DL |||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var str string
			if !tt.wantInvJSON {
				wantJSON, _ := json.Marshal(tt.vacancy)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					PutVacancy(tt.vacancy, tt.record.ID, tt.vacancy.ID).
					Return(true)
			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					PutVacancy(tt.vacancy, tt.record.ID, tt.vacancy.ID).
					Return(false)
			}
			err := h.PutVacancy(tt.vacancy.ID, r, tt.record)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
