package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	"2019_2_IBAT/pkg/app/recommends/repository"
	"2019_2_IBAT/pkg/app/recommends/service"
	"2019_2_IBAT/pkg/pkg/db_connect"
)

func main() {
	lis, err := net.Listen("tcp", ":8082")
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

	fmt.Println("starting server at :8082")
	server.Serve(lis)
}

// изменено:      .vscode/launch.json
// изменено:      cmd/auth/auth.go
// изменено:      cmd/main/main
// изменено:      cmd/main/main.go
// изменено:      cmd/main/testlogfile
