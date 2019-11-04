package handler


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
				require.Equal(t, tc.wantErrorMessage, gotError.Message, "The two values should be the same.")
			}
		})
	}
}