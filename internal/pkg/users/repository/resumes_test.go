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

// func init() {
// 	/* load test data */
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer db.Close()
// }

func TestDBUserStorage_GetResumes_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	rows := sqlmock.
		NewRows([]string{"id", "own_id", "email", "city", "phone_number",
			"first_name", "second_name", "birth_date", "sex", "citizenship",
			"profession", "position", "experience", "education", "wage", "about",
		})
	expect := []Resume{
		{
			ID:      uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"),
			OwnerID: uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
			// Email:       "",
			City:        "Moscow",
			PhoneNumber: "12345678910",
			FirstName:   "Vova",
			SecondName:  "Zyablikov",
			BirthDate:   "1999-01-08",
			Sex:         "male",
			Citizenship: "Russia",
			Profession:  "programmer",
			Position:    "middle",
			Experience:  "7 years",
			Education:   "MSU",
			Wage:        "100 500.00 руб",
			About:       "Hello employer",
		},
		{
			ID:          uuid.MustParse("f14c6104-3431-413b-ab4e-e31c8642ad8a"),
			OwnerID:     uuid.MustParse("92b77777-bac7-4597-ab71-7b5fbe53052d"),
			Email:       "email@mail.ru",
			City:        "Moscow",
			PhoneNumber: "12345678910",
			FirstName:   "Petya",
			SecondName:  "Zyablikov",
			BirthDate:   "1986-01-08",
			Sex:         "male",
			Citizenship: "Russia",
			Profession:  "programmer",
			Position:    "middle",
			Experience:  "7 years",
			Education:   "MSU",
			Wage:        "100 500.00 руб",
			About:       "Hello employer",
		},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID.String(), item.OwnerID.String(), item.Email, item.City, item.PhoneNumber, item.FirstName,
			item.SecondName, item.BirthDate, item.Sex, item.Citizenship, item.Profession,
			item.Position, item.Experience, item.Education, item.Wage, item.About,
		)
	}

	mock.
		ExpectQuery("SELECT id, own_id, first_name, second_name, email, " +
			"city, phone_number, birth_date, sex, citizenship, experience, profession, " +
			"position, wage, education, about FROM resumes").
		WithArgs().
		WillReturnRows(rows)

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	resumes, err := repo.GetResumes()
	fmt.Println(resumes)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(resumes, expect) {
		t.Errorf("results not match,\n want\n%v,\n have\n %v\n", expect, resumes)
		return
	}
}

func TestDBUserStorage_GetResumes_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	mock.
		ExpectQuery("SELECT id, own_id, first_name, second_name, email, " +
			"city, phone_number, birth_date, sex, citizenship, experience, profession, " +
			"position, wage, education, about FROM resumes").
		WithArgs().
		WillReturnError(errors.New("GetResume: error while querying"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	resumes, err := repo.GetResumes()
	fmt.Println(resumes)

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

func TestDBUserStorage_GetResume_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	expect := []Resume{
		{
			ID:          uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a"),
			OwnerID:     uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
			Email:       "",
			City:        "Moscow",
			PhoneNumber: "12345678910",
			FirstName:   "Vova",
			SecondName:  "Zyablikov",
			BirthDate:   "1999-01-08",
			Sex:         "male",
			Citizenship: "Russia",
			Profession:  "programmer",
			Position:    "middle",
			Experience:  "7 years",
			Education:   "MSU",
			Wage:        "100 500.00 руб",
			About:       "Hello employer",
		},
	}

	rows := sqlmock.
		NewRows([]string{"id", "own_id", "first_name", "second_name", "email",
			"city", "phone_number", "birth_date", "sex", "citizenship",
			"experience", "profession", "position", "wage", "education", "about",
		}).AddRow(expect[0].ID.String(), expect[0].OwnerID.String(), expect[0].FirstName, expect[0].SecondName, expect[0].Email, expect[0].City,
		expect[0].PhoneNumber, expect[0].BirthDate, expect[0].Sex, expect[0].Citizenship, expect[0].Experience,
		expect[0].Profession, expect[0].Position, expect[0].Wage, expect[0].Education, expect[0].About,
	)
	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a")
	mock.
		ExpectQuery("SELECT id, own_id, first_name, second_name, email, " +
			"city, phone_number, birth_date, sex, citizenship, experience, profession, " +
			"position, wage, education, about FROM resumes").
		WithArgs(id).
		WillReturnRows(rows)

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	item, err := repo.GetResume(id)
	fmt.Println()

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match,\n want\n%v,\n have\n %v\n", expect[0], item)
		return
	}
}

