package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	auth_rep "2019_2_IBAT/pkg/app/auth/repository"
	"2019_2_IBAT/pkg/app/auth/service"
	"2019_2_IBAT/pkg/app/auth/session"
	"2019_2_IBAT/pkg/pkg/config"
)

func main() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.AuthServicePort))
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	redisAddr := flag.String("redisServer", ":"+strconv.Itoa(config.ReddisPort), "")

	// aS := auth_serv.AuthService{
	// 	Storage: auth_rep.NewSessionManager(auth_rep.RedNewPool(*redisAddr)),
	// }

	server := grpc.NewServer()

	session.RegisterServiceServer(server, service.AuthService{
		Storage: auth_rep.NewSessionManager(auth_rep.RedNewPool(*redisAddr)),
	})

	fmt.Println("starting server at " + strconv.Itoa(config.AuthServicePort))
	server.Serve(lis)
}
