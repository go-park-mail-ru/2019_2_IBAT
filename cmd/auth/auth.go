package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	auth_rep "2019_2_IBAT/pkg/app/auth/repository"
	"2019_2_IBAT/pkg/app/auth/service"
	"2019_2_IBAT/pkg/app/auth/session"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	redisAddr := flag.String("redisServer", ":6379", "")

	// aS := auth_serv.AuthService{
	// 	Storage: auth_rep.NewSessionManager(auth_rep.RedNewPool(*redisAddr)),
	// }

	server := grpc.NewServer()

	session.RegisterServiceServer(server, service.AuthService{
		Storage: auth_rep.NewSessionManager(auth_rep.RedNewPool(*redisAddr)),
	})

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
