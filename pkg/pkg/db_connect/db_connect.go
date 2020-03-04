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
		Host:     config.DBHostname,
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

	fmt.Println("OpenSqlxViaPgxConnPool: test db")
	fmt.Println("OpenSqlxViaPgxConnPool: test base #2")
	fmt.Printf("OpenSqlxViaPgxConnPool: Hostname %s\n", config.DBHostname)
	if err != nil {
		fmt.Printf("OpenSqlxViaPgxConnPool: Failed to create connections pool: error - %s\n", err.Error())
		log.Fatalf("OpenSqlxViaPgxConnPool: Failed to create connections pool: error - %s\n", err.Error())
	}

	nativeDB := stdlib.OpenDBFromPool(connPool)

	log.Println("OpenSqlxViaPgxConnPool: the connection was created")
	return sqlx.NewDb(nativeDB, "pgx")
}
