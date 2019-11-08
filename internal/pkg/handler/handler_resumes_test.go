package handler

import (
	mock_users "2019_2_IBAT/internal/pkg/handler/mock_users"
	"context"
	"errors"
	"strings"

	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/golang/mock/gomock"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetResume(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	wantResumes := []Resume{
		{
			ID:          uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"),
			OwnerID:     uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
			FirstName:   "Vova",
			SecondName:  "Zyablikov",
			Region:      "Moscow",
			PhoneNumber: "12345678910",
			BirthDate:   "1994-21-08",
			Sex:         "male",
			Citizenship: "Ukraine",
			Experience:  "7 years",
			Profession:  "programmer",
			Position:    "middle",
			Wage:        "100500",
			Education:   "MSU",
			About:       "Hello employer",
		},
	}

	mockUserService.
		EXPECT().
		GetResume(uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8")).
		Return(wantResumes[0], nil)

	mockUserService.
		EXPECT().
		GetResume(uuid.MustParse("12222222-9dad-11d1-80b1-00c04fd435c8")).
		Return(Resume{}, errors.New(InvalidIdMsg))

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
			path := fmt.Sprintf("/resume/%s", tc.pathArg)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/resume/{id}", h.GetResume)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotResume Resume
				json.Unmarshal(bytes, &gotResume)

				require.Equal(t, wantResumes[0], gotResume, "The two values should be the same.")
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

func TestHandler_GetResumes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	expected := []Resume{
		{
			ID:          uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd435c8"),
			OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
			FirstName:   "Vova",
			SecondName:  "Zyablikov",
			Region:      "Moscow",
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
		{
			ID:          uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"),
			OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
			FirstName:   "Vova",
			SecondName:  "Zyablikov",
			Region:      "Moscow",
			PhoneNumber: "12345678910",
			BirthDate:   "1994-21-08",
			Sex:         "male",
			Citizenship: "Ukraine",
			Experience:  "7 years",
			Profession:  "programmer",
			Position:    "middle",
			Wage:        "100500",
			Education:   "MSU",
			About:       "Hello employer",
		},
	}

	mockUserService.
		EXPECT().
		GetResumes().
		Return(expected, nil)

	r := httptest.NewRequest("GET", "/resumes/", nil)
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
			h.GetResumes(w, r)
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

func TestHandler_CreateResume(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	tests := []struct {
		name             string
		resume           Resume
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
			resume: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				Region:      "Moscow",
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
				Region:      "Moscow",
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
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
		},
		{
			name: "Test3",
			resume: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				Region:      "Moscow",
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
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name:             "Test4",
			resume:           Resume{},
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: BadRequestMsg,
			wantInvJSON:      true,
			invJSON:          "{testx: fdsfsdf, fdsfsdf'sdfsdf / fdsfsdf}",
		},
		{
			name: "Test5",
			resume: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				Region:      "Moscow",
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
			wantUnauth:       false,
			wantStatusCode:   http.StatusInternalServerError,
			wantErrorMessage: InternalErrorMsg,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var str string
			if !tc.wantInvJSON {
				wantJSON, _ := json.Marshal(tc.resume)
				str = string(wantJSON)
			} else {
				str = tc.invJSON
			}

			reader := strings.NewReader(str)

			req, _ := http.NewRequest("POST", "/resume", reader)

			id1 := uuid.New()

			if !tc.wantFail {
				mockUserService.
					EXPECT().
					CreateResume(req.Body, tc.record).
					Return(id1, nil)
				mockUserService.
					EXPECT().
					GetResume(id1).
					Return(tc.resume, nil)
			} else if !tc.wantUnauth {
				mockUserService.
					EXPECT().
					CreateResume(req.Body, tc.record).
					Return(uuid.UUID{}, errors.New(tc.wantErrorMessage))
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), tc.record)
			}

			router := mux.NewRouter()

			router.HandleFunc("/resume", h.CreateResume)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var id Id
				err := json.Unmarshal(bytes, &id)
				if err != nil {
					t.Errorf("corrupted returned id: %s", err)
				} else {
					gotResume, _ := h.UserService.GetResume(uuid.MustParse(id.Id))

					if rr.Code != http.StatusOK {
						t.Error("status is not ok")
					}

					require.Equal(t, tc.resume, gotResume, "The two values should be the same.")
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

func TestHandler_DeleteResume(t *testing.T) {
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
		resume           Resume
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
			resume: Resume{
				ID:          uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				FirstName:   "Vova",
				SecondName:  "Zyablikov",
				Region:      "Moscow",
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
				ID:      uuid.MustParse("7ba7b810-9bad-12d1-80b1-00c04fd430c8"),
				Role:    SeekerStr,
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
			resume: Resume{
				ID:          uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				FirstName:   "Vova",
				SecondName:  "Zyablikov",
				Region:      "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-10-08",
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
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				Role:    SeekerStr,
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
			resume: Resume{
				ID:          uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				FirstName:   "Vova",
				SecondName:  "Zyablikov",
				Region:      "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-10-08",
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
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/resume/%s", tc.pathArg)
			req, err := http.NewRequest("DELETE", path, nil)

			if err != nil {
				t.Fatal(err)
			}

			if !tc.wantFail {
				mockUserService.
					EXPECT().
					DeleteResume(tc.resume.ID, tc.record).
					Return(nil)
				mockUserService.
					EXPECT().
					GetResume(tc.resume.ID).
					Return(Resume{}, errors.New(InvalidIdMsg))
			} else if !tc.wantUnauth && tc.wantErrorMessage != InvalidIdMsg {
				mockUserService.
					EXPECT().
					DeleteResume(tc.resume.ID, tc.record).
					Return(errors.New(tc.wantErrorMessage))
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), tc.record)
			}

			router := mux.NewRouter()

			router.HandleFunc("/resume/{id}", h.DeleteResume)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				gotResume, err := h.UserService.GetResume(uuid.MustParse(tc.pathArg))

				var empResume Resume
				if err != nil && gotResume != empResume {
					require.Equal(t, gotResume, empResume, "The two values should be the same.")
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

func TestHandler_PutResume(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: mockUserService,
	}

	tests := []struct {
		name             string
		pathArg          string
		resume           Resume
		record           AuthStorageValue
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:    "Test1",
			pathArg: "7ba7b810-9dad-12d1-80b1-00c04fd430c8",
			resume: Resume{
				ID:          uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				Region:      "Kaliningrad",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Sweeden",
				Experience:  "15 years",
				Profession:  "programmer",
				Position:    "senior",
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
			name:    "Test2",
			pathArg: "7ba7b810-9dad-12d1-80b1-00c04fd430c8",
			resume: Resume{
				ID:          uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				Region:      "Kaliningrad",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Sweeden",
				Experience:  "15 years",
				Profession:  "programmer",
				Position:    "senior",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c1"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
			},
		},
		{
			name:             "Test3",
			wantUnauth:       true,
			pathArg:          "7ba7b810-9dad-12d1-80b1-00c04fd430c8",
			resume:           Resume{},
			wantFail:         true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantJSON, _ := json.Marshal(tc.resume)

			reader := strings.NewReader(string(wantJSON))

			path := fmt.Sprintf("/resume/%s", tc.pathArg)
			req, err := http.NewRequest("PUT", path, reader)
			if err != nil {
				t.Fatal(err)
			}

			if !tc.wantFail {
				mockUserService.
					EXPECT().
					PutResume(tc.resume.ID, req.Body, tc.record).
					Return(nil)
				mockUserService.
					EXPECT().
					GetResume(tc.resume.ID).
					Return(tc.resume, nil)
			} else if !tc.wantUnauth {
				mockUserService.
					EXPECT().
					PutResume(tc.resume.ID, req.Body, tc.record).
					Return(errors.New(tc.wantErrorMessage))
			}

			req.Header.Set("Content-Type", "application/json")

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), tc.record)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/resume/{id}", h.PutResume)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				gotResume, _ := h.UserService.GetResume(uuid.MustParse(tc.pathArg))

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}
				if tc.resume != gotResume {
					require.Equal(t, tc.resume, gotResume, "The two values should be the same.")
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
