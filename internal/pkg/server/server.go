package server

import (
	"flag"
	"log"
	"net/http"
	"time"

	"2019_2_IBAT/internal/pkg/handler"
	usRep "2019_2_IBAT/internal/pkg/users/repository"
	usServ "2019_2_IBAT/internal/pkg/users/service"

	"2019_2_IBAT/internal/pkg/auth/session"
	. "2019_2_IBAT/internal/pkg/interfaces"
	recRep "2019_2_IBAT/internal/pkg/recommends/repository"
	recServ "2019_2_IBAT/internal/pkg/recommends/service"

	"2019_2_IBAT/internal/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"google.golang.org/grpc"

	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const staticDir = "/media/vltim/img"
const NotifChanSize = 50

var addr = flag.String("listen-address", ":8080",
	"The address to listen on for HTTP requests.")

func NewRouter() (*mux.Router, error) {
	// flag.Parse()

	// usersRegistered := prometheus.NewCounter(
	// 	prometheus.CounterOpts{
	// 		Name: "users_registered",
	// 	})
	// prometheus.MustRegister(usersRegistered)

	// usersOnline := prometheus.NewGauge(
	// 	prometheus.GaugeOpts{
	// 		Name: "users_online",
	// 	})
	// prometheus.MustRegister(usersOnline)

	// requestProcessingTimeSummaryMs := prometheus.NewSummary(
	// 	prometheus.SummaryOpts{
	// 		Name:       "request_processing_time_summary_ms",
	// 		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	// 	})
	// prometheus.MustRegister(requestProcessingTimeSummaryMs)

	// requestProcessingTimeHistogramMs := prometheus.NewHistogram(
	// 	prometheus.HistogramOpts{
	// 		Name:    "request_processing_time_histogram_ms",
	// 		Buckets: prometheus.LinearBuckets(0, 10, 20),
	// 	})
	// prometheus.MustRegister(requestProcessingTimeHistogramMs)

	// go func() {
	// 	for {
	// 		usersRegistered.Inc() // or: Add(5)
	// 		time.Sleep(1000 * time.Millisecond)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		for i := 0; i < 10000; i++ {
	// 			usersOnline.Set(float64(i)) // or: Inc(), Dec(), Add(5), Dec(5)
	// 			time.Sleep(10 * time.Millisecond)
	// 		}
	// 	}
	// }()

	// go func() {
	// 	src := rand.NewSource(time.Now().UnixNano())
	// 	rnd := rand.New(src)
	// 	for {
	// 		obs := float64(100 + rnd.Intn(30))
	// 		requestProcessingTimeSummaryMs.Observe(obs)
	// 		requestProcessingTimeHistogramMs.Observe(obs)
	// 		time.Sleep(10 * time.Millisecond)
	// 	}
	// }()

	router := mux.NewRouter()

	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
		return router, err
	}

	// defer grcpConn.Close()

	sessManager := session.NewServiceClient(grcpConn)

	rS := recServ.Service{
		Storage: &recRep.DBRecommendsStorage{
			DbConn: OpenSqlxViaPgxConnPool(),
		},
	}

	uS := usServ.UserService{
		Storage: &usRep.DBUserStorage{
			DbConn: OpenSqlxViaPgxConnPool(),
		},
		RecomService: rS,
		NotifChan:    make(chan NotifStruct, NotifChanSize),
	}

	h := handler.Handler{
		InternalDir: staticDir,
		AuthService: sessManager,
		UserService: &uS,
		WsConnects:  map[string]Connections{},
	}

	//should remove
	go uS.Notifications(h.WsConnects)

	loger := middleware.NewLogger()
	authMiddleware := middleware.AuthMiddlewareGenerator(sessManager)

	router.Use(loger.AccessLogMiddleware)
	// router.Use(middleware.CorsMiddleware)
	router.Use(authMiddleware)
	// router.Use(middleware.CSRFMiddleware)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	router.HandleFunc("/upload", h.UploadFile()).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/auth", h.CreateSession).Methods(http.MethodPost) //, http.MethodOptions)
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

	router.HandleFunc("/notifications", h.Notifications)

	// router.Handle("/metrics", promhttp.Handler())
	return router, nil
}

// func RunServer() {
// 	router, _ := NewRouter()
// 	httpsSrv := &http.Server{
// 		Handler: router,
// 		// Good practice: enforce timeouts for servers you create!
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}
// 	inProduction := true

// 	if inProduction {
// 		// Note: use a sensible value for data directory
// 		// this is where cached certificates are stored
// 		dataDir := "."
// 		hostPolicy := func(ctx context.Context, host string) error {
// 			allowedHost := "tko.vladimir.fvds.ru"
// 			if host == allowedHost {
// 				return nil
// 			}
// 			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
// 		}

// 		m := &autocert.Manager{
// 			Prompt:     autocert.AcceptTOS,
// 			HostPolicy: hostPolicy,
// 			Cache:      autocert.DirCache(dataDir),
// 		}

// 		httpsSrv.Addr = ":443"
// 		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

// 		log.Fatal(httpsSrv.ListenAndServeTLS("", ""))
// 	} else {

// 		httpsSrv.Addr = ":8080"

// 		log.Fatal(httpsSrv.ListenAndServe())
// 	}
// }

func RunServer() {
	router, err := NewRouter()
	if err != nil {
		log.Fatal("Failed to create router")
	}
	// log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router))
	log.Fatal(http.ListenAndServe(":8080", router))
	// go uS.Notifications(h.WsConnects)

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
