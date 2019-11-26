package users

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"
	mock_user_repo "2019_2_IBAT/pkg/app/users/service/mock_user_repo"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateSeeker(t *testing.T) {
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
		seekerReg        Seeker
		invJSON          string
	}{
		{
			name: "Test1",
			seekerReg: Seeker{
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
			},
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantErrorMessage: EmailExistsMsg,
		},
		{
			name:             "Test3",
			wantFail:         true,
			wantErrorMessage: InvalidJSONMsg,
			wantInvJSON:      true,
			invJSON:          "{fsdfsd,cvvlxcfp}|}><P@#@:W:ED?SAD<FAS:DL |||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var str string
			if !tt.wantInvJSON {
				wantJSON, _ := json.Marshal(tt.seekerReg)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					CreateSeeker(tt.seekerReg).
					Return(true)
			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					CreateSeeker(tt.seekerReg).
					Return(false)
			}
			_, err := h.CreateSeeker(r)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_PutSeeker(t *testing.T) {
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
		seekReg          SeekerReg
		invJSON          string
	}{
		{
			name: "Test1",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: EmployerStr,
			},
			seekReg: SeekerReg{
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
			},
		},
		{
			name: "Test2",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			seekReg: SeekerReg{
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
			},
			wantFail:         true,
			wantErrorMessage: BadRequestMsg,
		},
		{
			name:             "Test3",
			wantFail:         true,
			wantErrorMessage: InvalidJSONMsg,
			wantInvJSON:      true,
			invJSON:          "{fsdfsd,cvvlxcfp}|}><P@#@:W:ED?SAD<FAS:DL |||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var str string
			if !tt.wantInvJSON {
				wantJSON, _ := json.Marshal(tt.seekReg)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					PutSeeker(tt.seekReg, tt.record.ID).
					Return(true)
			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					PutSeeker(tt.seekReg, tt.record.ID).
					Return(false)
			}
			err := h.PutSeeker(r, tt.record.ID)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_GetSeeker(t *testing.T) {
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
		seeker           Seeker
	}{
		{
			name: "Test1",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			seeker: Seeker{
				ID:         uuid.New(),
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
			},
		},
		{
			name: "Test2",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			wantFail:         true,
			wantErrorMessage: pkgErrorMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetSeeker(tt.record.ID).
					Return(tt.seeker, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetSeeker(tt.record.ID).
					Return(Seeker{}, errors.New(tt.wantErrorMessage))
			}

			gotSeek, err := h.GetSeeker(tt.record.ID)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, tt.seeker, gotSeek, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_GetSeekers(t *testing.T) {
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
		seekers          []Seeker
	}{
		{
			name: "Test1",
			seekers: []Seeker{
				{
					ID:         uuid.New(),
					Email:      "third@mail.com",
					FirstName:  "Petr",
					SecondName: "Zyablikov",
					Password:   "12345",
				},
				{
					ID:         uuid.New(),
					Email:      "seconds@mail.com",
					FirstName:  "Petr",
					SecondName: "Zyablikov",
					Password:   "12345",
				},
			},
		},
		{
			name:             "Test2",
			wantFail:         false,
			wantErrorMessage: pkgErrorMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetSeekers().
					Return(tt.seekers, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetSeekers().
					Return([]Seeker{}, errors.New(tt.wantErrorMessage))
			}

			gotSeeks, err := h.GetSeekers()

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, tt.seekers, gotSeeks, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
