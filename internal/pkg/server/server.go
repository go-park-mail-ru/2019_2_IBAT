package server

import (
	"log"
	"net/http"
	"sync"

	"hh_workspace/2019_2_IBAT/internal/pkg/auth"
	"hh_workspace/2019_2_IBAT/internal/pkg/handler"
	. "hh_workspace/2019_2_IBAT/internal/pkg/interfaces"
	"hh_workspace/2019_2_IBAT/internal/pkg/users"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func NewServer() (*Server, error) {
	server := new(Server)

	router := mux.NewRouter()

	ah := auth.Handler{
		Storage: auth.MapAuthStorage{
			Storage: make(map[string]AuthStorageValue),
			Mu:      &sync.Mutex{},
		},
	}

	h := handler.Handler{
		AuthHandler: ah,
		UserControler: users.Controler{
			Storage: &users.MapUserStorage{
				SekMu:           &sync.Mutex{},
				EmplMu:          &sync.Mutex{},
				ResMu:           &sync.Mutex{},
				SeekerStorage:   map[uuid.UUID]Seeker{},
				EmployerStorage: map[uuid.UUID]Employer{}, //implement through new_handler()
				ResumeStorage:   map[uuid.UUID]Resume{},
			},
		},
	}

	router.HandleFunc("/auth", h.CreateSession).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/auth", h.DeleteSession).Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/seeker", h.CreateSeeker).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/seeker", h.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/seeker", h.GetSeeker).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/seeker", h.PutUser).Methods(http.MethodPut, http.MethodOptions)

	router.HandleFunc("/employer", h.CreateEmployer).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/employer", h.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/employer", h.GetEmployer).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/employer", h.PutUser).Methods(http.MethodPut, http.MethodOptions)

	router.HandleFunc("/resume", h.CreateResume).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/resume/{id}", h.DeleteResume).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/resume/{id}", h.PutResume).Methods(http.MethodPut, http.MethodOptions) // extra
	router.HandleFunc("/resume/{id}", h.GetResume).Methods(http.MethodGet, http.MethodOptions)

	// router.HandleFunc("/vacancy/{id}", h.HandleGetVacancy).Methods(http.MethodGet, http.MethodOptions)

	// router.HandleFunc("/vacancies", h.HandleGet).Methods(http.MethodGet, http.MethodOptions)
	// router.HandleFunc("/resumes", h.HandleGet).Methods(http.MethodGet, http.MethodOptions)
	// router.HandleFunc("/emloyers", h.HandleGet).Methods(http.MethodGet, http.MethodOptions)

	// router.HandleFunc("/vacancy", h.HandleVacancyResume).Methods(http.MethodPost, http.MethodOptions)

	server.Router = router

	return server, nil
}

func (server *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", server.Router))
}
