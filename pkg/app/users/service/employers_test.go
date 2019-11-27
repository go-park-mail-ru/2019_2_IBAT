package users

import (
	. "2019_2_IBAT/pkg/pkg/models"
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

func TestUserService_CreateEmployer(t *testing.T) {
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
		emplReg          Employer
		invJSON          string
	}{
		{
			name: "Test1",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			emplReg: Employer{
				CompanyName:      "MCDonalds",
				Site:             "petushki.com",
				Email:            "petushki@mail.com",
				FirstName:        "Vova",
				SecondName:       "Zyablikov",
				Password:         "1234",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				Region:           "Petushki",
				EmplNum:          "322",
			},
		},
		{
			name: "Test2",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
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
				wantJSON, _ := json.Marshal(tt.emplReg)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					CreateEmployer(tt.emplReg).
					Return(true)
			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					CreateEmployer(tt.emplReg).
					Return(false)
			}
			_, err := h.CreateEmployer(r)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_PutEmployer(t *testing.T) {
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
		emplReg          EmployerReg
		invJSON          string
	}{
		{
			name: "Test1",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: EmployerStr,
			},
			emplReg: EmployerReg{
				CompanyName:      "MCDonalds",
				Site:             "petushki.com",
				Email:            "petushki@mail.com",
				FirstName:        "Vova",
				SecondName:       "Zyablikov",
				Password:         "1234",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				Region:           "Petushki",
				EmplNum:          "322",
			},
		},
		{
			name: "Test2",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			emplReg: EmployerReg{
				CompanyName:      "MCDonalds",
				Site:             "petushki.com",
				Email:            "petushki@mail.com",
				FirstName:        "Vova",
				SecondName:       "Zyablikov",
				Password:         "1234",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				Region:           "Petushki",
				EmplNum:          "322",
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
				wantJSON, _ := json.Marshal(tt.emplReg)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			r := ioutil.NopCloser(strings.NewReader(string(str)))

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					PutEmployer(tt.emplReg, tt.record.ID).
					Return(true)
			} else if !tt.wantInvJSON {
				mockUserRepo.
					EXPECT().
					PutEmployer(tt.emplReg, tt.record.ID).
					Return(false)
			}
			err := h.PutEmployer(r, tt.record.ID)

			if !tt.wantFail {
				require.Equal(t, err, nil)
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_GetEmployer(t *testing.T) {
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
		employer         Employer
	}{
		{
			name: "Test1",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: EmployerStr,
			},
			employer: Employer{
				ID:               uuid.New(),
				CompanyName:      "MCDonalds",
				Site:             "petushki.com",
				Email:            "petushki@mail.com",
				FirstName:        "Vova",
				SecondName:       "Zyablikov",
				Password:         "1234",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				Region:           "Petushki",
				EmplNum:          "322",
			},
		},
		{
			name: "Test2",
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: EmployerStr,
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
					GetEmployer(tt.record.ID).
					Return(tt.employer, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetEmployer(tt.record.ID).
					Return(Employer{}, errors.New(tt.wantErrorMessage))
			}

			gotEmpl, err := h.GetEmployer(tt.record.ID)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, tt.employer, gotEmpl, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}

func TestUserService_GetEmployers(t *testing.T) {
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
		employers        []Employer
	}{
		{
			name: "Test1",
			employers: []Employer{
				{
					ID:               uuid.New(),
					CompanyName:      "MCDonalds",
					Site:             "petushki.com",
					Email:            "petushki@mail.com",
					FirstName:        "Vova",
					SecondName:       "Zyablikov",
					Password:         "1234",
					PhoneNumber:      "12345678911",
					ExtraPhoneNumber: "12345678910",
					Region:           "Petushki",
					EmplNum:          "322",
				},
				{
					ID:               uuid.New(),
					CompanyName:      "corp",
					Site:             "corp.com",
					Email:            "some@mail.com",
					FirstName:        "Vova",
					SecondName:       "Zyablikov",
					Password:         "1234",
					PhoneNumber:      "12345678911",
					ExtraPhoneNumber: "12345678910",
					Region:           "Petushki",
					EmplNum:          "322",
				},
			},
		},
		{
			name:             "Test2",
			wantFail:         false,
			wantErrorMessage: pkgErrorMsg,
		},
	}

	var dummy_map map[string]interface{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetEmployers(dummy_map).
					Return(tt.employers, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetEmployers(dummy_map).
					Return([]Employer{}, errors.New(tt.wantErrorMessage))
			}

			gotEmpls, err := h.GetEmployers(dummy_map)

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, tt.employers, gotEmpls, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
