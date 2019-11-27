package handler

import (
	"context"
	"encoding/json"
	"errors"
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

	"2019_2_IBAT/pkg/app/auth"
	"2019_2_IBAT/pkg/app/auth/session"
	mock_auth "2019_2_IBAT/pkg/app/server/handler/mock_auth"
	mock_users "2019_2_IBAT/pkg/app/server/handler/mock_users"
	. "2019_2_IBAT/pkg/pkg/models"
)

func TestHandler_CreateSession(t *testing.T) {

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
		name string
		// authInput         UserAuthInput
		// wantID            uuid.UUID
		authInput         UserAuthInput
		invJSON           string
		wantFail          bool
		wantRole          string
		wantStatusCode    int
		wantErrorMessage  string
		wantInvJSON       bool
		wantCreateSession bool
		ctx               context.Context
		sessionMsg        session.Session
	}{
		{
			name: "Test1",
			authInput: UserAuthInput{
				Email:    "petushki@mail.com",
				Password: "1234",
			},
			// wantID:   uuid.New(),
			wantRole: EmployerStr,
			wantFail: false,
			ctx:      context.Background(),
			sessionMsg: session.Session{
				Id:    uuid.New().String(),
				Class: EmployerStr,
			},
		},
		{
			name: "Test2",
			authInput: UserAuthInput{
				Email:    "some_another@mail.com",
				Password: "12345",
			},
			wantFail: false,
			wantRole: SeekerStr,
			ctx:      context.Background(),
			sessionMsg: session.Session{
				Id:    uuid.New().String(),
				Class: SeekerStr,
			},
		},
		{
			name: "Test3",
			authInput: UserAuthInput{
				Email:    "some_another@mail.com",
				Password: "1234567",
			},
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvPassOrEmailMsg,
		},
		{
			name: "Test4",
			authInput: UserAuthInput{
				Email:    "some_another@mail.com",
				Password: "1234567",
			},
			wantFail:          true,
			wantRole:          SeekerStr,
			wantStatusCode:    http.StatusInternalServerError,
			wantErrorMessage:  InternalErrorMsg,
			wantCreateSession: true,
			ctx:               context.Background(),
			sessionMsg: session.Session{
				Id:    uuid.New().String(),
				Class: SeekerStr,
			},
		},
		{
			name:             "Test5",
			authInput:        UserAuthInput{},
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidJSONMsg,
			wantInvJSON:      true,
			invJSON:          "{'lagin': sdfdfsdf pasword: sdfsdf }",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			var str string
			if !tc.wantInvJSON {
				wantJSON, _ := json.Marshal(tc.authInput)
				str = string(wantJSON)
			} else {
				str = tc.invJSON
			}

			reader := strings.NewReader(str)

			req, err := http.NewRequest("POST", "/auth", reader)

			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			if !tc.wantFail {
				mockUserService.
					EXPECT().
					CheckUser(tc.authInput.Email, tc.authInput.Password).
					Return(uuid.MustParse(tc.sessionMsg.Id), tc.sessionMsg.Class, true)
				mockAuthService.
					EXPECT().
					CreateSession(tc.ctx, &tc.sessionMsg).
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
					CheckUser(tc.authInput.Email, tc.authInput.Password).
					Return(uuid.MustParse(tc.sessionMsg.Id), tc.sessionMsg.Class, true)
				mockAuthService.
					EXPECT().
					CreateSession(tc.ctx, &tc.sessionMsg).
					Return(
						&session.CreateSessionInfo{}, errors.New("Create session error"))
			} else if !tc.wantInvJSON {
				mockUserService.
					EXPECT().
					CheckUser(tc.authInput.Email, tc.authInput.Password).
					Return(uuid.UUID{}, "", false)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/auth", h.CreateSession)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotRole Role
				json.Unmarshal(bytes, &gotRole)

				require.Equal(t, tc.wantRole, gotRole.Role, "The two values should be the same.")

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
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

func TestHandler_GetSession(t *testing.T) {
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
		name             string
		record           AuthStorageValue
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:             "Test1",
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name: "Test2",
			record: AuthStorageValue{
				Role:    SeekerStr,
				ID:      uuid.New(),
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader("") ///why
			req, err := http.NewRequest("GET", "/auth", reader)

			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), tc.record)
			}

			router := mux.NewRouter()
			router.HandleFunc("/auth", h.GetSession)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotRole Role
				json.Unmarshal(bytes, &gotRole)

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				if tc.record.Role != gotRole.Role {
					require.Equal(t, tc.record.Role, gotRole.Role, "The two values should be the same.")
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

func TestHandler_DeleteSession(t *testing.T) {
	mockCtrl2 := gomock.NewController(t)
	defer mockCtrl2.Finish()

	mockAuthService := mock_auth.NewMockServiceClient(mockCtrl2)

	h := Handler{
		AuthService: mockAuthService,
	}

	tests := []struct {
		name             string
		record           AuthStorageValue
		wantFail         bool
		wantUnauth       bool
		wantFailDelete   bool
		wantStatusCode   int
		wantErrorMessage string
		ctx              context.Context
		cookie           session.Cookie
	}{
		{
			name:       "Test1",
			wantFail:   false,
			wantUnauth: false,
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			ctx: context.Background(),
			cookie: session.Cookie{
				Cookie: "cookie",
			},
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
			ctx:              context.Background(),

			cookie: session.Cookie{
				Cookie: "aaabbbaaaaa",
			},
		},
		{
			name:           "Test3",
			wantFail:       true,
			wantUnauth:     false,
			wantFailDelete: true,
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: BadRequestMsg,
			ctx:              context.Background(),
			cookie: session.Cookie{
				Cookie: "cookie",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/auth", nil)

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			if !tc.wantFail {
				mockAuthService.
					EXPECT().
					DeleteSession(tc.ctx, &tc.cookie).
					Return(&session.Bool{Ok: true}, nil)
			} else if tc.wantFailDelete {
				mockAuthService.
					EXPECT().
					DeleteSession(tc.ctx, &tc.cookie).
					Return(&session.Bool{Ok: false}, nil)
			}
			if !tc.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tc.cookie.Cookie,
				}

				req.AddCookie(&cookie)
			}
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/auth", h.DeleteSession)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
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
