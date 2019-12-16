package db_connect

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"2019_2_IBAT/pkg/pkg/config"
)

func OpenSqlxViaPgxConnPool() *sqlx.DB {
	connConfig := pgx.ConnConfig{
		Host:     config.Hostname,
		Database: config.Database,
		User:     config.User,
		Password: config.Password,
	}
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AfterConnect:   nil,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		fmt.Printf("FFFFFFFFFailed to create connections pool: error - %s", err.Error())
		log.Fatalf("FFFFFFFFFailed to create connections pool: error - %s", err.Error())
	}

	nativeDB := stdlib.OpenDBFromPool(connPool)

	log.Println("OpenSqlxViaPgxConnPool: the connection was created")
	return sqlx.NewDb(nativeDB, "pgx")
}
