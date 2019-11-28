package repository

import (
	. "2019_2_IBAT/pkg/pkg/models"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDBUserStorage_DeleteUser_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	userID := uuid.New()
	mock.
		ExpectExec(`DELETE FROM persons WHERE `).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	err = repo.DeleteUser(userID)

	if err != nil {
		t.Error("Failed to delete user\n")
		return
	}
}

func TestDBUserStorage_DeleteUser_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	userID := uuid.New()
	mock.
		ExpectExec(`DELETE FROM persons WHERE `).
		WithArgs(userID).
		WillReturnError(fmt.Errorf("error"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	err = repo.DeleteUser(userID)

	if err == nil {
		t.Error("Expected error\n")
		return
	}
}

// func (m *DBUserStorage) CheckUser(email string, password string) (uuid.UUID, string, bool) {
// 	row := m.DbConn.QueryRow("SELECT id, role, password_hash FROM persons "+
// 		"WHERE email = $1", email)

// 	resId := uuid.UUID{}
// 	var class string
// 	var password_hash []byte
// 	err := row.Scan(&resId, &class, &password_hash)

// 	// if !passwords.CheckPass(password_hash, password) || err != nil {
// 	// 	return resId, class, false
// 	// }
// 	if password != string(password_hash) || err != nil {
// 		if err != nil {
// 			fmt.Printf("CheckUser: %s\n", err)
// 		}
// 		return resId, class, false
// 	}
// 	return resId, class, true
// }

func TestDBUserStorage_CheckUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	// id := uuid.New()
	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a")
	email := "ex@mail.com"
	password := "11111"
	role := SeekerStr
	ok := true

	rows := sqlmock.
		NewRows([]string{"id", "role", "password_hash"})
	rows = rows.AddRow(id.String(), role, password)

	mock.
		ExpectQuery(`SELECT id, role, password_hash FROM persons WHERE`).
		WithArgs(email).
		WillReturnRows(rows)

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	gotId, gotRole, gotOk := repo.CheckUser(email, password)

	require.Equal(t, id, gotId, "The two values should be the same.")
	require.Equal(t, role, gotRole, "The two values should be the same.")
	require.Equal(t, ok, gotOk, "The two values should be the same.")

	wrongPassword := "1111111"
	ok = false

	rows = rows.AddRow(id.String(), role, wrongPassword)

	mock.
		ExpectQuery(`SELECT id, role, password_hash FROM persons WHERE`).
		WithArgs(email).
		WillReturnRows(rows)

	gotId, gotRole, gotOk = repo.CheckUser(email, password)

	require.Equal(t, id, gotId, "The two values should be the same.")
	require.Equal(t, role, gotRole, "The two values should be the same.")
	require.Equal(t, ok, gotOk, "The two values should be the same.")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// func TestDBUserStorage_GetSeekers_Correct(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}

// 	rows := sqlmock.
// 		NewRows([]string{"id", "email", "first_name", "second_name", "path_to_image"})
// 	rows_resumes_id1 := sqlmock.NewRows([]string{"id"})
// 	rows_resumes_id2 := sqlmock.NewRows([]string{"id"})

// 	expect := []Seeker{
// 		{
// 			ID:         uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"),
// 			Email:      "some@mail.ru",
// 			FirstName:  "Victor",
// 			SecondName: "Timofeev",
// 			PathToImg:  "",
// 			Resumes:    []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
// 		},
// 		{
// 			ID:         uuid.MustParse("f14c6111-3430-413b-ab4e-e31c8642ad8a"),
// 			Email:      "some@mail.ru",
// 			FirstName:  "Victor",
// 			SecondName: "Timofeev",
// 			PathToImg:  "",
// 			Resumes:    []uuid.UUID{uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d")},
// 		},
// 	}

// 	for _, item := range expect {
// 		rows = rows.AddRow(item.ID.String(), item.Email, item.FirstName, item.SecondName,
// 			item.PathToImg,
// 		)
// 	}

// 	mock.
// 		ExpectQuery("SELECT id, email, first_name, second_name," +
// 			"path_to_image FROM persons WHERE").
// 		WithArgs(SeekerStr).
// 		WillReturnRows(rows)

// 	rows_resumes_id1 = rows_resumes_id1.AddRow(uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d").String())
// 	mock.
// 		ExpectQuery("SELECT r.id FROM resumes AS r WHERE").
// 		WithArgs(expect[0].ID).
// 		WillReturnRows(rows_resumes_id1)

// 	rows_resumes_id2 = rows_resumes_id2.AddRow(uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d").String())
// 	mock.
// 		ExpectQuery("SELECT r.id FROM resumes AS r WHERE").
// 		WithArgs(expect[1].ID).
// 		WillReturnRows(rows_resumes_id2)

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	seekers, err := repo.GetSeekers()

// 	if err != nil {
// 		t.Errorf("unexpected err: %s", err)
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 		return
// 	}
// 	if !reflect.DeepEqual(seekers, expect) {
// 		t.Errorf("results not match,\n want\n%v,\n have\n %v\n", expect, seekers)
// 		return
// 	}
// }

// func TestDBUserStorage_GetSeekers_Fail(t *testing.T) { //ADD SECOND SELECT TEST CASE
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}

// 	mock.
// 		ExpectQuery("SELECT id, email, first_name, second_name," +
// 			"path_to_image FROM persons WHERE").
// 		WithArgs().
// 		WillReturnError(errors.New("GetSeeker: error while query seekers"))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	seekers, err := repo.GetSeekers()
// 	fmt.Println(seekers)

