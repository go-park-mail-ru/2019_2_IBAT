package server

import (
	"2019_2_IBAT/pkg/app/auth/session"
	. "2019_2_IBAT/pkg/app/chat/models"
	"2019_2_IBAT/pkg/app/chat/repository"
	"2019_2_IBAT/pkg/app/chat/service"

	"2019_2_IBAT/pkg/pkg/config"
	"2019_2_IBAT/pkg/pkg/db_connect"
	"2019_2_IBAT/pkg/pkg/middleware"
	"strconv"
	"sync"

	"net/http"

	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func RunServer() error {
	s := service.Service{
		MainChan: make(chan InChatMessage, 100),
		ConnectsPool: service.WsConnects{
			Connects: map[uuid.UUID]*service.ConnectsPerUser{},
			ConsMu:   &sync.Mutex{},
		},
		Storage: repository.DBStorage{
			DbConn: db_connect.OpenSqlxViaPgxConnPool(),
		},
	}

	router := mux.NewRouter()

	loger := middleware.NewLogger()

	authGrcpConn, err := grpc.Dial(
		"127.0.0.1:"+strconv.Itoa(config.AuthServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to authGrcpConn")
		return err
	}
	sessManager := session.NewServiceClient(authGrcpConn)
	authMiddleware := middleware.AuthMiddlewareGenerator(sessManager)

	router.Use(loger.AccessLogMiddleware)
	router.Use(middleware.CorsMiddleware)
	router.Use(authMiddleware)

	// router.Use(middleware.CSRFMiddleware)
	// router = router.PathPrefix("/api").Subrouter()

	router.HandleFunc("/api/chat/{companion_id}", s.HandleCreateChat).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/chat/history/{id}", s.HandlerGetChatHistory).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/chat/ws", s.HandleChat)
	router.HandleFunc("/api/chat/list", s.HandlerGetChats).Methods(http.MethodGet, http.MethodOptions)

	for i := 0; i < config.ChatWorkers; i++ {
		go s.ProcessMessage()
	}
	fmt.Println("starting main server at :" + strconv.Itoa(config.ChatAppPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.ChatAppPort), router))

	return nil
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	// if r.URL.Path != "/" {
	// 	http.Error(w, "Not found", http.StatusNotFound)
	// 	return
	// }
	// if r.Method != "GET" {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	http.ServeFile(w, r, "home.html")
}
