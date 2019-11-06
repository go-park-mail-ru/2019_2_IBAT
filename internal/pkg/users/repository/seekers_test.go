package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// // func init() {
// // 	/* load test data */
// // 	db, mock, err := sqlmock.New()
// // 	if err != nil {
// // 		t.Fatalf("cant create mock: %s", err)
// // 	}
// // 	defer db.Close()
// // }

// seekers := []Seeker{}

// 	rows, err := m.DbConn.Queryx("SELECT id, email, first_name, second_name,"+
// 		"path_to_image FROM persons WHERE role = $1;", SeekerStr)
// 	defer rows.Close()

// 	if err != nil {
// 		fmt.Println("GetSeeker: error while query seekers")
// 		return seekers, err
// 	}

// 	for rows.Next() {
// 		seek := Seeker{}
// 		_ = rows.StructScan(&seek)
// 		// if err != nil {
// 		// 	return seekers, err
// 		// }

// 		id_rows, err := m.DbConn.Query("SELECT v.id FROM resumes AS v WHERE v.own_id = $1;", seek.ID)
// 		if err != nil {
// 			fmt.Println("GetSeeker: error while query resumes")
// 			return seekers, err
// 		}
// 		defer id_rows.Close()

// 		resumes := make([]uuid.UUID, 0)

// 		for id_rows.Next() {
// 			var id uuid.UUID
// 			_ = id_rows.Scan(&id)
// 			// if err != nil {
// 			// 	return employers, err
// 			// }
// 			resumes = append(resumes, id)
// 		}

// 		seek.Resumes = resumes
// 		seekers = append(seekers, seek)
// 	}

