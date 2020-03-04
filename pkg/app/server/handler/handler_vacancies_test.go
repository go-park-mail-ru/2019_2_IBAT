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

	mock_users "2019_2_IBAT/pkg/app/server/handler/mock_users"
	. "2019_2_IBAT/pkg/pkg/models"
)

func TestHandler_GetVacancies(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	expected := []Vacancy{
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

	params := make(map[string]interface{})

	mockUserService.
		EXPECT().
		GetVacancies(gomock.Any(), params, gomock.Any()).
		Return(expected, nil)

	r := httptest.NewRequest("GET", "/vacancies/", nil)
	w := httptest.NewRecorder()

	expectedJSON, _ := json.Marshal(expected)

	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Test1",
			expected: string(expectedJSON),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.GetVacancies(w, r)
			if w.Code != http.StatusOK {
				t.Error("status is not ok")
			}
			bytes, _ := ioutil.ReadAll(w.Body)

			if string(bytes) != tt.expected {
				require.Equal(t, tt.expected, string(bytes), "The two values should be the same.")
			}
		})
	}
}

func TestHandler_GetVacancy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	wantVacancies := []Vacancy{
		{
			ID:           uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"),
			OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
			CompanyName:  "PETUH",
			Experience:   "None",
			Position:     "",
			Tasks:        "drive",
			Requirements: "middle school education",
			WageFrom:     "50000 RUB",
			Conditions:   "nice team",
			About:        "nice job",
		},
	}

	mockUserService.
		EXPECT().
		GetVacancy(uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"), gomock.Any()).
		Return(wantVacancies[0], nil)

	mockUserService.
		EXPECT().
		GetVacancy(uuid.MustParse("12222222-9dad-11d1-80b1-00c04fd435c8"), gomock.Any()).
		Return(Vacancy{}, errors.New(InvalidIdMsg))

	tests := []struct {
		name             string
		pathArg          string
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:     "Test1",
			pathArg:  "22222222-9dad-11d1-80b1-00c04fd435c8",
			wantFail: false,
		},
		{
			name:             "Test2",
			pathArg:          "12222222-9dad-11d1-80b1-00c04fd435c8",
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidIdMsg,
		},
		{
			name:             "Test3",
			pathArg:          "фвапвапвпа_а<#ыва/||s",
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidIdMsg,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/vacancy/%s", tc.pathArg)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/vacancy/{id}", h.GetVacancy)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotVacancy Vacancy
				json.Unmarshal(bytes, &gotVacancy)

				require.Equal(t, wantVacancies[0], gotVacancy, "The two values should be the same.")
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

func TestHandler_CreateVacancy(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	tests := []struct {
		name             string
		vacancy          Vacancy
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
		wantInvJSON      bool
		invJSON          string
		record           AuthStorageValue
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
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
		},
		{
			name: "Test3",
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
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name:             "Test4",
			vacancy:          Vacancy{},
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: BadRequestMsg,
			wantInvJSON:      true,
			invJSON:          "{testx: fdsfsdf, fdsfsdf'sdfsdf / fdsfsdf}",
		},
		{
			name: "Test5",
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
			wantFail:   true,
			wantUnauth: false,
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			wantStatusCode:   http.StatusInternalServerError,
			wantErrorMessage: InternalErrorMsg,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var str string
			if !tc.wantInvJSON {
				wantJSON, _ := json.Marshal(tc.vacancy)
				str = string(wantJSON)
			} else {
				str = tc.invJSON
			}

			reader := strings.NewReader(str)

			req, _ := http.NewRequest("POST", "/vacancy", reader)

			id1 := uuid.New()

			if !tc.wantFail {
				mockUserService.
					EXPECT().
					CreateVacancy(req.Body, tc.record).
					Return(id1, nil)
				mockUserService.
					EXPECT().
					GetVacancy(id1, gomock.Any()).
					Return(tc.vacancy, nil)
			} else if !tc.wantUnauth {
				mockUserService.
					EXPECT().
					CreateVacancy(req.Body, tc.record).
					Return(uuid.UUID{}, errors.New(tc.wantErrorMessage))
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), tc.record)
			}

			router := mux.NewRouter()

			router.HandleFunc("/vacancy", h.CreateVacancy)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var id Id
				err := json.Unmarshal(bytes, &id)
				if err != nil {
					t.Errorf("corrupted returned id: %s", err)
				} else {
					gotVacancy, _ := h.UserService.GetVacancy(uuid.MustParse(id.Id), AuthStorageValue{})

					if rr.Code != http.StatusOK {
						t.Error("status is not ok")
					}

					require.Equal(t, tc.vacancy, gotVacancy, "The two values should be the same.")
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

func TestHandler_DeleteVacancy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	tests := []struct {
		name             string
		pathArg          string
		cookieValue      string
		wantUnauth       bool
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
		vacancy          Vacancy
		record           AuthStorageValue
	}{
		{
			name:             "Test1",
			pathArg:          "7aa7b810-9dad-11d1-72b5-04c04fd430c8",
			cookieValue:      "aaaaaaaaaa",
			wantUnauth:       false,
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
			vacancy: Vacancy{
				ID:           uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"),
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
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-11004fd430c8"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name:             "Test2",
			pathArg:          "11111111-9dad-11d1-80b1-00c04fd430c8",
			wantUnauth:       true,
			wantFail:         true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name:       "Test3",
			pathArg:    "7aa7b810-9dad-11d1-72b5-04c04fd430c8",
			wantFail:   false,
			wantUnauth: false,
			vacancy: Vacancy{
				ID:           uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"),
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
			name:             "Test4",
			pathArg:          "11111111-9dad-11d1-80b1-00#|<>c04fd430c8",
			wantUnauth:       false,
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidIdMsg,
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name:             "Test5",
			pathArg:          "7aa7b810-9dad-11d1-72b5-04c04fd430c8",
			wantUnauth:       false,
			wantFail:         true,
			wantStatusCode:   http.StatusInternalServerError,
			wantErrorMessage: InternalErrorMsg,
			vacancy: Vacancy{
				ID:           uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"),
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
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/vacancy/%s", tc.pathArg)
			req, err := http.NewRequest("DELETE", path, nil)

			if err != nil {
				t.Fatal(err)
			}

			if !tc.wantFail {
				mockUserService.
					EXPECT().
					DeleteVacancy(tc.vacancy.ID, tc.record).
					Return(nil)
				mockUserService.
					EXPECT().
					GetVacancy(tc.vacancy.ID, gomock.Any()).
					Return(Vacancy{}, errors.New(InvalidIdMsg))
			} else if !tc.wantUnauth && tc.wantErrorMessage != InvalidIdMsg {
				mockUserService.
					EXPECT().
					DeleteVacancy(tc.vacancy.ID, tc.record).
					Return(errors.New(tc.wantErrorMessage))
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), tc.record)
			}

			router := mux.NewRouter()

			router.HandleFunc("/vacancy/{id}", h.DeleteVacancy)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				gotVacancy, err := h.UserService.GetVacancy(uuid.MustParse(tc.pathArg), AuthStorageValue{})

				var empVacancy Vacancy
				if err != nil {
					require.Equal(t, gotVacancy, empVacancy, "The two values should be the same.")
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

func TestHandler_PutVacancy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	tests := []struct {
		name             string
		pathArg          string
		vacancy          Vacancy
		record           AuthStorageValue
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:    "Test1",
			pathArg: "7ba7b810-9dad-12d1-80b1-00c04fd430c8",
			vacancy: Vacancy{
				ID:           uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
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
			name:    "Test2",
			pathArg: "7ba7b810-9dad-12d1-80b1-00c04fd430c8",
			vacancy: Vacancy{
				ID:           uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
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
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c1"),
				Role:    EmployerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
		},
		{
			name:             "Test3",
			wantUnauth:       true,
			pathArg:          "7ba7b810-9dad-12d1-80b1-00c04fd430c8",
			vacancy:          Vacancy{},
			wantFail:         true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantJSON, _ := json.Marshal(tc.vacancy)

			reader := strings.NewReader(string(wantJSON))

			path := fmt.Sprintf("/vacancy/%s", tc.pathArg)
			req, err := http.NewRequest("PUT", path, reader)
			if err != nil {
				t.Fatal(err)
			}

			if !tc.wantFail {
				mockUserService.
					EXPECT().
					PutVacancy(tc.vacancy.ID, req.Body, tc.record).
					Return(nil)
				mockUserService.
					EXPECT().
					GetVacancy(tc.vacancy.ID, gomock.Any()).
					Return(tc.vacancy, nil)
			} else if !tc.wantUnauth {
				mockUserService.
					EXPECT().
					PutVacancy(tc.vacancy.ID, req.Body, tc.record).
					Return(errors.New(tc.wantErrorMessage))
			}

			req.Header.Set("Content-Type", "application/json")

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), tc.record)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/vacancy/{id}", h.PutVacancy)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				gotVacancy, _ := h.UserService.GetVacancy(uuid.MustParse(tc.pathArg), tc.record)

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				require.Equal(t, tc.vacancy, gotVacancy, "The two values should be the same.")
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
