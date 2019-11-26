package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	mock_auth "2019_2_IBAT/pkg/app/server/handler/mock_auth"
	mock_users "2019_2_IBAT/pkg/app/server/handler/mock_users"
	. "2019_2_IBAT/pkg/pkg/interfaces"
)

func TestHandler_CreateEmployer(t *testing.T) {
	mockCtrl1 := gomock.NewController(t)
	defer mockCtrl1.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl1)

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()

	mockAuthService := mock_auth.NewMockService(mockCtrl2)

	h := Handler{
		UserService: mockUserService,
		AuthService: mockAuthService,
	}

	tests := []struct {
		name              string
		emplReg           EmployerReg
		wantRole          string
		wantFail          bool
		wantStatusCode    int
		wantErrorMessage  string
		wantInvJSON       bool
		wantCreateSession bool
		invJSON           string
	}{
		{
			name: "Test1",
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
			wantRole: EmployerStr,
		},
		{
			name: "Test2",
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
			wantRole:         EmployerStr,
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: EmailExistsMsg,
		},
		{
			name: "Test3",
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
			wantRole:         SeekerStr,
			wantFail:         true,
			wantInvJSON:      true,
			invJSON:          "{sfsdf: some email, login: password: sdfdf}",
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidJSONMsg,
		},
		{
			name: "Test4",
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
			wantRole:          SeekerStr,
			wantFail:          true,
			wantStatusCode:    http.StatuspkgServerError,
			wantErrorMessage:  pkgErrorMsg,
			wantCreateSession: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var str string
			if !tc.wantInvJSON {
				wantJSON, _ := json.Marshal(tc.emplReg)
				str = string(wantJSON)
			} else {
				str = tc.invJSON
			}

			reader := strings.NewReader(string(str))

			req, err := http.NewRequest("POST", "/employer", reader)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			id1 := uuid.New()
			if !tc.wantFail {
				mockUserService.
					EXPECT().
					CreateEmployer(req.Body).
					Return(id1, nil)
				mockAuthService.
					EXPECT().
					CreateSession(id1, EmployerStr).
					Return(
						AuthStorageValue{
							ID:      id1,
							Role:    EmployerStr,
							Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
						},
						"cookie", nil)
			} else if tc.wantCreateSession {
				mockUserService.
					EXPECT().
					CreateEmployer(req.Body).
					Return(id1, nil)
				mockAuthService.
					EXPECT().
					CreateSession(id1, EmployerStr).
					Return(
						AuthStorageValue{},
						"", errors.New("Create session error"))
			} else {
				mockUserService.
					EXPECT().
					CreateEmployer(req.Body).
					Return(uuid.UUID{}, errors.New(tc.wantErrorMessage))
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/employer", h.CreateEmployer)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var Role Role
				json.Unmarshal(bytes, &Role)

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				if tc.wantRole != Role.Role {
					require.Equal(t, tc.wantRole, Role.Role, "The two values should be the same.")
				}
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Message, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_GetEmployerById(t *testing.T) {
	mockCtrl1 := gomock.NewController(t)
	defer mockCtrl1.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl1)

	h := Handler{
		UserService: mockUserService,
	}

	tests := []struct {
		name             string
		pathArg          string
		employer         Employer
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:    "Test1",
			pathArg: "6ba7b811-9dad-11d1-80b1-00c04fd430c8",
			employer: Employer{
				ID:               uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
				CompanyName:      "MCDonalds",
				Site:             "petushki.com",
				Email:            "petushki@mail.com",
				FirstName:        "Vova",
				SecondName:       "Zyablikov",
				Password:         "",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				Region:           "Petushki",
				EmplNum:          "322",
			},
			wantFail: false,
		},
		{
			name:             "Test2",
			pathArg:          "6ba7b810-9bad-11d1-80b1-00c04fd430c8",
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidIdMsg,
		},
	}

	for _, tc := range tests {
		path := fmt.Sprintf("/employer/%s", tc.pathArg)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !tc.wantFail {
			mockUserService.
				EXPECT().
				GetEmployer(tc.employer.ID).
				Return(tc.employer, nil)
		} else {
			mockUserService.
				EXPECT().
				GetEmployer(uuid.MustParse(tc.pathArg)).
				Return(Employer{}, errors.New(InvalidIdMsg))
		}

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/employer/{id}", h.GetEmployerById)
		router.ServeHTTP(rr, req)

		if !tc.wantFail {
			if rr.Code != http.StatusOK {
				t.Error("status is not ok")
			}
			bytes, _ := ioutil.ReadAll(rr.Body)

			var gotEmployer Employer
			json.Unmarshal(bytes, &gotEmployer)

			require.Equal(t, tc.employer, gotEmployer, "The two values should be the same.")
		} else {
			bytes, _ := ioutil.ReadAll(rr.Body)
			var gotError Error
			json.Unmarshal(bytes, &gotError)

			require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
			require.Equal(t, tc.wantErrorMessage, gotError.Message, "The two values should be the same.")
		}
	}
}

func TestHandler_GetEmployers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	expected := []Employer{
		{
			ID:               uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
			CompanyName:      "MCDonalds",
			Site:             "petushki.com",
			Email:            "petushki@mail.com",
			FirstName:        "Vova",
			SecondName:       "Zyablikov",
			Password:         "",
			PhoneNumber:      "12345678911",
			ExtraPhoneNumber: "12345678910",
			Region:           "Petushki",
			EmplNum:          "322",
		},
		{
			ID:               uuid.MustParse("1ba7b811-9dad-11d1-80b1-00c04fd430c8"),
			CompanyName:      "IDs",
			Site:             "IDS.com",
			Email:            "ids@mail.com",
			FirstName:        "Kostya",
			SecondName:       "Zyablikov",
			Password:         "1234",
			PhoneNumber:      "12345678911",
			ExtraPhoneNumber: "12345678910",
			Region:           "Moscow",
			EmplNum:          "322",
		},
	}

	expectedJSON, _ := json.Marshal(expected)

	tests := []struct {
		name             string
		expected         string
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:     "Test1",
			expected: string(expectedJSON),
			wantFail: false,
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantStatusCode:   http.StatuspkgServerError,
			wantErrorMessage: pkgErrorMsg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantFail {
				mockUserService.
					EXPECT().
					GetEmployers().
					Return(expected, nil)
			} else {
				mockUserService.
					EXPECT().
					GetEmployers().
					Return([]Employer{}, errors.New(pkgErrorMsg))
			}

			r := httptest.NewRequest("GET", "/employers/", nil)
			w := httptest.NewRecorder()
			h.GetEmployers(w, r)

			if !tt.wantFail {
				if w.Code != http.StatusOK {
					t.Error("status is not ok")
				}
				bytes, _ := ioutil.ReadAll(w.Body)

				if string(bytes) != tt.expected {
					require.Equal(t, tt.expected, string(bytes), "The two values should be the same.")
				}
			} else {
				bytes, _ := ioutil.ReadAll(w.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tt.wantStatusCode, w.Code, "The two values should be the same.")
				require.Equal(t, tt.wantErrorMessage, gotError.Message, "The two values should be the same.")
			}
		})
	}
}
