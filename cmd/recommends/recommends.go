package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	"2019_2_IBAT/pkg/app/recommends/repository"
	"2019_2_IBAT/pkg/app/recommends/service"
	"2019_2_IBAT/pkg/pkg/config"
	"2019_2_IBAT/pkg/pkg/db_connect"
)

func main() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.RecommendsServicePort))
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	recomsproto.RegisterServiceServer(server,
		service.Service{
			Storage: repository.DBRecommendsStorage{
				DbConn: db_connect.OpenSqlxViaPgxConnPool(),
			},
		},
	)

	fmt.Println("starting server at :" + strconv.Itoa(config.RecommendsServicePort))
	server.Serve(lis)
}
