package handler

import (
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"2019_2_IBAT/internal/pkg/users"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func newHandler() *Handler {
	return &Handler{
		AuthService: auth.AuthService{
			Storage: auth.MapAuthStorage{
				Storage: make(map[string]AuthStorageValue),
				Mu:      &sync.Mutex{},
			},
		},
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:           &sync.Mutex{},
				EmplMu:          &sync.Mutex{},
				ResMu:           &sync.Mutex{},
				SeekerStorage:   map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{},
				ResumeStorage:   map[uuid.UUID]Resume{},
			},
		},
	}
}

func TestHandler_GetResumes(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: make(map[string]AuthStorageValue),
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						Email:      "some@mail.com",
						FirstName:  "Vova",
						SecondName: "Zyablikov",
						Password:   "1234",
						Resumes: []uuid.UUID{
							uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd435c8"),
						},
					},
					uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"): {
						Email:      "third@mail.com",
						FirstName:  "Petr",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{},
				ResumeStorage: map[uuid.UUID]Resume{
					uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd435c8"): {
						OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						FirstName:   "Vova",
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
					uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"): {
						OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						FirstName:   "Vova",
						SecondName:  "Zyablikov",
						City:        "Moscow",
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
				},
			},
		},
	}
	r := httptest.NewRequest("GET", "/resumes/", nil)
	w := httptest.NewRecorder()

	expectedJSON, _ := json.Marshal(h.UserService.Storage.GetResumes())

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

func TestHandler_GetEmployers(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: make(map[string]AuthStorageValue),
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:         &sync.Mutex{},
				EmplMu:        &sync.Mutex{},
				ResMu:         &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies:        make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "BMSTU",
						Site:             "bmstu.ru",
						Email:            "bmstu@mail.com",
						FirstName:        "Tolya",
						SecondName:       "Alex",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Moscow",
						EmplNum:          "1830",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{},
			},
		},
	}
	r := httptest.NewRequest("GET", "/employers/", nil)
	w := httptest.NewRecorder()

	expectedEmployers := h.UserService.Storage.GetEmployers()

	for i, item := range expectedEmployers {
		item.Password = ""
		expectedEmployers[i] = item
	}

	expectedJSON, _ := json.Marshal(expectedEmployers)

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
			h.GetEmployers(w, r)
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

func TestHandler_GetEmployerById(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: make(map[string]AuthStorageValue),
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:         &sync.Mutex{},
				EmplMu:        &sync.Mutex{},
				ResMu:         &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies:        make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "BMSTU",
						Site:             "bmstu.ru",
						Email:            "bmstu@mail.com",
						FirstName:        "Tolya",
						SecondName:       "Alex",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Moscow",
						EmplNum:          "1830",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{},
			},
		},
	}

	expectedEmployers := h.UserService.Storage.GetEmployers()

	for i, item := range expectedEmployers {
		item.Password = ""
		expectedEmployers[i] = item
	}

	employerJSON1, _ := json.Marshal(expectedEmployers[uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8")])
	employerJSON2, _ := json.Marshal(expectedEmployers[uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8")])

	tests := []struct {
		name     string
		pathArg  string
		expected string
	}{
		{
			name:     "Test1",
			pathArg:  "6ba7b810-9dad-11d1-80b1-00c04fd430c8",
			expected: string(employerJSON1),
		},
		{
			name:     "Test2",
			pathArg:  "6ba7b811-9dab-11d1-80b1-00c04fd430c8",
			expected: string(employerJSON2),
		},
	}

	for _, tc := range tests {
		path := fmt.Sprintf("/employers/%s", tc.pathArg)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/employers/{id}", h.GetEmployerById)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error("status is not ok")
		}
		bytes, _ := ioutil.ReadAll(rr.Body)

		if string(bytes) != tc.expected {
			require.Equal(t, tc.expected, string(bytes), "The two values should be the same.")
		}
	}
}

func TestHandler_GetSeekerById(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: make(map[string]AuthStorageValue),
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"): {
						Email:      "some@mail.com",
						FirstName:  "Vova",
						SecondName: "Zyablikov",
						Password:   "1234",
						Resumes: []uuid.UUID{
							uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
						},
					},

					uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"): {
						Email:      "third@mail.com",
						FirstName:  "Petr",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes: []uuid.UUID{
							uuid.MustParse("7ba7b810-9dad-11d1-71b5-04c04fd430c8"),
						},
					},
					uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{},
				ResumeStorage: map[uuid.UUID]Resume{
					uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"): {
						OwnerID:     uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
						FirstName:   "Vova",
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
					uuid.MustParse("7ba7b810-9dad-11d1-71b5-04c04fd430c8"): {
						OwnerID:     uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
						FirstName:   "Petr",
						SecondName:  "Zyablikov",
						City:        "Moscow",
						PhoneNumber: "12345678910",
						BirthDate:   "1994-21-08",
						Sex:         "male",
						Citizenship: "Russia",
						Experience:  "8 years",
						Profession:  "programmer",
						Position:    "senior",
						Wage:        "100500",
						Education:   "MSU",
						About:       "Hello employer",
					},
				},
			},
		},
	}

	expectedSeekers := h.UserService.Storage.GetSeekers()

	for i, item := range expectedSeekers {
		item.Password = ""
		expectedSeekers[i] = item
	}

	seekJSON1, _ := json.Marshal(expectedSeekers[uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8")])
	seekJSON2, _ := json.Marshal(expectedSeekers[uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8")])

	tests := []struct {
		name     string
		pathArg  string
		expected string
	}{
		{
			name:     "Test1",
			pathArg:  "6ba7b811-9dad-11d1-80b1-00c04fd430c8",
			expected: string(seekJSON1),
		},
		{
			name:     "Test2",
			pathArg:  "6ba7b810-9bad-11d1-80b1-00c04fd430c8",
			expected: string(seekJSON2),
		},
	}

	for _, tc := range tests {
		path := fmt.Sprintf("/seekers/%s", tc.pathArg)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/seekers/{id}", h.GetSeekerById)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error("status is not ok")
		}
		bytes, _ := ioutil.ReadAll(rr.Body)

		if string(bytes) != tc.expected {
			require.Equal(t, tc.expected, string(bytes), "The two values should be the same.")
		}
	}
}

func TestHandler_GetEmployer(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"aaaaaaaaab": {
					ID:      uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:         &sync.Mutex{},
				EmplMu:        &sync.Mutex{},
				ResMu:         &sync.Mutex{},
				VacMu:         &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies:        make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "BMSTU",
						Site:             "bmstu.ru",
						Email:            "bmstu@mail.com",
						FirstName:        "Tolya",
						SecondName:       "Alex",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Moscow",
						EmplNum:          "1830",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{},
			},
		},
	}

	expectedEmployers := h.UserService.Storage.GetEmployers()

	emplJSON1, _ := json.Marshal(expectedEmployers[uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8")])
	emplJSON2, _ := json.Marshal(expectedEmployers[uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8")])

	tests := []struct {
		name        string
		cookieValue string
		expected    string
	}{
		{
			name:        "Test1",
			cookieValue: "aaaaaaaaaa",
			expected:    string(emplJSON1),
		},
		{
			name:        "Test1",
			cookieValue: "aaaaaaaaab",
			expected:    string(emplJSON2),
		},
	}

	for _, tc := range tests {
		req, err := http.NewRequest("GET", "/employer", nil)
		if err != nil {
			t.Fatal(err)
		}

		cookie := http.Cookie{
			Name:  auth.CookieName,
			Value: tc.cookieValue,
		}

		req.AddCookie(&cookie)

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/employer", h.GetEmployer)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error("status is not ok")
		}
		bytes, _ := ioutil.ReadAll(rr.Body)

		if string(bytes) != tc.expected {
			require.Equal(t, tc.expected, string(bytes), "The two values should be the same.")
		}
	}
}

func TestHandler_CreateVacancy(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"bbbbbbbbbb": {
					ID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"cccccccc": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				VacMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies:        make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "newpetushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage:  map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{},
			},
		},
	}

	tests := []struct {
		name             string
		cookieValue      string
		vacancy          Vacancy
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:        "Test1",
			cookieValue: "aaaaaaaaaa",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				CompanyName:  "MCDonalds",
				Experience:   "None",
				Profession:   "waiter",
				Position:     "",
				Tasks:        "bring food to costumers",
				Requirements: "middle school education",
				Wage:         "1000 USD",
				Conditions:   "nice team",
				About:        "nice job",
			},
		},
		{
			name:        "Test2",
			cookieValue: "bbbbbbbbbb",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
				CompanyName:  "PETUH",
				Experience:   "None",
				Profession:   "driver",
				Position:     "",
				Tasks:        "drive",
				Requirements: "middle school education",
				Wage:         "50000 RUB",
				Conditions:   "nice team",
				About:        "nice job",
			},
		},
		{
			name:        "Test3",
			cookieValue: "aaaa",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				CompanyName:  "MCDonalds",
				Experience:   "None",
				Profession:   "waiter",
				Position:     "",
				Tasks:        "bring food to costumers",
				Requirements: "middle school education",
				Wage:         "1000 USD",
				Conditions:   "nice team",
				About:        "nice job",
			},
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name:        "InvTest4",
			cookieValue: "cccccccc",
			vacancy: Vacancy{
				CompanyName:  "MCDonalds",
				Experience:   "None",
				Profession:   "waiter",
				Position:     "",
				Tasks:        "bring food to costumers",
				Requirements: "middle school education",
				Wage:         "1000 USD",
				Conditions:   "nice team",
				About:        "nice job",
			},
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantJSON, _ := json.Marshal(tc.vacancy)

			reader := strings.NewReader(string(wantJSON))

			req, err := http.NewRequest("POST", "/vacancy", reader)

			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			if !tc.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tc.cookieValue,
				}

				req.AddCookie(&cookie)
			}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/vacancy", h.CreateVacancy)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var id Id
				json.Unmarshal(bytes, &id)
				gotVacancy, _ := h.UserService.Storage.GetVacancy(uuid.MustParse(id.Id))

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				if tc.vacancy != gotVacancy {
					require.Equal(t, tc.vacancy, string(bytes), "The two values should be the same.")
				}
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_GetVacancies(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{},
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:         &sync.Mutex{},
				EmplMu:        &sync.Mutex{},
				ResMu:         &sync.Mutex{},
				VacMu:         &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies:        make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "newpetushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{
					uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd430c8"): Vacancy{
						OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						CompanyName:  "MCDonalds",
						Experience:   "None",
						Profession:   "waiter",
						Position:     "",
						Tasks:        "bring food to costumers",
						Requirements: "middle school education",
						Wage:         "1000 USD",
						Conditions:   "nice team",
						About:        "nice job",
					},
					uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"): Vacancy{
						OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
						CompanyName:  "PETUH",
						Experience:   "None",
						Profession:   "driver",
						Position:     "",
						Tasks:        "drive",
						Requirements: "middle school education",
						Wage:         "50000 RUB",
						Conditions:   "nice team",
						About:        "nice job",
					},
				},
			},
		},
	}

	tests := []struct {
		name string
	}{
		{
			name: "Test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantVacancies := h.UserService.Storage.GetVacancies()

			req, err := http.NewRequest("GET", "/vacancies", nil)

			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/vacancies", h.GetVacancies)
			router.ServeHTTP(rr, req)

			bytes, _ := ioutil.ReadAll(rr.Body)
			var gotVacancies map[uuid.UUID]Vacancy
			json.Unmarshal(bytes, &gotVacancies)

			if rr.Code != http.StatusOK {
				t.Error("status is not ok")
			}

			if !reflect.DeepEqual(wantVacancies, gotVacancies) {
				require.Equal(t, wantVacancies, gotVacancies, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_DeleteVacancy(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"bbbbbbbbbb": {
					ID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:         &sync.Mutex{},
				EmplMu:        &sync.Mutex{},
				ResMu:         &sync.Mutex{},
				VacMu:         &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies: []uuid.UUID{
							uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd430c8"),
						},
					},
					uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "newpetushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies: []uuid.UUID{
							uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"),
						},
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{
					uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd430c8"): Vacancy{
						OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						CompanyName:  "MCDonalds",
						Experience:   "None",
						Profession:   "waiter",
						Position:     "",
						Tasks:        "bring food to costumers",
						Requirements: "middle school education",
						Wage:         "1000 USD",
						Conditions:   "nice team",
						About:        "nice job",
					},
					uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"): Vacancy{
						OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
						CompanyName:  "PETUH",
						Experience:   "None",
						Profession:   "driver",
						Position:     "",
						Tasks:        "drive",
						Requirements: "middle school education",
						Wage:         "50000 RUB",
						Conditions:   "nice team",
						About:        "nice job",
					},
				},
			},
		},
	}
	tests := []struct {
		name             string
		pathArg          string
		cookieValue      string
		wantUnauth       bool
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:        "Test1",
			pathArg:     "11111111-9dad-11d1-1111-00c04fd430c8",
			cookieValue: "bbbbbbbbbb",
			wantFail:    false,
		},
		{
			name:             "Test2",
			pathArg:          "11111111-9dad-11d1-80b1-00c04fd430c8",
			cookieValue:      "bbbbbbbaaaa",
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
		},
		{
			name:        "Test3",
			pathArg:     "11111111-9dad-11d1-80b1-00c04fd430c8",
			cookieValue: "aaaaaaaaaa",
			wantFail:    false,
		},
		{
			name:             "Test4",
			pathArg:          "sdfsfdsd",
			cookieValue:      "bbbb",
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: "Invalid id",
		},
		{
			name:             "Test4",
			pathArg:          "sdfsfdsd",
			cookieValue:      "bbbb",
			wantUnauth:       true,
			wantFail:         true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: "Unauthorized",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/vacancy/%s", tc.pathArg)
			req, err := http.NewRequest("DELETE", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			if !tc.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tc.cookieValue,
				}

				req.AddCookie(&cookie)
			}
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/vacancy/{id}", h.DeleteVacancy)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				gotVacancy, ok := h.UserService.Storage.GetVacancy(uuid.MustParse(tc.pathArg))

				var empVac Vacancy
				if ok != false && gotVacancy != empVac {
					require.Equal(t, empVac, gotVacancy, "The two values should be the same.")
				}
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_PutVacancy(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"bbbbbbbbbb": {
					ID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"cccccccc": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				VacMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies:        make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "newpetushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies: []uuid.UUID{
							uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"),
						},
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{
					uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd430c8"): Vacancy{
						OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						CompanyName:  "MCDonalds",
						Experience:   "None",
						Profession:   "waiter",
						Position:     "",
						Tasks:        "bring food to costumers",
						Requirements: "middle school education",
						Wage:         "1000 USD",
						Conditions:   "nice team",
						About:        "nice job",
					},
					uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"): Vacancy{
						OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
						CompanyName:  "PETUH",
						Experience:   "None",
						Profession:   "driver",
						Position:     "",
						Tasks:        "drive",
						Requirements: "middle school education",
						Wage:         "50000 RUB",
						Conditions:   "nice team",
						About:        "nice job",
					},
				},
			},
		},
	}

	tests := []struct {
		name             string
		cookieValue      string
		pathArg          string
		vacancy          Vacancy
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:        "Test1",
			pathArg:     "11111111-9dad-11d1-1111-00c04fd430c8",
			cookieValue: "bbbbbbbbbb",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
				CompanyName:  "newValue",
				Experience:   "newValue",
				Profession:   "newValue",
				Position:     "newValue",
				Tasks:        "newValue",
				Requirements: "newValue middle school education",
				Wage:         "newValue RUB",
				Conditions:   "newValue team",
				About:        "newValue job",
			},
		},
		{
			name:    "Test2",
			pathArg: "11111111-9dad-11d1-1111-00c04fd430c8",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
				CompanyName:  "newValue",
				Experience:   "newValue",
				Profession:   "newValue",
				Position:     "newValue",
				Tasks:        "newValue",
				Requirements: "newValue middle school education",
				Wage:         "newValue RUB",
				Conditions:   "newValue team",
				About:        "newValue job",
			},
			wantUnauth:       true,
			wantFail:         true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name:        "Test3",
			pathArg:     "11111111-9dad-11d1-1111-00c04fd430c8",
			cookieValue: "cccccccccc",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				CompanyName:  "newValue",
				Experience:   "newValue",
				Profession:   "newValue",
				Position:     "newValue",
				Tasks:        "newValue",
				Requirements: "newValue middle school education",
				Wage:         "newValue RUB",
				Conditions:   "newValue team",
				About:        "newValue job",
			},
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
		},
		{
			name:        "Test4",
			pathArg:     "sdfываыва..,,",
			cookieValue: "bbbbbbbbbb",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
				CompanyName:  "newValue",
				Experience:   "newValue",
				Profession:   "newValue",
				Position:     "newValue",
				Tasks:        "newValue",
				Requirements: "newValue middle school education",
				Wage:         "newValue RUB",
				Conditions:   "newValue team",
				About:        "newValue job",
			},
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidIdMsg,
		},
		{
			name:        "Test5",
			pathArg:     "11111222-9dad-11d1-1111-00c04fd430c8",
			cookieValue: "bbbbbbbbbb",
			vacancy: Vacancy{
				OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
				CompanyName:  "newValue",
				Experience:   "newValue",
				Profession:   "newValue",
				Position:     "newValue",
				Tasks:        "newValue",
				Requirements: "newValue middle school education",
				Wage:         "newValue RUB",
				Conditions:   "newValue team",
				About:        "newValue job",
			},
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
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

			if !tc.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tc.cookieValue,
				}

				req.AddCookie(&cookie)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/vacancy/{id}", h.PutVacancy)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				gotVacancy, _ := h.UserService.Storage.GetVacancy(uuid.MustParse(tc.pathArg))

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}
				if tc.vacancy != gotVacancy {
					require.Equal(t, tc.vacancy, gotVacancy, "The two values should be the same.")
				}
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_CreateSession(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{},
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				VacMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies:        make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "newpetushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage:  map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{},
			},
		},
	}

	tests := []struct {
		name             string
		authInput        UserAuthInput
		invJSON          string
		wantFail         bool
		wantRole         string
		wantStatusCode   int
		wantErrorMessage string
		wantInvJSON      bool
	}{
		{
			name: "Test1",
			authInput: UserAuthInput{
				Email:    "petushki@mail.com",
				Password: "1234",
			},
			wantFail: false,
			wantRole: EmployerStr, //make deep check
		},
		{
			name: "Test2",
			authInput: UserAuthInput{
				Email:    "some_another@mail.com",
				Password: "12345",
			},
			wantFail: false,
			wantRole: SeekerStr, //make deep check
		},
		{
			name: "Test3",
			authInput: UserAuthInput{
				Email:    "some_another@mail.com",
				Password: "1234567",
			},
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: "Invalid password or email",
		},
		{
			name:             "Test4",
			authInput:        UserAuthInput{},
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: "Invalid json",
			wantInvJSON:      true,
			invJSON:          "{'lagin': sdfdfsdf pasword: sdfsdf }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var str string
			if !tt.wantInvJSON {
				wantJSON, _ := json.Marshal(tt.authInput)
				str = string(wantJSON)
			} else {
				str = tt.invJSON
			}

			reader := strings.NewReader(str)

			req, err := http.NewRequest("POST", "/auth", reader)

			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/auth", h.CreateSession)
			router.ServeHTTP(rr, req)

			if !tt.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotRole Role
				json.Unmarshal(bytes, &gotRole)

				require.Equal(t, tt.wantRole, gotRole.Role, "The two values should be the same.")

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tt.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tt.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}

		})
	}
}

func TestHandler_GetSession(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
				"bbbbbbbbbb": {
					ID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				VacMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "newpetushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage:  map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{},
			},
		},
	}
	tests := []struct {
		name             string
		cookieValue      string
		wantRole         string
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:             "Test1",
			cookieValue:      "bbbbbbbbbb",
			wantRole:         EmployerStr,
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name:        "Test2",
			cookieValue: "aaaaaaaaaa",
			wantRole:    SeekerStr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader("") ///why
			req, err := http.NewRequest("GET", "/auth", reader)

			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			if !tt.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tt.cookieValue,
				}

				req.AddCookie(&cookie)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/auth", h.GetSession)
			router.ServeHTTP(rr, req)

			if !tt.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotRole Role
				json.Unmarshal(bytes, &gotRole)

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				if tt.wantRole != gotRole.Role {
					require.Equal(t, tt.wantRole, gotRole.Role, "The two values should be the same.")
				}
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tt.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tt.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_DeleteSession(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"bbbbbbbbbb": {
					ID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"cccccccc": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				VacMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies:        make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "newpetushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies: []uuid.UUID{
							uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"),
						},
					},
				},
				ResumeStorage:  map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{},
			},
		},
	}

	tests := []struct {
		name             string
		cookieValue      string
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:        "Test1",
			cookieValue: "aaaaaaaaaa",
			wantFail:    false,
			wantUnauth:  false,
		},
		{
			name:             "Test2",
			cookieValue:      "aaabbbaaaaa",
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/auth", nil)

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			if !tt.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tt.cookieValue,
				}

				req.AddCookie(&cookie)
			}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/auth", h.DeleteSession)
			router.ServeHTTP(rr, req)

			if !tt.wantFail {
				authInfo, ok := h.AuthService.Storage.Get(tt.cookieValue)

				if ok {
					t.Error("Deleting was failed")
				}
				require.Equal(t, AuthStorageValue{}, authInfo, "The two values should be the same.")
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tt.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tt.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_CreateSeeker(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{},
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:           &sync.Mutex{},
				EmplMu:          &sync.Mutex{},
				ResMu:           &sync.Mutex{},
				VacMu:           &sync.Mutex{},
				SeekerStorage:   map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{},
				ResumeStorage:   map[uuid.UUID]Resume{},
				VacancyStorage:  map[uuid.UUID]Vacancy{},
			},
		},
	}

	tests := []struct {
		name             string
		seekReg          SeekerReg
		wantRole         string
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
		wantInvJSON      bool
		invJSON          string
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
		},
		{
			name: "Test2",
			seekReg: SeekerReg{
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
			},
			wantRole:         SeekerStr,
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: "Email already exists",
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
			wantErrorMessage: "Invalid JSON",
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
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_CreateEmployer(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{},
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:           &sync.Mutex{},
				EmplMu:          &sync.Mutex{},
				ResMu:           &sync.Mutex{},
				VacMu:           &sync.Mutex{},
				SeekerStorage:   map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{},
				ResumeStorage:   map[uuid.UUID]Resume{},
				VacancyStorage:  map[uuid.UUID]Vacancy{},
			},
		},
	}

	tests := []struct {
		name             string
		emplReg          EmployerReg
		wantRole         string
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
		wantInvJSON      bool
		invJSON          string
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
				City:             "Petushki",
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
				City:             "Petushki",
				EmplNum:          "322",
			},
			wantRole:         EmployerStr,
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: "Email already exists",
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
				City:             "Petushki",
				EmplNum:          "322",
			},
			wantRole:         EmployerStr,
			wantFail:         true,
			wantInvJSON:      true,
			invJSON:          "{sfsdf: some email, login: password: sdfdf}",
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: "Invalid JSON",
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
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_CreateResume(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"cccccccc": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
				"cfv": {
					ID:      uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				VacMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"): {
						Email:      "dmasis@mail.com",
						FirstName:  "Dima",
						SecondName: "Sidorov",
						Password:   "12345",
						Resumes:    make([]uuid.UUID, 0),
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage:  map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{},
			},
		},
	}

	tests := []struct {
		name             string
		cookieValue      string
		resume           Resume
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
		wantInvJSON      bool
		invJSON          string
	}{
		{
			name:        "Test1",
			cookieValue: "cccccccc",
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
		},
		{
			name:        "Test2",
			cookieValue: "aaaaaaaaaa",
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
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
		},
		{
			name:        "Test3",
			cookieValue: "cbfdgdfg",
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
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name:             "Test4",
			cookieValue:      "cbfdgdfg",
			resume:           Resume{},
			wantFail:         true,
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
			wantInvJSON:      true,
			invJSON:          "{testx: fdsfsdf, fdsfsdf'sdfsdf / fdsfsdf}",
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

			req, err := http.NewRequest("POST", "/resume", reader)

			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			if !tc.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tc.cookieValue,
				}

				req.AddCookie(&cookie)
			}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/resume", h.CreateResume)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var id Id
				json.Unmarshal(bytes, &id)
				gotResume, _ := h.UserService.Storage.GetResume(uuid.MustParse(id.Id))

				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				require.Equal(t, tc.resume, gotResume, "The two values should be the same.")
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_DeleteResume(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"cccccccc": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
				"cfv": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				VacMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes: []uuid.UUID{
							uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"),
						},
					},
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						Email:      "dmasis@mail.com",
						FirstName:  "Vova",
						SecondName: "Sidorov",
						Password:   "12345",
						Resumes: []uuid.UUID{
							uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
						},
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{
					uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"): {
						OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						FirstName:   "Vova",
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
					uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"): {
						OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
						FirstName:   "Petr",
						SecondName:  "Zyablikov",
						City:        "Moscow",
						PhoneNumber: "12345678910",
						BirthDate:   "1994-21-08",
						Sex:         "male",
						Citizenship: "Russia",
						Experience:  "8 years",
						Profession:  "programmer",
						Position:    "senior",
						Wage:        "100500",
						Education:   "MSU",
						About:       "Hello employer",
					},
				},
				VacancyStorage: map[uuid.UUID]Vacancy{},
			},
		},
	}
	tests := []struct {
		name             string
		pathArg          string
		cookieValue      string
		wantUnauth       bool
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:             "Test1",
			pathArg:          "7aa7b810-9dad-11d1-72b5-04c04fd430c8",
			cookieValue:      "aaaaaaaaaa",
			wantUnauth:       false,
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
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
			name:        "Test3",
			pathArg:     "7aa7b810-9dad-11d1-72b5-04c04fd430c8",
			cookieValue: "cccccccc",
			wantFail:    false,
		},
		{
			name:             "Test3",
			pathArg:          "ываывадлд",
			cookieValue:      "cccccccc",
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidIdMsg,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/resume/%s", tc.pathArg)
			req, err := http.NewRequest("DELETE", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			if !tc.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tc.cookieValue,
				}

				req.AddCookie(&cookie)
			}
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/resume/{id}", h.DeleteResume)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}

				gotResume, ok := h.UserService.Storage.GetResume(uuid.MustParse(tc.pathArg))

				var empResume Resume
				if ok != false && gotResume != empResume {
					require.Equal(t, gotResume, empResume, "The two values should be the same.")
				}
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_GetResume(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: make(map[string]AuthStorageValue),
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						Email:      "some@mail.com",
						FirstName:  "Vova",
						SecondName: "Zyablikov",
						Password:   "1234",
						Resumes: []uuid.UUID{
							uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd435c8"),
						},
					},
					uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"): {
						Email:      "third@mail.com",
						FirstName:  "Petr",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes: []uuid.UUID{
							uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"),
						},
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{},
				ResumeStorage: map[uuid.UUID]Resume{
					uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd435c8"): {
						OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						FirstName:   "Vova",
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
					uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"): {
						OwnerID:     uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
						FirstName:   "Vova",
						SecondName:  "Zyablikov",
						City:        "Moscow",
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
				},
			},
		},
	}

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
			pathArg:          "222222-9dad-11d1-80b1-00c04fd435c8",
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidIdMsg,
		},
		{
			name:             "Test3",
			pathArg:          "фвапвапвпа_аыва",
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

				wantResume, _ := h.UserService.Storage.GetResume(uuid.MustParse(tc.pathArg))

				require.Equal(t, wantResume, gotResume, "The two values should be the same.")
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_GetVacancy(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: make(map[string]AuthStorageValue),
			Mu:      &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:         &sync.Mutex{},
				EmplMu:        &sync.Mutex{},
				ResMu:         &sync.Mutex{},
				VacMu:         &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies: []uuid.UUID{
							uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd430c8"),
						},
					},
					uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"): {
						CompanyName:      "Petushki",
						Site:             "petushki.com",
						Email:            "newpetushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "1488",
						Vacancies: []uuid.UUID{
							uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"),
						},
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{},
				VacancyStorage: map[uuid.UUID]Vacancy{
					uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd430c8"): Vacancy{
						OwnerID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
						CompanyName:  "MCDonalds",
						Experience:   "None",
						Profession:   "waiter",
						Position:     "",
						Tasks:        "bring food to costumers",
						Requirements: "middle school education",
						Wage:         "1000 USD",
						Conditions:   "nice team",
						About:        "nice job",
					},
					uuid.MustParse("11111111-9dad-11d1-1111-00c04fd430c8"): Vacancy{
						OwnerID:      uuid.MustParse("6ba7b810-9bbb-1111-1111-00c04fd430c8"),
						CompanyName:  "PETUH",
						Experience:   "None",
						Profession:   "driver",
						Position:     "",
						Tasks:        "drive",
						Requirements: "middle school education",
						Wage:         "50000 RUB",
						Conditions:   "nice team",
						About:        "nice job",
					},
				},
			},
		},
	}

	tests := []struct {
		name             string
		pathArg          string
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:     "Test1",
			pathArg:  "11111111-9dad-11d1-80b1-00c04fd430c8",
			wantFail: false,
		},
		{
			name:             "Test2",
			pathArg:          "222222-9dad-11d1-80b1-00c04fd435c8",
			wantFail:         true,
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: InvalidIdMsg,
		},
		{
			name:             "Test3",
			pathArg:          "фвапвапвпа_аыва",
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

				wantVacancy, _ := h.UserService.Storage.GetVacancy(uuid.MustParse(tc.pathArg))

				require.Equal(t, wantVacancy, gotVacancy, "The two values should be the same.")
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}