// 	if err == nil {
// 		fmt.Println(err)
// 		t.Errorf("Expected err")
// 		return
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 		return
// 	}

// }

// func TestDBUserStorage_GetSeeker_Correct(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	rows := sqlmock.
// 		NewRows([]string{"id", "email", "first_name", "second_name", "path_to_image"})
// 	rows_resumes_id1 := sqlmock.NewRows([]string{"id"})

// 	expect := Seeker{
// 		ID:         uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"),
// 		Email:      "some@mail.ru",
// 		FirstName:  "Victor",
// 		SecondName: "Timofeev",
// 		PathToImg:  "",
// 		Resumes:    []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
// 	}

// 	rows = rows.AddRow(expect.ID.String(), expect.Email, expect.FirstName,
// 		expect.SecondName, expect.PathToImg,
// 	)

// 	mock.
// 		ExpectQuery("SELECT id, email, first_name, second_name, " +
// 			"path_to_image FROM persons WHERE").
// 		WithArgs(expect.ID).
// 		WillReturnRows(rows)

// 	rows_resumes_id1 = rows_resumes_id1.AddRow(uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d").String())

// 	mock.
// 		ExpectQuery("SELECT r.id FROM resumes AS r WHERE").
// 		WithArgs(expect.ID).
// 		WillReturnRows(rows_resumes_id1)

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	item, err := repo.GetSeeker(expect.ID)

// 	if err != nil {
// 		t.Errorf("unexpected err: %s", err)
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 		return
// 	}
// 	if !reflect.DeepEqual(item, expect) {
// 		t.Errorf("results not match,\n want\n%v,\n have\n %v\n", expect, item)
// 		return
// 	}
// }

// func TestDBUserStorage_GetSeeker_Fail(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")

// 	mock.
// 		ExpectQuery("SELECT id, email, first_name, second_name, " +
// 			"path_to_image FROM persons WHERE ").
// 		WithArgs(id).
// 		WillReturnError(errors.New("sql: no rows in result set"))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	seeker, err := repo.GetSeeker(id)
// 	fmt.Println(seeker)

// 	if err == nil {
// 		fmt.Println(err)
// 		t.Errorf("Expected err")
// 		return
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 		return
// 	}
// }

// func TestDBUserStorage_GetSeeker_Fail2(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	rows := sqlmock.
// 		NewRows([]string{"id", "email", "first_name", "second_name", "path_to_image"})

// 	expect := Seeker{
// 		ID:         uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba"),
// 		Email:      "some@mail.ru",
// 		FirstName:  "Victor",
// 		SecondName: "Timofeev",
// 		PathToImg:  "",
// 		Resumes:    []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
// 	}

// 	rows = rows.AddRow(expect.ID.String(), expect.Email, expect.FirstName,
// 		expect.SecondName, expect.PathToImg,
// 	)

// 	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")

// 	mock.
// 		ExpectQuery("SELECT id, email, first_name, second_name, " +
// 			"path_to_image FROM persons WHERE ").
// 		WithArgs(id).
// 		WillReturnRows(rows)

// 	mock.
// 		ExpectQuery("SELECT r.id FROM resumes AS r WHERE").
// 		WithArgs(id).
// 		WillReturnError(errors.New("GetSeeker: Invalid id"))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	seeker, err := repo.GetSeeker(id)
// 	fmt.Println(seeker)

// 	if err == nil {
// 		fmt.Println(err)
// 		t.Errorf("Expected err")
// 		return
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 		return
// 	}
// }

// func TestDBUserStorage_CreateSeeker_False(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	seeker := Seeker{
// 		ID:         uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba"),
// 		Email:      "some@mail.ru",
// 		FirstName:  "Victor",
// 		SecondName: "Timofeev",
// 		PathToImg:  "",
// 		Resumes:    []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
// 	}

// 	mock.
// 		ExpectExec(`INSERT INTO persons`).
// 		WithArgs(
// 			seeker.ID, seeker.Email, seeker.FirstName,
// 			seeker.SecondName, seeker.Password, SeekerStr, seeker.PathToImg,
// 		).
// 		WillReturnError(fmt.Errorf("bad query"))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	ok := repo.CreateSeeker(seeker)

// 	if ok {
// 		t.Errorf("expected false, got true")
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }

// func TestDBUserStorage_PutSeeker_Correct(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	seeker := SeekerReg{
// 		Email:      "some@mail.ru",
// 		FirstName:  "Victor",
// 		SecondName: "Timofeev",
// 		Password:   "sdfsdf",
// 	}

// 	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")

// 	mock.
// 		ExpectExec(`UPDATE persons SET`).
// 		WithArgs(
// 			seeker.Email, seeker.FirstName,
// 			seeker.SecondName, seeker.Password, id,
// 		).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	ok := repo.PutSeeker(seeker, id)

// 	if !ok {
// 		t.Error("Failed to put seeker\n")
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }

// func TestDBUserStorage_PutSeeker_False(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	seeker := SeekerReg{
// 		Email:      "some@mail.ru",
// 		FirstName:  "Victor",
// 		SecondName: "Timofeev",
// 		Password:   "sdfsdf",
// 	}

// 	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")

// 	mock.
// 		ExpectExec(`UPDATE persons SET`).
// 		WithArgs(
// 			seeker.Email, seeker.FirstName,
// 			seeker.SecondName, seeker.Password, id,
// 		).
// 		WillReturnError(fmt.Errorf("bad query"))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	ok := repo.PutSeeker(seeker, id)

// 	if ok {
// 		t.Errorf("expected false, got true")
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }
