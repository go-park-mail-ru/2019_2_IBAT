package handler

import (
	"encoding/json"
	"fmt"
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"2019_2_IBAT/internal/pkg/users"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

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
						Number:      "12345678910",
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
						Number:      "12345678910",
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
						CompanyName: "Petushki",
						Site:        "petushki.com",
						Email:       "petushki@mail.com",
						FirstName:   "Vova",
						SecondName:  "Zyablikov",
						Password:    "1234",
						Number:      "12345678911",
						ExtraNumber: "12345678910",
						City:        "Petushki",
						EmplNum:     1488,
						Vacancies:   make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"): {
						CompanyName: "BMSTU",
						Site:        "bmstu.ru",
						Email:       "bmstu@mail.com",
						FirstName:   "Tolya",
						SecondName:  "Alex",
						Password:    "1234",
						Number:      "12345678911",
						ExtraNumber: "12345678910",
						City:        "Moscow",
						EmplNum:     1830,
						Vacancies:   make([]uuid.UUID, 0),
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
						CompanyName: "Petushki",
						Site:        "petushki.com",
						Email:       "petushki@mail.com",
						FirstName:   "Vova",
						SecondName:  "Zyablikov",
						Password:    "1234",
						Number:      "12345678911",
						ExtraNumber: "12345678910",
						City:        "Petushki",
						EmplNum:     1488,
						Vacancies:   make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"): {
						CompanyName: "BMSTU",
						Site:        "bmstu.ru",
						Email:       "bmstu@mail.com",
						FirstName:   "Tolya",
						SecondName:  "Alex",
						Password:    "1234",
						Number:      "12345678911",
						ExtraNumber: "12345678910",
						City:        "Moscow",
						EmplNum:     1830,
						Vacancies:   make([]uuid.UUID, 0),
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
						Number:      "12345678910",
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
						Number:      "12345678910",
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
					Class:   EmployerStr,
				},
				"aaaaaaaaab": {
					ID:      uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"),
					Expires: time.Now().In(auth.Loc).Add(24 * time.Hour).Format(auth.TimeFormat),
					Class:   EmployerStr,
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
				SeekerStorage: map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{
					uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
						CompanyName: "Petushki",
						Site:        "petushki.com",
						Email:       "petushki@mail.com",
						FirstName:   "Vova",
						SecondName:  "Zyablikov",
						Password:    "1234",
						Number:      "12345678911",
						ExtraNumber: "12345678910",
						City:        "Petushki",
						EmplNum:     1488,
						Vacancies:   make([]uuid.UUID, 0),
					},
					uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"): {
						CompanyName: "BMSTU",
						Site:        "bmstu.ru",
						Email:       "bmstu@mail.com",
						FirstName:   "Tolya",
						SecondName:  "Alex",
						Password:    "1234",
						Number:      "12345678911",
						ExtraNumber: "12345678910",
						City:        "Moscow",
						EmplNum:     1830,
						Vacancies:   make([]uuid.UUID, 0),
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
