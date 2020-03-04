package repository

// import (
// 	. "2019_2_IBAT/pkg/pkg/models"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/jmoiron/sqlx"
// 	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
// )

// func TestDBUserStorage_CreateRespond_Correct(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	defer db.Close()
// 	sqlxDB := sqlx.NewDb(db, "sqlmock")

// 	if err != nil {
// 		t.Fatalf("cant create mock: %s", err)
// 	}
// 	defer sqlxDB.Close()

// 	personID := uuid.New()
// 	resumeID := uuid.New()

// 	resume := Resume{
// 		ID:               resumeID,
// 		OwnerID:          personID,
// 		Email:            "",
// 		Region:           "Moscow",
// 		PhoneNumber:      "12345678910",
// 		FirstName:        "Vova",
// 		SecondName:       "Zyablikov",
// 		BirthDate:        "1999-01-08",
// 		Sex:              "male",
// 		Citizenship:      "Russia",
// 		Position:         "programmer",
// 		Experience:       "7 years",
// 		Education:        "MSU",
// 		Wage:             "100 500.00 руб",
// 		About:            "Hello employer",
// 		TypeOfEmployment: "",
// 		WorkSchedule:     "",
// 	}

// 	respond := Respond{
// 		Status:    AwaitSt,
// 		ResumeID:  resumeID,
// 		VacancyID: uuid.New(),
// 	}

// 	rows := sqlmock.
// 		NewRows([]string{"id", "own_id", "first_name", "second_name", "email",
// 			"region", "phone_number", "birth_date", "sex", "citizenship",
// 			"experience", "position", "wage", "education", "about", "work_schedule", "type_of_employment",
// 		}).AddRow(resume.ID.String(), resume.OwnerID.String(), resume.FirstName, resume.SecondName, resume.Email, resume.Region,
// 		resume.PhoneNumber, resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience,
// 		resume.Position, resume.Wage, resume.Education, resume.About,
// 		resume.WorkSchedule, resume.TypeOfEmployment,
// 	)
// 	id := uuid.MustParse("f14c6104-3430-413b-ab4e-e31c8642ad8a")
// 	mock.
// 		ExpectQuery("SELECT id, own_id, first_name, second_name, email, " +
// 			"region, phone_number, birth_date, sex, citizenship, experience, " +
// 			"position, wage, education, about, work_schedule, type_of_employment FROM resumes").
// 		WithArgs(id).
// 		WillReturnRows(rows)

// 	mock.
// 		ExpectQuery("INSERT INTO responds").
// 		WithArgs(
// 			respond.ResumeID, respond.VacancyID, respond.Status,
// 		).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 		// WillReturnResult(sqlmock.NewResult(1, 1))
// 	repo := DBUserStorage{
// 		DbConn: sqlxDB,
// 	}
// 	ok := repo.CreateRespond(respond, personID)

// 	if !ok {
// 		t.Error("Failed to create vacancy\n")
// 		return
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }
