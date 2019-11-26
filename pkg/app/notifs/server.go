package notifs

import (
	"2019_2_IBAT/pkg/app/auth/session"
	"2019_2_IBAT/pkg/app/notifs/notifsproto"
	"2019_2_IBAT/pkg/app/notifs/service"
	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	"2019_2_IBAT/pkg/pkg/config"
	"2019_2_IBAT/pkg/pkg/middleware"
	"strconv"

	"net/http"

	"fmt"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func RunServer() error {
	authGrcpConn, err := grpc.Dial(
		"127.0.0.1:"+strconv.Itoa(config.AuthServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to authGrcpConn")
		return err
	}

	// defer grcpConn.Close()

	sessManager := session.NewServiceClient(authGrcpConn)

	recomsGrcpConn, err := grpc.Dial(
		"127.0.0.1:"+strconv.Itoa(config.RecommendsServicePort),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("cant connect to recomsGrcpConn")
		return err
	}

	recommsManager := recomsproto.NewServiceClient(recomsGrcpConn)

	notifService := service.Service{
		ConnectsPool: service.WsConnects{
			Connects: map[uuid.UUID]*service.ConnectsPerUser{},
			ConsMu:   &sync.Mutex{},
		},
		NotifChan:    make(chan service.NotifStruct, 64),
		AuthService:  sessManager,
		RecomService: recommsManager,
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.NotifsServicePort))
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()
	notifsproto.RegisterServiceServer(server, notifService)

	fmt.Println("starting rpc server at :" + strconv.Itoa(config.NotifsServicePort))
	go server.Serve(lis)

	router := mux.NewRouter()

	loger := middleware.NewLogger()

	authMiddleware := middleware.AuthMiddlewareGenerator(sessManager)

	router.Use(loger.AccessLogMiddleware)
	// router.Use(middleware.CorsMiddleware)
	router.Use(authMiddleware)
	// router.Use(middleware.CSRFMiddleware)

	// router = router.PathPrefix("/api/").Subrouter()

	router.Use(loger.AccessLogMiddleware)
	router.Use(middleware.CorsMiddleware)
	router.Use(authMiddleware)

	// fmt.Println(notifService.HandleNotifications)
	router.HandleFunc("/api/notifications", notifService.HandleNotifications)

	go notifService.Notifications()

	fmt.Println("starting main server at :" + strconv.Itoa(config.NotifsAppPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.NotifsAppPort), router))

	return nil
}