// 	return seekers, nil
func TestDBUserStorage_GetSeekers_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	rows := sqlmock.
		NewRows([]string{"id", "email", "first_name", "second_name", "path_to_image"})
	rows_resumes_id1 := sqlmock.NewRows([]string{"id"})
	rows_resumes_id2 := sqlmock.NewRows([]string{"id"})

	expect := []Seeker{
		{
			ID:         uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"),
			Email:      "some@mail.ru",
			FirstName:  "Victor",
			SecondName: "Timofeev",
			PathToImg:  "",
			Resumes:    []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
		},
		{
			ID:         uuid.MustParse("f14c6111-3430-413b-ab4e-e31c8642ad8a"),
			Email:      "some@mail.ru",
			FirstName:  "Victor",
			SecondName: "Timofeev",
			PathToImg:  "",
			Resumes:    []uuid.UUID{uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d")},
		},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID.String(), item.Email, item.FirstName, item.SecondName,
			item.PathToImg,
		)
	}

	mock.
		ExpectQuery("SELECT id, email, first_name, second_name," +
			"path_to_image FROM persons WHERE").
		WithArgs(SeekerStr).
		WillReturnRows(rows)

	rows_resumes_id1 = rows_resumes_id1.AddRow(uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d").String())
	mock.
		ExpectQuery("SELECT r.id FROM resumes AS r WHERE").
		WithArgs(expect[0].ID).
		WillReturnRows(rows_resumes_id1)

	rows_resumes_id2 = rows_resumes_id2.AddRow(uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d").String())
	mock.
		ExpectQuery("SELECT r.id FROM resumes AS r WHERE").
		WithArgs(expect[1].ID).
		WillReturnRows(rows_resumes_id2)

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	seekers, err := repo.GetSeekers()

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(seekers, expect) {
		t.Errorf("results not match,\n want\n%v,\n have\n %v\n", expect, seekers)
		return
	}
}

func TestDBUserStorage_GetSeekers_Fail(t *testing.T) { //ADD SECOND SELECT TEST CASE
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	mock.
		ExpectQuery("SELECT id, email, first_name, second_name," +
			"path_to_image FROM persons WHERE").
		WithArgs().
		WillReturnError(errors.New("GetSeeker: error while query seekers"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	seekers, err := repo.GetSeekers()
	fmt.Println(seekers)

	if err == nil {
		fmt.Println(err)
		t.Errorf("Expected err")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

}

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
// 		ExpectQuery("SELECT id, email, first_name, second_name," +
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

func TestDBUserStorage_GetSeeker_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")
	rows := sqlmock.
		NewRows([]string{"id", "email", "first_name", "second_name", "path_to_image"})

	mock.
		ExpectQuery("SELECT id, email, first_name, second_name," +
			"path_to_image FROM persons WHERE").
		WithArgs(id).
		WillReturnRows(rows)

	// mock.
	// 	ExpectQuery("SELECT r.id FROM resumes AS r WHERE").
	// 	WithArgs(id).
	// 	WillReturnError(errors.New("GetSeeker: Invalid id"))
	// mock.
	// 	ExpectQuery("SELECT v.id, v.own_id, c.company_name, v.experience," +
	// 		"v.profession, v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about" +
	// 		" FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id WHERE").
	// 	WithArgs(id).
	// 	WillReturnError(errors.New("GetVacancy: error while querying"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	seeker, err := repo.GetSeeker(id)
	fmt.Println(seeker)

	if err == nil {
		fmt.Println(err)
		t.Errorf("Expected err")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

// func TestDBUserStorage_CreateVacancy_Correct(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	vacancy := Vacancy{
// 		ID:           uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"),
// 		OwnerID:      uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
// 		CompanyName:  "MC",
// 		Experience:   "7 years",
// 		Profession:   "cleaner",
// 		Position:     "mid",
// 		Tasks:        "cleaning rooms",
// 		Requirements: "work for 24 hours per week",
// 		WageFrom:     "100 500.00 руб",
// 		WageTo:       "120 500.00 руб",
// 		Conditions:   "Nice geolocation",
// 		About:        "Hello employer",
// 	}

// 	mock.
// 		ExpectExec(`INSERT INTO vacancies`).
// 		WithArgs(
// 			vacancy.ID, vacancy.OwnerID, vacancy.Experience, vacancy.Profession, vacancy.Position, vacancy.Tasks,
// 			vacancy.Requirements, vacancy.Conditions, vacancy.WageFrom, vacancy.WageTo, vacancy.About,
// 		).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	ok := repo.CreateVacancy(vacancy)

// 	if !ok {
// 		t.Error("Failed to create vacancy\n")
// 		return
// 	}
// }

// func TestDBUserStorage_CreateVacancy_False(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	vacancy := Vacancy{
// 		ID:           uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"),
// 		OwnerID:      uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
// 		CompanyName:  "MC",
// 		Experience:   "7 years",
// 		Profession:   "cleaner",
// 		Position:     "mid",
// 		Tasks:        "cleaning rooms",
// 		Requirements: "work for 24 hours per week",
// 		WageFrom:     "100 500.00 руб",
// 		WageTo:       "120 500.00 руб",
// 		Conditions:   "Nice geolocation",
// 		About:        "Hello employer",
// 	}

// 	mock.
// 		ExpectExec(`INSERT INTO vacancies`).
// 		WithArgs(
// 			vacancy.ID, vacancy.OwnerID, vacancy.Experience, vacancy.Profession, vacancy.Position, vacancy.Tasks,
// 			vacancy.Requirements, vacancy.Conditions, vacancy.WageFrom, vacancy.WageTo, vacancy.About,
// 		).
// 		WillReturnError(fmt.Errorf("bad query"))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	ok := repo.CreateVacancy(vacancy)

// 	if ok {
// 		t.Errorf("expected false, got true")
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }

// func TestDBUserStorage_DeleteVacancy_Correct(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a")

// 	mock.
// 		ExpectExec(`DELETE FROM vacancies`).
// 		WithArgs(id).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	err = repo.DeleteVacancy(id)
// 	fmt.Println()

// 	if err != nil {
// 		t.Errorf("unexpected err: %s", err)
// 		return
// 	}

// 	mock.
// 		ExpectQuery("SELECT v.id, v.own_id, c.company_name, v.experience," +
// 			"v.profession, v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about" +
// 			" FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id WHERE").
// 		WithArgs(id).
// 		WillReturnError(errors.New("GetVacancy: error while querying"))

// 	_, err = repo.GetVacancy(id)
// 	fmt.Println()

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

// func TestDBUserStorage_DeleteVacancy_False(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a")

// 	mock.
// 		ExpectExec(`DELETE FROM vacancies`).
// 		WithArgs(id).
// 		WillReturnError(errors.Errorf("error"))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	err = repo.DeleteVacancy(id)

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

// func TestDBUserStorage_PutVacancy_Correct(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	vacancy := Vacancy{
// 		ID:           uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"),
// 		OwnerID:      uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
// 		CompanyName:  "MC",
// 		Experience:   "7 years",
// 		Profession:   "cleaner",
// 		Position:     "mid",
// 		Tasks:        "cleaning rooms",
// 		Requirements: "work for 24 hours per week",
// 		WageFrom:     "100 500.00 руб",
// 		WageTo:       "101 500.00 руб",
// 		Conditions:   "Nice geolocation",
// 		About:        "Hello employer",
// 	}

// 	mock.
// 		ExpectExec(`UPDATE vacancies SET`).
// 		WithArgs(
// 			vacancy.Experience, vacancy.Profession, vacancy.Position, vacancy.Tasks, vacancy.Requirements,
// 			vacancy.Conditions, vacancy.WageFrom, vacancy.WageTo, vacancy.About, vacancy.ID, vacancy.OwnerID,
// 		).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	ok := repo.PutVacancy(vacancy, vacancy.OwnerID, vacancy.ID)

// 	if !ok {
// 		t.Error("Failed to put vacancy\n")
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }

// func TestDBUserStorage_PutVacancy_False(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	vacancy := Vacancy{
// 		ID:           uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"), //invalid id
// 		OwnerID:      uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
// 		CompanyName:  "MC",
// 		Experience:   "7 years",
// 		Profession:   "cleaner",
// 		Position:     "mid",
// 		Tasks:        "cleaning rooms",
// 		Requirements: "work for 24 hours per week",
// 		WageFrom:     "100 500.00 руб",
// 		WageTo:       "101 500.00 руб",
// 		Conditions:   "Nice geolocation",
// 		// About:        "Hello employer",
// 	}

// 	mock.
// 		ExpectExec(`UPDATE vacancies SET`).
// 		WithArgs(
// 			vacancy.Experience, vacancy.Profession, vacancy.Position, vacancy.Tasks, vacancy.Requirements,
// 			vacancy.Conditions, vacancy.WageFrom, vacancy.WageTo, vacancy.About, vacancy.ID, vacancy.OwnerID,
// 		).
// 		WillReturnError(fmt.Errorf("bad query"))

// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}

// 	ok := repo.PutVacancy(vacancy, vacancy.OwnerID, vacancy.ID)

// 	if ok {
// 		t.Errorf("expected false, got true")
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }
