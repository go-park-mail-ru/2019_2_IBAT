package main

import (
	"hh_workspace/2019_2_IBAT/internal/pkg/server"
)

func main() {
	server, _ := server.NewServer()
	server.Run()
}