func TestHandler_PutResume(t *testing.T) {
	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: map[string]AuthStorageValue{
				"aaaaaaaaaa": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    EmployerStr,
				},
				"cccccccc": {
					ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Role:    SeekerStr,
				},
			},
			Mu: &sync.Mutex{},
		},
	}

	h := &Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:  &sync.Mutex{},
				EmplMu: &sync.Mutex{},
				ResMu:  &sync.Mutex{},
				VacMu:  &sync.Mutex{},
				SeekerStorage: map[uuid.UUID]Seeker{
					uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"): {
						Email:      "some_another@mail.com",
						FirstName:  "Petya",
						SecondName: "Zyablikov",
						Password:   "12345",
						Resumes: []uuid.UUID{
							uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
						},
					},
				},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName:      "MCDonalds",
						Site:             "petushki.com",
						Email:            "petushki@mail.com",
						FirstName:        "Vova",
						SecondName:       "Zyablikov",
						Password:         "1234",
						PhoneNumber:      "12345678911",
						ExtraPhoneNumber: "12345678910",
						City:             "Petushki",
						EmplNum:          "322",
						Vacancies:        make([]uuid.UUID, 0),
					},
				},
				ResumeStorage: map[uuid.UUID]Resume{
					uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"): {
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
				VacancyStorage: map[uuid.UUID]Vacancy{},
			},
		},
	}

	tests := []struct {
		name             string
		cookieValue      string
		pathArg          string
		resume           Resume
		wantFail         bool
		wantUnauth       bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			name:        "Test1",
			cookieValue: "cccccccc",
			pathArg:     "7ba7b810-9dad-12d1-80b1-00c04fd430c8",
			resume: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				FirstName:   "Petya",
				SecondName:  "Zyablikov",
				City:        "Kaliningrad",
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
		},
		{
			name:             "Test1",
			cookieValue:      "aaaaaaaaaa",
			pathArg:          "7ba7b810-9dad-12d1-80b1-00c04fd430c8",
			resume:           Resume{},
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
		},
		{
			name:             "Test1",
			cookieValue:      "",
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

			if !tc.wantUnauth {
				cookie := http.Cookie{
					Name:  auth.CookieName,
					Value: tc.cookieValue,
				}

				req.AddCookie(&cookie)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/resume/{id}", h.PutResume)
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				gotResume, _ := h.UserService.Storage.GetResume(uuid.MustParse(tc.pathArg))

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
				require.Equal(t, tc.wantErrorMessage, gotError.Error, "The two values should be the same.")
			}
		})
	}
}
