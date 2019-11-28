package repository

import (
	. "2019_2_IBAT/pkg/pkg/models"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDBUserStorage_GetEmployers_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	rows := sqlmock.
		NewRows([]string{"id", "email", "company_name", "first_name", "second_name",
			"site", "empl_num", "phone_number", "extra_phone_number", "spheres_of_work", "path_to_image",
			"region", "description"})
	rows_vacancies_id1 := sqlmock.NewRows([]string{"id"})
	rows_vacancies_id2 := sqlmock.NewRows([]string{"id"})

	expected := []Employer{
		{
			ID:               uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
			CompanyName:      "MCDonalds",
			Site:             "petushki.com",
			Email:            "petushki@mail.com",
			FirstName:        "Vova",
			SecondName:       "Zyablikov",
			PhoneNumber:      "12345678911",
			ExtraPhoneNumber: "12345678910",
			Region:           "Petushki",
			EmplNum:          "322",
			SpheresOfWork:    "string",
			Vacancies:        []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
		},
		{
			ID:               uuid.MustParse("1ba7b811-9dad-11d1-80b1-00c04fd430c8"),
			CompanyName:      "IDs",
			Site:             "IDS.com",
			Email:            "ids@mail.com",
			FirstName:        "Kostya",
			SecondName:       "Zyablikov",
			PhoneNumber:      "12345678911",
			ExtraPhoneNumber: "12345678910",
			Region:           "Moscow",
			EmplNum:          "322",
			SpheresOfWork:    "string",
			Vacancies:        []uuid.UUID{uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d")},
		},
	}

	for _, item := range expected {
		rows = rows.AddRow(item.ID.String(), item.Email, item.CompanyName, item.FirstName,
			item.SecondName, item.Site, item.EmplNum, item.PhoneNumber, item.ExtraPhoneNumber,
			item.SpheresOfWork, item.PathToImg, item.Region, item.Description,
		)
	}
	mock.
		ExpectQuery("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site," +
			"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image, c.region, " +
			" c.description FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE").
		WithArgs(EmployerStr).
		WillReturnRows(rows)

	rows_vacancies_id1 = rows_vacancies_id1.AddRow(uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d").String())
	mock.
		ExpectQuery("SELECT v.id FROM vacancies AS v WHERE ").
		WithArgs(expected[0].ID).
		WillReturnRows(rows_vacancies_id1)

	rows_vacancies_id2 = rows_vacancies_id2.AddRow(uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d").String())
	mock.
		ExpectQuery("SELECT v.id FROM vacancies AS v WHERE ").
		WithArgs(expected[1].ID).
		WillReturnRows(rows_vacancies_id2)

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	dummyMap := map[string]interface{}{}
	employers, err := repo.GetEmployers(dummyMap)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.Equal(t, employers, expected, "The two values should be the same.")
}

func TestDBUserStorage_GetEmployers_Fail(t *testing.T) { //ADD SECOND SELECT TEST CASE
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	mock.
		ExpectQuery("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site," +
			"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image, c.region, " +
			" c.description FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE").
		WithArgs(EmployerStr).
		WillReturnError(errors.New("GetEmployers: error while query employers"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	dummyMap := map[string]interface{}{}
	employers, err := repo.GetEmployers(dummyMap)
	fmt.Println(employers)

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

func TestDBUserStorage_GetEmployer_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	rows := sqlmock.
		NewRows([]string{"id", "email", "company_name", "first_name", "second_name",
			"site", "empl_num", "phone_number", "extra_phone_number", "spheres_of_work", "path_to_image",
			"region", "description"})
	rows_vacancies_id1 := sqlmock.NewRows([]string{"id"})

	expect := Employer{
		ID:               uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
		CompanyName:      "MCDonalds",
		Site:             "petushki.com",
		Email:            "petushki@mail.com",
		FirstName:        "Vova",
		SecondName:       "Zyablikov",
		PhoneNumber:      "12345678911",
		ExtraPhoneNumber: "12345678910",
		Region:           "Petushki",
		EmplNum:          "322",
		SpheresOfWork:    "string",
		Vacancies:        []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
	}

	rows = rows.AddRow(expect.ID.String(), expect.Email, expect.CompanyName, expect.FirstName,
		expect.SecondName, expect.Site, expect.EmplNum, expect.PhoneNumber, expect.ExtraPhoneNumber,
		expect.SpheresOfWork, expect.PathToImg, expect.Region, expect.Description,
	)

	mock.
		ExpectQuery("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site," +
			"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image, c.region, " +
			" c.description FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE").
		WithArgs(expect.ID).
		WillReturnRows(rows)

	rows_vacancies_id1 = rows_vacancies_id1.AddRow(uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d").String())

	mock.
		ExpectQuery("SELECT v.id FROM vacancies AS v WHERE").
		WithArgs(expect.ID).
		WillReturnRows(rows_vacancies_id1)

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	item, err := repo.GetEmployer(expect.ID)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.Equal(t, item, expect, "The two values should be the same.")
}

func TestDBUserStorage_GetEmployer_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")

	mock.
		ExpectQuery("SELECT id, email, first_name, second_name, " +
			"path_to_image FROM persons WHERE ").
		WithArgs(id).
		WillReturnError(errors.New("sql: no rows in result set"))

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

func TestDBUserStorage_GetEmployer_Fail2(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	rows := sqlmock.
		NewRows([]string{"id", "email", "company_name", "first_name", "second_name",
			"site", "empl_num", "phone_number", "extra_phone_number", "spheres_of_work", "path_to_image",
			"region", "description"})
	// rows_vacancies_id1 := sqlmock.NewRows([]string{"id"})

	expect := Employer{
		ID:               uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
		CompanyName:      "MCDonalds",
		Site:             "petushki.com",
		Email:            "petushki@mail.com",
		FirstName:        "Vova",
		SecondName:       "Zyablikov",
		PhoneNumber:      "12345678911",
		ExtraPhoneNumber: "12345678910",
		Region:           "Petushki",
		EmplNum:          "322",
		SpheresOfWork:    "string",
		Vacancies:        []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
	}

	rows = rows.AddRow(expect.ID.String(), expect.Email, expect.CompanyName, expect.FirstName,
		expect.SecondName, expect.Site, expect.EmplNum, expect.PhoneNumber, expect.ExtraPhoneNumber,
		expect.SpheresOfWork, expect.PathToImg, expect.Region, expect.Description,
	)

	mock.
		ExpectQuery("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site," +
			"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image, c.region, " +
			" c.description FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE").
		WithArgs(expect.ID).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT v.id FROM vacancies AS v WHERE").
		WithArgs(expect.ID).
		WillReturnError(errors.New("GetEmployer: Invalid id"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	seeker, err := repo.GetEmployer(expect.ID)
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

func TestDBUserStorage_CreateEmployer_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	employer := Employer{
		ID:               uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
		CompanyName:      "MCDonalds",
		Site:             "petushki.com",
		Email:            "petushki@mail.com",
		FirstName:        "Vova",
		SecondName:       "Zyablikov",
		PhoneNumber:      "12345678911",
		ExtraPhoneNumber: "12345678910",
		Region:           "Petushki",
		EmplNum:          "322",
		SpheresOfWork:    "string",
		Vacancies:        []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`INSERT INTO persons`).
		WithArgs(
			employer.ID, employer.Email, employer.FirstName,
			employer.SecondName, employer.Password, EmployerStr, employer.PathToImg,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec(`INSERT INTO companies`).
		WithArgs(
			employer.ID, employer.CompanyName, employer.Site, employer.PhoneNumber,
			employer.ExtraPhoneNumber, employer.Region, employer.EmplNum,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	ok := repo.CreateEmployer(employer)

	if !ok {
		t.Error("Failed to create vacancy\n")
		return
	}
}

func TestDBUserStorage_CreateEmployer_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	employer := Employer{
		ID:               uuid.MustParse("6ba7b811-9dad-11d1-80b1-00c04fd430c8"),
		CompanyName:      "MCDonalds",
		Site:             "petushki.com",
		Email:            "petushki@mail.com",
		FirstName:        "Vova",
		SecondName:       "Zyablikov",
		PhoneNumber:      "12345678911",
		ExtraPhoneNumber: "12345678910",
		Region:           "Petushki",
		EmplNum:          "322",
		SpheresOfWork:    "string",
		Vacancies:        []uuid.UUID{uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d")},
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`INSERT INTO persons`).
		WithArgs(
			employer.ID, employer.Email, employer.FirstName,
			employer.SecondName, employer.Password, EmployerStr, employer.PathToImg,
		).
		WillReturnError(errors.New("bad query"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	ok := repo.CreateEmployer(employer)

	if ok {
		t.Errorf("expected false, got true")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStorage_PutEmployer_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	employer := EmployerReg{
		CompanyName:      "MCDonalds",
		Site:             "petushki.com",
		Email:            "petushki@mail.com",
		FirstName:        "Vova",
		SecondName:       "Zyablikov",
		PhoneNumber:      "12345678911",
		ExtraPhoneNumber: "12345678910",
		Region:           "Petushki",
		EmplNum:          "322",
	}

	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE persons SET`).
		WithArgs(
			employer.Email, employer.FirstName,
			employer.SecondName, employer.Password, id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("UPDATE companies SET").
		WithArgs(
			employer.CompanyName, employer.Site, employer.PhoneNumber,
			employer.ExtraPhoneNumber, employer.Region, employer.EmplNum, id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	ok := repo.PutEmployer(employer, id)

	if !ok {
		t.Error("Failed to put employer\n")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStorage_PutEmployer_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	employer := EmployerReg{
		CompanyName:      "MCDonalds",
		Site:             "petushki.com",
		Email:            "petushki@mail.com",
		FirstName:        "Vova",
		SecondName:       "Zyablikov",
		PhoneNumber:      "12345678911",
		ExtraPhoneNumber: "12345678910",
		Region:           "Petushki",
		EmplNum:          "322",
	}

	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE persons SET`).
		WithArgs(
			employer.Email, employer.FirstName,
			employer.SecondName, employer.Password, id,
		).
		WillReturnError(fmt.Errorf("bad query"))
	mock.ExpectRollback()
	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	ok := repo.PutEmployer(employer, id)

	if ok {
		t.Errorf("expected false, got true")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
