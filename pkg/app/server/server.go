package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"2019_2_IBAT/pkg/app/auth/session"
	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"

	"2019_2_IBAT/pkg/app/server/handler"
	usRep "2019_2_IBAT/pkg/app/users/repository"
	usServ "2019_2_IBAT/pkg/app/users/service"
	"2019_2_IBAT/pkg/pkg/config"
	"2019_2_IBAT/pkg/pkg/db_connect"
	"2019_2_IBAT/pkg/pkg/middleware"
)

const staticDir = "/media/vltim/img"

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

	authGrcpConn, err := grpc.Dial(
		config.AuthHostname+":"+strconv.Itoa(config.AuthServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to authGrcpConn")
		return router, err
	}
	fmt.Printf("authGrcpConn: RedisHostname %s\n", config.RedisHostname)

	sessManager := session.NewServiceClient(authGrcpConn)

	recomsGrcpConn, err := grpc.Dial(
		"127.0.0.1:"+strconv.Itoa(config.RecommendsServicePort),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("cant connect to recomsGrcpConn")
		return router, err
	}

	recommsManager := recomsproto.NewServiceClient(recomsGrcpConn)

	notifGrcpConn, err := grpc.Dial(
		"127.0.0.1:"+strconv.Itoa(config.NotifsServicePort),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("cant connect to recomsGrcpConn")
		return router, err
	}

	notifManager := notifsproto.NewServiceClient(notifGrcpConn)

	uS := usServ.UserService{
		Storage: &usRep.DBUserStorage{
			DbConn: db_connect.OpenSqlxViaPgxConnPool(),
		},
		RecomService: recommsManager,
		NotifService: notifManager,
	}

	h := handler.Handler{
		InternalDir: staticDir,
		AuthService: sessManager,
		UserService: &uS,
	}

	loger := middleware.NewLogger()
	authMiddleware := middleware.AuthMiddlewareGenerator(sessManager)

	router = router.PathPrefix("/api/").Subrouter()

	router.Use(loger.AccessLogMiddleware)
	router.Use(middleware.CorsMiddleware)
	router.Use(authMiddleware)
	// router.Use(middleware.CSRFMiddleware)

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

	// router.Handle("/metrics", promhttp.Handler())

	return router, nil
}

func RunServer() {
	router, err := NewRouter()
	if err != nil {
		log.Fatal("Failed to create router")
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.MainAppPort), router))
}
