package handler

import (
	mock_auth "2019_2_IBAT/pkg/app/server/handler/mock_auth"
	mock_users "2019_2_IBAT/pkg/app/server/handler/mock_users"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"strings"

	. "2019_2_IBAT/pkg/pkg/models"
	"encoding/json"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

//cookie needs to be set for the test to run
// func TestHandler_DeleteUser(t *testing.T) {
// 	mockCtrl1 := gomock.NewController(t)
// 	defer mockCtrl1.Finish()

// 	mockUserService := mock_users.NewMockService(mockCtrl1)

// 	mockCtrl2 := gomock.NewController(t)
// 	defer mockCtrl2.Finish()

// 	mockAuthService := mock_auth.NewMockServiceClient(mockCtrl2)

// 	h := Handler{
// 		UserService: mockUserService,
// 		AuthService: mockAuthService,
// 	}

// 	tests := []struct {
// 		name              string
// 		wantRole          string
// 		wantFail          bool
// 		wantStatusCode    int
// 		wantErrorMessage  string
// 		wantInvJSON       bool
// 		wantCreateSession bool
// 		wantUnauth        bool
// 		invJSON           string
// 		userId            uuid.UUID
// 		class             string
// 		ctx               context.Context
// 	}{
// 		{
// 			name:     "Test1",
// 			wantRole: EmployerStr,
// 			userId:   uuid.UUID{},
// 			ctx:      context.Background(),
// 		},
// 		// {
// 		// 	name: "Test2",
// 		// 	wantRole:         EmployerStr,
// 		// 	wantFail:         true,
// 		// 	wantStatusCode:   http.StatusBadRequest,
// 		// 	wantErrorMessage: EmailExistsMsg,
// 		// },
// 		// {
// 		// 	name: "Test3",
// 		// 	wantRole:         SeekerStr,
// 		// 	wantFail:         true,
// 		// 	wantInvJSON:      true,
// 		// 	invJSON:          "{sfsdf: some email, login: password: sdfdf}",
// 		// 	wantStatusCode:   http.StatusBadRequest,
// 		// 	wantErrorMessage: InvalidJSONMsg,
// 		// },
// 		// {
// 		// 	name: "Test4",
// 		// 	wantRole:          EmployerStr,
// 		// 	wantFail:          true,
// 		// 	wantStatusCode:    http.StatusInternalServerError,
// 		// 	wantErrorMessage:  InternalErrorMsg,
// 		// 	wantCreateSession: true,
// 		// 	ctx:               context.Background(),
// 		// 	sessionMsg: session.Session{
// 		// 		Id:    uuid.New().String(),
// 		// 		Class: SeekerStr,
// 		// 	},
// 		// },
// 	}
// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var str string
// 			// if !tc.wantInvJSON {
// 			// 	wantJSON, _ := json.Marshal(tc.emplReg)
// 			// 	str = string(wantJSON)
// 			// } else {
// 			// 	str = tc.invJSON
// 			// }

// 			reader := strings.NewReader(string(str))

// 			req, err := http.NewRequest("POST", "/"+tc.class, reader)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			record := AuthStorageValue{
// 				ID:      tc.userId,
// 				Role:    tc.class,
// 				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
// 			}

// 			req.Header.Set("session-id", "cookie")
// 			// http.SetCookie(req, cookie)

// 			// id1 := uuid.New()
// 			if !tc.wantFail {
// 				mockUserService.
// 					EXPECT().
// 					DeleteUser(record).
// 					Return(nil)

// 				mockAuthService.
// 					EXPECT().
// 					DeleteSession(gomock.Any(), &session.Cookie{
// 						Cookie: "cookie",
// 					}).
// 					Return(&session.Bool{Ok: true}, nil)
// 			}
// 			//  else if tc.wantCreateSession {
// 			// 	mockUserService.
// 			// 	DeleteUser(AuthStorageValue{
// 			// 		ID: sessionMsg.ID,
// 			// 		Role: sessionMsg.Role,
// 			// 	}).
// 			// 		Return(uuid.MustParse(tc.sessionMsg.Id), nil)
// 			// 	mockAuthService.
// 			// 		EXPECT().
// 			// 		CreateSession(tc.ctx, &session.Session{
// 			// 			Id:    tc.sessionMsg.Id,
// 			// 			Class: tc.wantRole,
// 			// 		}).
// 			// 		Return(&session.CreateSessionInfo{}, errors.New("Create session error"))

// 			// } else {
// 			// 	mockUserService.
// 			// 		EXPECT().
// 			// 		CreateEmployer(req.Body).
// 			// 		Return(uuid.UUID{}, errors.New(tc.wantErrorMessage))
// 			// }

// 			rr := httptest.NewRecorder()

// 			ctx := context.TODO()
// 			if !tc.wantUnauth {
// 				ctx = NewContext(req.Context(), record)
// 			}

// 			router := mux.NewRouter()
// 			router.HandleFunc("/"+tc.class, h.DeleteUser)
// 			router.ServeHTTP(rr, req.WithContext(ctx))

// 			if !tc.wantFail {
// 				if rr.Code != http.StatusOK {
// 					t.Error("status is not ok")
// 				}
// 			} else {
// 				var gotError Error
// 				var bytes []byte
// 				json.Unmarshal(bytes, &gotError)

// 				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
// 				require.Equal(t, tc.wantErrorMessage, gotError.Message, "The two values should be the same.")
// 			}
// 		})
// 	}
// }

func TestHandler_PutUser(t *testing.T) {
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

	employer := EmployerReg{
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
	}

	seeker := SeekerReg{
		Email:      "third@mail.com",
		FirstName:  "Petr",
		SecondName: "Zyablikov",
		Password:   "12345",
	}

	tests := []struct {
		name             string
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
		wantForbidden    bool
		wantUnauth       bool
		userId           uuid.UUID
		role             string
		ctx              context.Context
	}{
		{
			name:   "Test1",
			role:   EmployerStr,
			userId: uuid.UUID{},
			ctx:    context.Background(),
		},
		{
			name:     "Test2",
			role:     SeekerStr,
			userId:   uuid.UUID{},
			wantFail: false,
			ctx:      context.Background(),
		},
		{
			name:             "Test3",
			wantUnauth:       true,
			role:             SeekerStr,
			wantFail:         true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},

		{
			name:             "Test4",
			role:             EmployerStr,
			userId:           uuid.UUID{},
			wantFail:         true,
			wantForbidden:    true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
			ctx:              context.Background(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var str string
			var wantJSON []byte
			if tc.role == SeekerStr {
				wantJSON, _ = json.Marshal(seeker)
			} else {
				wantJSON, _ = json.Marshal(employer)
			}
			str = string(wantJSON)

			reader := strings.NewReader(string(str))

			req, err := http.NewRequest("POST", "/"+tc.role, reader)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			record := AuthStorageValue{
				ID:      tc.userId,
				Role:    tc.role,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			}

			// id1 := uuid.New()
			if !tc.wantFail {
				if tc.role == SeekerStr {
					mockUserService.
						EXPECT().
						PutSeeker(gomock.Any(), tc.userId).
						Return(nil)
				} else if tc.role == EmployerStr {
					mockUserService.
						EXPECT().
						PutEmployer(gomock.Any(), tc.userId).
						Return(nil)
				}
			} else {
				if !tc.wantUnauth && tc.wantForbidden {
					if tc.role == SeekerStr {
						mockUserService.
							EXPECT().
							PutSeeker(gomock.Any(), tc.userId).
							Return(fmt.Errorf(ForbiddenMsg))
					} else {
						mockUserService.
							EXPECT().
							PutEmployer(gomock.Any(), tc.userId).
							Return(fmt.Errorf(ForbiddenMsg))
					}
				}
			}
			rr := httptest.NewRecorder()

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), record)
			}

			router := mux.NewRouter()
			router.HandleFunc("/"+tc.role, h.PutUser)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				_ = gotError.UnmarshalJSON(bytes)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Message, "The two values should be the same.")
			}
		})
	}
}
