package server

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	auth_rep "2019_2_IBAT/internal/pkg/auth/repository"
	auth_serv "2019_2_IBAT/internal/pkg/auth/service"

	"2019_2_IBAT/internal/pkg/handler"
	usRep "2019_2_IBAT/internal/pkg/users/repository"
	usServ "2019_2_IBAT/internal/pkg/users/service"

	"2019_2_IBAT/internal/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"

	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"golang.org/x/crypto/acme/autocert"
)

// type Server struct {
// 	Router *mux.Router
// }

const staticDir = "/media/vltim/img"

func NewRouter() (*mux.Router, error) {
	// server := new(Server)

	router := mux.NewRouter()

	redisAddr := flag.String("redisServer", ":6379", "")

	aS := auth_serv.AuthService{
		Storage: auth_rep.NewSessionManager(auth_rep.RedNewPool(*redisAddr)),
	}

	uS := usServ.UserService{
		Storage: &usRep.DBUserStorage{
			DbConn: OpenSqlxViaPgxConnPool(),
		},
	}

	h := handler.Handler{
		InternalDir: staticDir,
		AuthService: &aS,
		UserService: &uS,
	}

	loger := middleware.NewLogger()
	// AccessLogOut := new(middleware.AccessLogger)
	// AccessLogOut.StdLogger = log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile)

	router.Use(loger.AccessLogMiddleware)
	router.Use(middleware.CorsMiddleware)
	router.Use(aS.AuthMiddleware)
	router.Use(middleware.CSRFMiddleware)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	router.HandleFunc("/upload", h.UploadFile()).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/auth", h.CreateSession).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/auth", h.GetSession).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth", h.DeleteSession).Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/seeker", h.CreateSeeker).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/seeker", h.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/seeker", h.PutUser).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/seeker/{id}", h.GetSeekerById).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/employer", h.CreateEmployer).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/employer", h.DeleteUser).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/employer", h.PutUser).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/employer/{id}", h.GetEmployerById).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/profile", h.GetUser).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/resume", h.CreateResume).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/resume/{id}", h.DeleteResume).Methods(http.MethodDelete, http.MethodOptions) //test
	router.HandleFunc("/resume/{id}", h.PutResume).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/resume/{id}", h.GetResume).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/vacancy", h.CreateVacancy).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/vacancy/{id}", h.DeleteVacancy).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/vacancy/{id}", h.GetVacancy).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/vacancy/{id}", h.PutVacancy).Methods(http.MethodPut, http.MethodOptions)

	router.HandleFunc("/employers", h.GetEmployers).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/resumes", h.GetResumes).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/vacancies", h.GetVacancies).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/responds", h.GetResponds).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/respond", h.CreateRespond).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/favorite_vacancies", h.GetFavoriteVacancies).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/favorite_vacancy/{id}", h.CreateFavorite).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/favorite_vacancy/{id}", h.DeleteFavoriteVacancy).Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/tags", h.GetTags).Methods(http.MethodGet, http.MethodOptions)

	return router, nil
}

func RunServer() {

	router, _ := NewRouter()
	httpsSrv := &http.Server{
		Handler: router,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	inProduction := true

	if inProduction {
		// Note: use a sensible value for data directory
		// this is where cached certificates are stored
		dataDir := "."
		hostPolicy := func(ctx context.Context, host string) error {
			// Note: change to your real domain
			allowedHost := "*"
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}

		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		// log.Fatal(httpsSrv.ListenAndServe())
		// err := httpsSrv.ListenAndServeTLS("", "")
		// if err != nil {
		// 	log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
		// }
	} else {

		httpsSrv.Addr = ":8080"

		// err := httpsSrv.ListenAndServe()
		// log.Fatal(httpsSrv.ListenAndServe())
		// if err != nil {
		// 	log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		// }
	}
	log.Fatal(httpsSrv.ListenAndServe())
}

func OpenSqlxViaPgxConnPool() *sqlx.DB {
	connConfig := pgx.ConnConfig{
		Host:     "localhost",
		Database: "hh",
		User:     "postgres",
		Password: "newPassword",
	}
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AfterConnect:   nil,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to create connections pool")
	}

	nativeDB := stdlib.OpenDBFromPool(connPool)

	log.Println("OpenSqlxViaPgxConnPool: the connection was created")
	return sqlx.NewDb(nativeDB, "pgx")
}

// git commit -m "https cert resolution added"
