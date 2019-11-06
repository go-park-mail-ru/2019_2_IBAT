package users

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	mock_user_repo "2019_2_IBAT/internal/pkg/users/service/mock_user_repo"
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

func TestUserService_CreateResume(t *testing.T) {
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
		resume           Resume
		invJSON          string
	}{
		{
			name: "Test1",
			resume: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name: "Test2",
			resume: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
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
				Role:    SeekerStr,
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
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			invJSON: "{fsdfsd,cvvlxcfp}|}><P@#@:W:ED?SAD<FAS:DL |||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var str string
			if !tt.wantInvJSON {
				wantJSON, _ := json.Marshal(tt.resume)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					CreateResume(tt.resume).
					Return(true)
			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					CreateResume(tt.resume).
					Return(false)
			}
			_, err := h.CreateResume(r, tt.record)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_DeleteResume(t *testing.T) {
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
		resumeId         uuid.UUID
		record           AuthStorageValue
		resume           Resume
		wantFail         bool
		wantUnauth       bool
		wantErrorMessage string
	}{
		{
			name:     "Test1",
			resumeId: uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
			resume: Resume{
				ID:          uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
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
					GetResume(tt.resumeId).
					Return(tt.resume, nil)
				mockUserRepo.
					EXPECT().
					DeleteResume(tt.resumeId).
					Return(nil)
			} else if !tt.wantUnauth {
				mockUserRepo.
					EXPECT().
					GetResume(tt.resumeId).
					Return(tt.resume, nil)
				mockUserRepo.
					EXPECT().
					DeleteResume(tt.resumeId).
					Return(errors.New(tt.wantErrorMessage))
			}

			err := h.DeleteResume(tt.resumeId, tt.record)

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

func TestUserService_GetResume(t *testing.T) {
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
		resumeId         uuid.UUID
		resume           Resume
		wantFail         bool
		wantErrorMessage string
	}{
		{
			name:     "Test1",
			resumeId: uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
			resume: Resume{
				ID:          uuid.MustParse("1ba7b811-9dad-11d1-0000-00004fd430c8"),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
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
					GetResume(tt.resumeId).
					Return(tt.resume, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetResume(tt.resumeId).
					Return(Resume{}, errors.New(tt.wantErrorMessage))
			}

			gotResume, err := h.GetResume(tt.resumeId)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, tt.resume, gotResume, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_PutResume(t *testing.T) {
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
		resume           Resume
		invJSON          string
	}{
		{
			name: "Test1",
			resume: Resume{
				ID:          uuid.New(),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name: "Test2",
			resume: Resume{
				ID:          uuid.New(),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			wantFail:         true,
			wantErrorMessage: InternalErrorMsg,
		},
		{
			name:             "Test3",
			wantFail:         true,
			wantErrorMessage: InvalidJSONMsg,
			wantInvJSON:      true,
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
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
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			invJSON: "{fsdfsd,cvvlxcfp}|}><P@#@:W:ED?SAD<FAS:DL |||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var str string
			if !tt.wantInvJSON {
				wantJSON, _ := json.Marshal(tt.resume)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					PutResume(tt.resume, tt.record.ID, tt.resume.ID).
					Return(true)
			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					PutResume(tt.resume, tt.record.ID, tt.resume.ID).
					Return(false)
			}
			err := h.PutResume(tt.resume.ID, r, tt.record)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
