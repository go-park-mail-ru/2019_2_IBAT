package users

import (
	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	"fmt"

	mock_notifs "2019_2_IBAT/pkg/app/users/service/mock_notifs"
	"2019_2_IBAT/pkg/app/users/service/mock_recommends"
	mock_user_repo "2019_2_IBAT/pkg/app/users/service/mock_user_repo"

	. "2019_2_IBAT/pkg/pkg/models"
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

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()
	mockNotifRepo := mock_notifs.NewMockServiceClient(mockCtrl2)

	h := UserService{
		Storage:      mockUserRepo,
		NotifService: mockNotifRepo,
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
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				CompanyName:  "PETUH",
				Experience:   "None",
				Position:     "",
				Tasks:        "drive",
				Requirements: "middle school education",
				WageFrom:     "50000 RUB",
				Conditions:   "nice team",
				About:        "nice job",
				Spheres:      []Pair{},
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
				Spheres:      []Pair{},
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
				wantJSON, _ := tt.vacancy.MarshalJSON()
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))
			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					CreateVacancy(gomock.Any()).
					Return(true)
				mockUserRepo.
					EXPECT().
					GetTagIDs(tt.vacancy.Spheres).
					Return([]uuid.UUID{}, nil)
				mockNotifRepo.
					EXPECT().
					SendNotification(gomock.Any(), gomock.Any()).
					Return(&notifsproto.Bool{Ok: true}, nil)

			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					CreateVacancy(gomock.Any()).
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
					GetVacancy(tt.vacancyId, tt.record.ID).
					Return(tt.vacancy, nil)
				mockUserRepo.
					EXPECT().
					DeleteVacancy(tt.vacancyId).
					Return(nil)
			} else if !tt.wantUnauth {
				mockUserRepo.
					EXPECT().
					GetVacancy(tt.vacancyId, tt.record.ID).
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

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()
	mockRecomRepo := mock_recommends.NewMockServiceClient(mockCtrl2)

	h := UserService{
		Storage:      mockUserRepo,
		RecomService: mockRecomRepo,
	}

	tests := []struct {
		name             string
		vacancy          Vacancy
		wantFail         bool
		wantErrorMessage string
		record           AuthStorageValue
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
				Spheres:      []Pair{},
			},
			record: AuthStorageValue{
				ID: uuid.New(),
			},
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantErrorMessage: InvalidIdMsg,
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
				Spheres:      []Pair{},
			},
			record: AuthStorageValue{
				ID: uuid.New(),
			},
		},
	}

	// mockRecomRepo.EXPECT().SetTagIDs()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetVacancyTagIDs(tt.vacancy.ID).
					Return([]uuid.UUID{}, nil)

				mockRecomRepo.
					EXPECT().
					SetTagIDs(gomock.Any(), gomock.Any()).
					Return(&recomsproto.Bool{Ok: true}, nil)

				mockUserRepo.
					EXPECT().
					GetVacancy(tt.vacancy.ID, tt.record.ID).
					Return(tt.vacancy, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetVacancyTagIDs(tt.vacancy.ID).
					Return([]uuid.UUID{}, nil)
				mockRecomRepo.
					EXPECT().
					SetTagIDs(gomock.Any(), gomock.Any()).
					Return(&recomsproto.Bool{Ok: true}, nil)
				mockUserRepo.
					EXPECT().
					GetVacancy(tt.vacancy.ID, tt.record.ID).
					Return(Vacancy{}, errors.New(tt.wantErrorMessage))
			}

			gotVacancy, err := h.GetVacancy(tt.vacancy.ID, tt.record)

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

func TestUserService_GetVacancies(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserRepo := mock_user_repo.NewMockRepository(mockCtrl)

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()
	mockRecomRepo := mock_recommends.NewMockServiceClient(mockCtrl2)

	h := UserService{
		Storage:      mockUserRepo,
		RecomService: mockRecomRepo,
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
		seekers          []Seeker
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
			wantErrorMessage: InternalErrorMsg,
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantRecomms:      true,
			wantErrorMessage: InternalErrorMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetTagIDs(gomock.Any()).
					Return([]uuid.UUID{}, nil)
				mockUserRepo.
					EXPECT().
					GetVacancies(tt.record, gomock.Any()).
					Return(expVacancies, nil)
				mockRecomRepo.
					EXPECT().
					SetTagIDs(gomock.Any(), gomock.Any()).
					Return(&recomsproto.Bool{Ok: true}, nil)
			} else if tt.wantRecomms {
				mockUserRepo.
					EXPECT().
					GetTagIDs(gomock.Any()).
					Return([]uuid.UUID{}, nil)
				mockUserRepo.
					EXPECT().
					GetVacancies(tt.record, gomock.Any()).
					Return(expVacancies, nil)
				mockRecomRepo.
					EXPECT().
					SetTagIDs(gomock.Any(), gomock.Any()).
					Return(&recomsproto.Bool{Ok: false}, fmt.Errorf(InternalErrorMsg))
			} else {
				mockUserRepo.
					EXPECT().
					GetTagIDs(gomock.Any()).
					Return([]uuid.UUID{}, nil)
				mockUserRepo.
					EXPECT().
					GetVacancies(tt.record, gomock.Any()).
					Return([]Vacancy{}, fmt.Errorf(InternalErrorMsg))
			}
			paramsDummyMap := map[string]interface{}{}
			tagDummyMap := map[string]interface{}{}

			gotSeeks, err := h.GetVacancies(tt.record, paramsDummyMap, tagDummyMap)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, expVacancies, gotSeeks, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
