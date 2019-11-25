package db_connect

import (
	"log"
	"time"

	"github.com/jackc/pgx"

	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func OpenSqlxViaPgxConnPool() *sqlx.DB {
	connConfig := pgx.ConnConfig{
		Host:     "localhost",
		Database: "hh",
		User:     "postgres",
		Password: "newPassword",
	}
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AfterConnect:   nil,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to create connections pool")
	}

	nativeDB := stdlib.OpenDBFromPool(connPool)

	log.Println("OpenSqlxViaPgxConnPool: the connection was created")
	return sqlx.NewDb(nativeDB, "pgx")
}