func TestDBUserStorage_GetResume_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642bbba")
	mock.
		ExpectQuery("SELECT id, own_id, first_name, second_name, email, " +
			"city, phone_number, birth_date, sex, citizenship, experience, profession, " +
			"position, wage, education, about FROM resumes").
		WithArgs(id).
		WillReturnError(errors.New("GetResume: error while querying"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	resume, err := repo.GetResume(id)
	fmt.Println(resume)

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

func TestDBUserStorage_CreateResume_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	resume := Resume{
		ID:          uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d"),
		OwnerID:     uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
		Email:       "email@mail.ru",
		City:        "Moscow",
		PhoneNumber: "12345678910",
		FirstName:   "Vova",
		SecondName:  "Zyablikov",
		BirthDate:   "1999-01-08",
		Sex:         "male",
		Citizenship: "Russia",
		Profession:  "programmer",
		Position:    "middle",
		Experience:  "7 years",
		Education:   "MSU",
		Wage:        "100 500.00 руб",
		About:       "Hello employer",
	}

	mock.
		ExpectExec(`INSERT INTO resumes`).
		WithArgs(
			resume.ID, resume.OwnerID, resume.FirstName, resume.SecondName, resume.Email, resume.City,
			resume.PhoneNumber, resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience,
			resume.Profession, resume.Position, resume.Wage, resume.Education, resume.About,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	ok := repo.CreateResume(resume)

	if !ok {
		t.Error("Failed to create resume\n")
		return
	}
}

func TestDBUserStorage_CreateResume_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	resume := Resume{
		ID:          uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d"),
		OwnerID:     uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
		Email:       "email@mail.ru",
		City:        "Moscow",
		PhoneNumber: "12345678910",
		BirthDate:   "1999-112-100",
		Citizenship: "Russia",
		Profession:  "programmer",
		Position:    "middle",
		Experience:  "7 years",
		Education:   "MSU",
		Wage:        "100 500.00 руб",
		About:       "Hello employer",
	}

	mock.
		ExpectExec(`INSERT INTO resumes`).
		WithArgs(
			resume.ID, resume.OwnerID, resume.FirstName, resume.SecondName, resume.Email, resume.City,
			resume.PhoneNumber, resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience,
			resume.Profession, resume.Position, resume.Wage, resume.Education, resume.About,
		).
		WillReturnError(fmt.Errorf("bad query"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	ok := repo.CreateResume(resume)

	if ok {
		t.Errorf("expected false, got true")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStorage_DeleteResume_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a")

	mock.
		ExpectExec(`DELETE FROM resumes`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	err = repo.DeleteResume(id)
	fmt.Println()

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	mock.
		ExpectQuery("SELECT id, own_id, first_name, second_name, email, " +
			"city, phone_number, birth_date, sex, citizenship, experience, profession, " +
			"position, wage, education, about FROM resumes").
		WithArgs(id).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.GetResume(id)
	fmt.Println()

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

func TestDBUserStorage_DeleteResume_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a")

	mock.
		ExpectExec(`DELETE FROM resumes`).
		WithArgs(id).
		WillReturnError(errors.Errorf("error"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	err = repo.DeleteResume(id)

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

func TestDBUserStorage_PutResume_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	resume := Resume{
		ID:          uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d"),
		OwnerID:     uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
		Email:       "ema@mail.ru",
		City:        "Moscow",
		PhoneNumber: "12345678910",
		FirstName:   "Vova",
		SecondName:  "Zyablikov",
		BirthDate:   "1999-01-08",
		Sex:         "male",
		Citizenship: "Russia",
		Profession:  "programmer",
		Position:    "middle",
		Experience:  "7 years",
		Education:   "MSU",
		Wage:        "100 500.00 руб",
		About:       "Hello employer",
	}

	mock.
		ExpectExec(`UPDATE resumes SET`).
		WithArgs(
			resume.FirstName, resume.SecondName, resume.Email, resume.City, resume.PhoneNumber,
			resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience, resume.Profession,
			resume.Position, resume.Wage, resume.Education, resume.About, resume.ID, resume.OwnerID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	ok := repo.PutResume(resume, resume.OwnerID, resume.ID)

	if !ok {
		t.Error("Failed to put resume\n")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStorage_PutResume_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	resume := Resume{
		ID:          uuid.MustParse("11b77a73-bac7-4597-ab71-7b5fbe53052d"),
		OwnerID:     uuid.MustParse("92b77a73-bac7-4597-ab71-7b5fbe53052d"),
		Email:       "email@mail.ru",
		City:        "Moscow",
		PhoneNumber: "12345678910",
		BirthDate:   "1999-112-100",
		Citizenship: "Russia",
		Profession:  "programmer",
		Position:    "middle",
		Experience:  "7 years",
		Education:   "MSU",
		Wage:        "100 500.00 руб",
		About:       "Hello employer",
	}

	mock.
		ExpectExec(`UPDATE resumes SET`).
		WithArgs(
			resume.FirstName, resume.SecondName, resume.Email, resume.City, resume.PhoneNumber,
			resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience, resume.Profession,
			resume.Position, resume.Wage, resume.Education, resume.About, resume.ID, resume.OwnerID,
		).
		WillReturnError(fmt.Errorf("bad query"))

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	ok := repo.PutResume(resume, resume.OwnerID, resume.ID)

	if ok {
		t.Errorf("expected false, got true")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
