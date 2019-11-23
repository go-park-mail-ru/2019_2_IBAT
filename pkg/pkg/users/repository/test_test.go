// dbconn := OpenSqlxViaPgxConnPool()
package repository

// type Record struct {
// 	Name string `db:"name"`
// 	Wage string `db:"wage"`
// 	Age  string `db:"age"`
// }

// func TestDBUserStorage_TestDynQueries_Correct(t *testing.T) {

// 	dbconn := OpenSqlxViaPgxConnPool()

// 	args := map[string]interface{}{
// 		"wage": "1000",
// 		`age`:  `25`,
// 	}
// 	// query := make([]string, 0)
// 	var query string
// 	count := 1
// 	for i, _ := range args {
// 		query += i + `=:` + i
// 		if count != len(args) {
// 			query += ` AND `
// 		}
// 		count++
// 	}

// 	fmt.Println("Query: " + query)
// 	// rows, err := dbconn.Queryx("SELECT * FROM users WHERE " + query, args)
// 	nmst, err := dbconn.PrepareNamed(`SELECT * FROM test WHERE ` + query + " ;")
// 	// rows.
// 	if err != nil {
// 		t.Errorf("unexpected err: %s", err)
// 		return
// 	}

// 	rows, err := nmst.Queryx(args)
// 	var records []Record

// 	for rows.Next() {
// 		var record Record

// 		_ = rows.StructScan(&record)

// 		records = append(records, record)
// 	}
// 	fmt.Println(records)

// 	// records
// 	if err != nil {
// 		t.Errorf("unexpected err: %s", err)
// 		return
// 	}
// }

// func OpenSqlxViaPgxConnPool() *sqlx.DB {
// 	connConfig := pgx.ConnConfig{
// 		Host:     "localhost",
// 		Database: "hh",
// 		User:     "postgres",
// 		Password: "newPassword",
// 	}
// 	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
// 		ConnConfig:     connConfig,
// 		AfterConnect:   nil,
// 		MaxConnections: 20,
// 		AcquireTimeout: 30 * time.Second,
// 	})
// 	if err != nil {
// 		log.Fatal("Failed to create connections pool")
// 	}

// 	// Apply any migrations...

// 	// Then set up sqlx and return the created DB reference
// 	// nativeDB, err := stdlib.OpenFromConnPool(connPool)
// 	nativeDB := stdlib.OpenDBFromPool(connPool)

// 	log.Println("OpenSqlxViaPgxConnPool: the connection was created")
// 	return sqlx.NewDb(nativeDB, "pgx")
// }
