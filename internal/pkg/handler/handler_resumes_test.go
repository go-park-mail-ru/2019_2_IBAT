package handler

<<<<<<< HEAD
import (
	mock_users "2019_2_IBAT/internal/pkg/handler/mock_users"
	"context"
	"sync"

	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/golang/mock/gomock"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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
		{
			ID:          uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"),
			OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
=======

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"github.com/golang/mock/gomock"
	"encoding/json"
	"io/ioutil"
	mock_users "2019_2_IBAT/internal/pkg/handler/mock_users"
	"fmt"
	"net/http"
	"github.com/google/uuid"

)

// (1) Define an interface that you wish to mock.
//       type MyInterface interface {
//         SomeMethod(x int64, y string)
//       }
// (2) Use mockgen to generate a mock from the interface.
// (3) Use the mock in a test:
//       func TestMyThing(t *testing.T) {
//         mockCtrl := gomock.NewController(t)
//         defer mockCtrl.Finish()

//         mockObj := something.NewMockMyInterface(mockCtrl)
//         mockObj.EXPECT().SomeMethod(4, "blah")
//         // pass mockObj to a real object and play with it.
//       }


// func (h *Handler) GetResume(w http.ResponseWriter, r *http.Request) { //+
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	log.Println("Handle GetResume: start")

// 	resId, err := uuid.Parse(mux.Vars(r)["id"])

// 	if err != nil {
// 		log.Println("Handle GetResume: invalid id")
// 		w.WriteHeader(http.StatusBadRequest)
// 		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
// 		w.Write([]byte(errJSON))
// 		return
// 	}

// 	resume, err := h.UserService.GetResume(resId)

// 	if err != nil {
// 		log.Println("Handle GetResume: failed to get resume")
// 		w.WriteHeader(http.StatusBadRequest)
// 		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
// 		w.Write([]byte(errJSON))
// 		return
// 	}

// 	resumeJSON, _ := json.Marshal(resume)

// 	w.Write([]byte(resumeJSON))
// }

func TestHandler_GetResume(t *testing.T) {
	// ah := auth.AuthService{
	// 	Storage: auth.MapAuthStorage{
	// 		Storage: make(map[string]AuthStorageValue),
	// 		Mu:      &sync.Mutex{},
	// 	},
	// }

	// h := &Handler{
	// 	AuthService: ah,
	// 	UserService: users.UserService{
	// 		Storage: &users.MapUserStorage{
	// 			SekMu:  &sync.Mutex{},
	// 			EmplMu: &sync.Mutex{},
	// 			ResMu:  &sync.Mutex{},
	// 			SeekerStorage: map[uuid.UUID]Seeker{
	// 				uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
	// 					Email:      "some@mail.com",
	// 					FirstName:  "Vova",
	// 					SecondName: "Zyablikov",
	// 					Password:   "1234",
	// 					Resumes: []uuid.UUID{
	// 						uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd435c8"),
	// 					},
	// 				},
	// 				uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"): {
	// 					Email:      "third@mail.com",
	// 					FirstName:  "Petr",
	// 					SecondName: "Zyablikov",
	// 					Password:   "12345",
	// 					Resumes: []uuid.UUID{
	// 						uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"),
	// 					},
	// 				},
	// 			},
	// 			EmployerStorage: map[uuid.UUID]Employer{},
	// 			ResumeStorage: map[uuid.UUID]Resume{
	// 				uuid.MustParse("11111111-9dad-11d1-80b1-00c04fd435c8"): {
	// 					OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
	// 					FirstName:   "Vova",
	// 					SecondName:  "Zyablikov",
	// 					City:        "Moscow",
	// 					PhoneNumber: "12345678910",
	// 					BirthDate:   "1994-21-08",
	// 					Sex:         "male",
	// 					Citizenship: "Russia",
	// 					Experience:  "7 years",
	// 					Profession:  "programmer",
	// 					Position:    "middle",
	// 					Wage:        "100500",
	// 					Education:   "MSU",
	// 					About:       "Hello employer",
	// 				},
	// 				uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"): {
	// 					OwnerID:     uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
	// 					FirstName:   "Vova",
	// 					SecondName:  "Zyablikov",
	// 					City:        "Moscow",
	// 					PhoneNumber: "12345678910",
	// 					BirthDate:   "1994-21-08",
	// 					Sex:         "male",
	// 					Citizenship: "Ukraine",
	// 					Experience:  "7 years",
	// 					Profession:  "programmer",
	// 					Position:    "middle",
	// 					Wage:        "100500",
	// 					Education:   "MSU",
	// 					About:       "Hello employer",
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	
	userService := mock_users.NewMockService(mockCtrl)
	h := Handler{
		UserService: userService,
	}

	want_resumes := []Resume{
		{
			ID: uuid.MustParse("22222222-9dad-11d1-80b1-00c04fd435c8"),
			OwnerID:     uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
>>>>>>> 0ffc1e987c6bc6b935425bf07d111edecc90a0f3
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
	}

<<<<<<< HEAD
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
		cookieValue      string
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
			record: AuthStorageValue{
				ID:      uuid.MustParse("6ba7b810-9dad-11d1-0000-00004fd430c8"),
				Role:    SeekerStr,
				Expires: time.Now().In(Loc).Add(24 * time.Hour).Format(TimeFormat),
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
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: BadRequestMsg,
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

			req, _ := http.NewRequest("POST", "/resume", reader)

			id1 := uuid.New()
			log.Printf("id1 generated: %s\n", id1)

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

				log.Println("TEST Create resume req:")
				log.Println(req)
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
					log.Printf("id.Id after unmarshaling: %s\n", id.Id)
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
=======
	tests := []struct {
		name             string
		pathArg          string
>>>>>>> 0ffc1e987c6bc6b935425bf07d111edecc90a0f3
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
<<<<<<< HEAD
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
=======
			name:     "Test1",
			pathArg:  "22222222-9dad-11d1-80b1-00c04fd435c8",
			wantFail: false,
		},
		// {
		// 	name:             "Test2",
		// 	pathArg:          "222222-9dad-11d1-80b1-00c04fd435c8",
		// 	wantFail:         true,
		// 	wantStatusCode:   http.StatusBadRequest,
		// 	wantErrorMessage: InvalidIdMsg,
		// },
		// {
		// 	name:             "Test3",
		// 	pathArg:          "фвапвапвпа_аыва",
		// 	wantFail:         true,
		// 	wantStatusCode:   http.StatusBadRequest,
		// 	wantErrorMessage: InvalidIdMsg,
		// },
>>>>>>> 0ffc1e987c6bc6b935425bf07d111edecc90a0f3
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := fmt.Sprintf("/resume/%s", tc.pathArg)
<<<<<<< HEAD
			req, err := http.NewRequest("DELETE", path, nil)
=======
			req, err := http.NewRequest("GET", path, nil)
>>>>>>> 0ffc1e987c6bc6b935425bf07d111edecc90a0f3
			if err != nil {
				t.Fatal(err)
			}

<<<<<<< HEAD
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
=======
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/resume/{id}", h.GetResume)
>>>>>>> 0ffc1e987c6bc6b935425bf07d111edecc90a0f3
			router.ServeHTTP(rr, req)

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}
<<<<<<< HEAD

				gotResume, ok := h.UserService.Storage.GetResume(uuid.MustParse(tc.pathArg))

				var empResume Resume
				if ok != false && gotResume != empResume {
					require.Equal(t, gotResume, empResume, "The two values should be the same.")
				}
=======
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotResume Resume
				json.Unmarshal(bytes, &gotResume)

				wantResume, _ := h.UserService.Storage.GetResume(uuid.MustParse(tc.pathArg))

				require.Equal(t, wantResume, gotResume, "The two values should be the same.")
>>>>>>> 0ffc1e987c6bc6b935425bf07d111edecc90a0f3
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				json.Unmarshal(bytes, &gotError)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Message, "The two values should be the same.")
			}
		})
	}
<<<<<<< HEAD
}
=======
}
>>>>>>> 0ffc1e987c6bc6b935425bf07d111edecc90a0f3
