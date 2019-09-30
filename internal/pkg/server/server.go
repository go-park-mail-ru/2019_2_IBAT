package server

import (
	"log"
	"net/http"
	"sync"

	"2019_2_IBAT/internal/pkg/auth"
	"2019_2_IBAT/internal/pkg/handler"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"2019_2_IBAT/internal/pkg/users"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func NewServer() (*Server, error) {
	server := new(Server)

	router := mux.NewRouter()

	ah := auth.AuthService{
		Storage: auth.MapAuthStorage{
			Storage: make(map[string]AuthStorageValue),
			Mu:      &sync.Mutex{},
		},
	}

	h := handler.Handler{
		AuthService: ah,
		UserService: users.UserService{
			Storage: &users.MapUserStorage{
				SekMu:           &sync.Mutex{},
				EmplMu:          &sync.Mutex{},
				ResMu:           &sync.Mutex{},
				VacMu:           &sync.Mutex{},
				SeekerStorage:   map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{}, //implement through new_handler()
				ResumeStorage:   map[uuid.UUID]Resume{},
				VacancyStorage:  map[uuid.UUID]Vacancy{},
			},
		},
	}

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/tmp/img"))))

	router.Use(handler.CorsMiddleware)

	router.HandleFunc("/upload", h.UploadFile()).Methods(http.MethodPost, http.MethodOptions)

	// staticHandler := http.FileServer(http.Dir("/tmp/img"))

	// router.Handle("/static/{id}", http.StripPrefix("/static/", staticHandler)).Methods(http.MethodGet, http.MethodOptions)
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/tmp/img"))))

	router.HandleFunc("/auth", h.CreateSession).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/auth", h.GetSession).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth", h.DeleteSession).Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/seeker", h.CreateSeeker).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/seeker", h.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/seeker", h.GetSeeker).Methods(http.MethodGet, http.MethodOptions) //use for getting all seeker info including password
	router.HandleFunc("/seeker", h.PutUser).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/seeker/{id}", h.GetSeekerById).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/employer", h.CreateEmployer).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/employer", h.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/employer", h.GetEmployer).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/employer", h.PutUser).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/employer/{id}", h.GetEmployerById).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/resume", h.CreateResume).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/resume/{id}", h.DeleteResume).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/resume/{id}", h.PutResume).Methods(http.MethodPut, http.MethodOptions) // extra
	router.HandleFunc("/resume/{id}", h.GetResume).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/vacancy", h.CreateVacancy).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/vacancy/{id}", h.DeleteVacancy).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/vacancy/{id}", h.GetVacancy).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/vacancy/{id}", h.PutVacancy).Methods(http.MethodPut, http.MethodOptions)

	router.HandleFunc("/employers", h.GetEmployers).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/resumes", h.GetResumes).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/vacancies", h.GetVacancies).Methods(http.MethodGet, http.MethodOptions)

	server.Router = router

	return server, nil
}

func (server *Server) Run() {
	// log.Fatal(http.ListenAndServe(":8080", server.Router))
	log.Fatal(http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", server.Router))
}
