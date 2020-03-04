package handler

import (
	"context"
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

	"2019_2_IBAT/pkg/app/auth/session"
	mock_auth "2019_2_IBAT/pkg/app/server/handler/mock_auth"
	mock_users "2019_2_IBAT/pkg/app/server/handler/mock_users"
	. "2019_2_IBAT/pkg/pkg/models"
)

func TestHandler_CreateSeeker(t *testing.T) {
	mockCtrl1 := gomock.NewController(t)
	defer mockCtrl1.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl1)

	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()

	mockAuthService := mock_auth.NewMockServiceClient(mockCtrl2)

	h := Handler{
		UserService: mockUserService,
		AuthService: mockAuthService,
	}

	tests := []struct {
		name              string
		seekReg           SeekerReg
		wantRole          string
		wantFail          bool
		wantStatusCode    int
		wantErrorMessage  string
		wantInvJSON       bool
		wantCreateSession bool
		invJSON           string
		ctx               context.Context
		sessionMsg        session.Session
	}{
		{
			name: "Test1",
			seekReg: SeekerReg{
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
			},
			wantRole: SeekerStr,
			ctx:      context.Background(),
			sessionMsg: session.Session{
				Id:    uuid.New().String(),
				Class: SeekerStr,
			},
		},
		{
			name: "Test2",
			seekReg: SeekerReg{
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
			},
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: EmailExistsMsg,
		},
		{
			name: "Test3",
			seekReg: SeekerReg{
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
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
			seekReg: SeekerReg{
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
			},
			wantRole:          SeekerStr,
			wantFail:          true,
			wantStatusCode:    http.StatusInternalServerError,
			wantErrorMessage:  InternalErrorMsg,
			wantCreateSession: true,
			ctx:               context.Background(),
			sessionMsg: session.Session{
				Id:    uuid.New().String(),
				Class: SeekerStr,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var str string
			if !tc.wantInvJSON {
				wantJSON, _ := json.Marshal(tc.seekReg)
				str = string(wantJSON)
			} else {
				str = tc.invJSON
			}

			reader := strings.NewReader(string(str))

			req, err := http.NewRequest("POST", "/seeker", reader)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			if !tc.wantFail {
				mockUserService.
					EXPECT().
					CreateSeeker(req.Body).
					Return(uuid.MustParse(tc.sessionMsg.Id), nil)
				mockAuthService.
					EXPECT().
					CreateSession(gomock.Any(), &tc.sessionMsg).
					Return(
						&session.CreateSessionInfo{
							ID:      tc.sessionMsg.Id,
							Role:    tc.sessionMsg.Class,
							Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
							Cookie:  "cookie",
						}, nil)
			} else if tc.wantCreateSession {
				mockUserService.
					EXPECT().
					CreateSeeker(req.Body).
					Return(uuid.MustParse(tc.sessionMsg.Id), nil)
				mockAuthService.
					EXPECT().
					CreateSession(gomock.Any(), &tc.sessionMsg).
					Return(
						&session.CreateSessionInfo{}, errors.New("Create session error"))
			} else {
				mockUserService.
					EXPECT().
					CreateSeeker(req.Body).
					Return(uuid.UUID{}, errors.New(tc.wantErrorMessage))
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/seeker", h.CreateSeeker)
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

func TestHandler_GetSeekerById(t *testing.T) {
	mockCtrl1 := gomock.NewController(t)
	defer mockCtrl1.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl1)

	h := Handler{
		UserService: mockUserService,
	}

	tests := []struct {
		name             string
		pathArg          string
		seeker           Seeker
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:    "Test1",
			pathArg: "6ba7b811-9dad-11d1-80b1-00c04fd430c8",
			seeker: Seeker{
				ID:         uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
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
		path := fmt.Sprintf("/seeker/%s", tc.pathArg)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		if !tc.wantFail {
			mockUserService.
				EXPECT().
				GetSeeker(tc.seeker.ID).
				Return(tc.seeker, nil)
		} else {
			mockUserService.
				EXPECT().
				GetSeeker(uuid.MustParse(tc.pathArg)).
				Return(Seeker{}, errors.New(InvalidIdMsg))
		}

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/seeker/{id}", h.GetSeekerById)
		router.ServeHTTP(rr, req)

		if !tc.wantFail {
			if rr.Code != http.StatusOK {
				t.Error("status is not ok")
			}
			bytes, _ := ioutil.ReadAll(rr.Body)

			var gotSeeker Seeker
			json.Unmarshal(bytes, &gotSeeker)

			require.Equal(t, tc.seeker, gotSeeker, "The two values should be the same.")
		} else {
			bytes, _ := ioutil.ReadAll(rr.Body)
			var gotError Error
			json.Unmarshal(bytes, &gotError)

			require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
			require.Equal(t, tc.wantErrorMessage, gotError.Message, "The two values should be the same.")
		}
	}
}
