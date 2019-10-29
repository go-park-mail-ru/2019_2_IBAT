package users

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DBUserStorage struct {
	DbConn *sqlx.DB
}

func (m *DBUserStorage) CreateSeeker(seekerInput SeekerReg) (uuid.UUID, bool) {
	id := uuid.New()

	_, err := m.DbConn.Exec(
		"INSERT INTO persons(id, email, first_name, second_name, password_hash, role)"+
			"VALUES($1, $2, $3, $4, $5, $6)", id, seekerInput.Email, seekerInput.FirstName,
		seekerInput.SecondName, seekerInput.Password, SeekerStr,
	)

	if err != nil {
		fmt.Println("error while creating user")
		return uuid.UUID{}, false
	}

	return id, true
}

func (m DBUserStorage) CreateEmployer(employerInput EmployerReg) (uuid.UUID, bool) {
	tx, err := m.DbConn.Begin()
	if err != nil {
		fmt.Println("error while creating user")
		return uuid.UUID{}, false
	}
	id := uuid.New()

	_, err = m.DbConn.Exec(
		"INSERT INTO persons(id, email, first_name, second_name, password_hash, role)"+
			"VALUES($1, $2, $3, $4, $5, $6)", id, employerInput.Email, employerInput.FirstName,
		employerInput.SecondName, employerInput.Password, EmployerStr,
	)

	if err != nil {
		fmt.Println("error while creating person")
		return uuid.UUID{}, false
	}

	_, err = m.DbConn.Exec(
		"INSERT INTO companies(own_id, company_name, site, phone_number, "+
			"extra_phone_number, city, empl_num)VALUES($1, $2, $3, $4, $5, $6, $7)",
		id, employerInput.CompanyName, employerInput.Site, employerInput.PhoneNumber,
		employerInput.ExtraPhoneNumber, employerInput.City, employerInput.EmplNum,
	)

	if err != nil {
		fmt.Println("error while creating company")
		tx.Rollback()
		return uuid.UUID{}, false
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println("error while creating user")
		tx.Rollback()
		return uuid.UUID{}, false
	}

	return id, true
}

func (m DBUserStorage) DeleteUser(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM persons WHERE id = $1", id)
	if err != nil {
		fmt.Println("DeleteUser: error while deleting")
		return err
	}
	return nil
}

func (m DBUserStorage) GetSeekers() ([]Seeker, error) { //not tested
	seekers := []Seeker{}

	rows, err := m.DbConn.Queryx("SELECT id, email, first_name, second_name,"+
		"path_to_image FROM persons WHERE role = $1;", SeekerStr)
	defer rows.Close()

	if err != nil {
		fmt.Println("GetSeeker: error while query seekers")
		return seekers, err
	}

	for rows.Next() {
		seek := Seeker{}
		_ = rows.StructScan(&seek)
		// if err != nil {
		// 	return seekers, err
		// }

		id_rows, err := m.DbConn.Query("SELECT v.id FROM resumes AS v WHERE v.own_id = $1;", seek.ID)
		if err != nil {
			fmt.Println("GetSeeker: error while query resumes")
			return seekers, err
		}
		defer id_rows.Close()

		resumes := make([]uuid.UUID, 0)

		for id_rows.Next() {
			var id uuid.UUID
			_ = id_rows.Scan(&id)
			// if err != nil {
			// 	return employers, err
			// }
			resumes = append(resumes, id)
		}

		seek.Resumes = resumes
		seekers = append(seekers, seek)
	}

	return seekers, nil
}

func (m DBUserStorage) CheckUser(email string, password string) (uuid.UUID, string, bool) {
	resId := uuid.UUID{}
	var class string
	fmt.Printf("CheckUser: %s\n", email)
	fmt.Printf("CheckUser: %s\n", password)

	row := m.DbConn.QueryRowx("SELECT id, role FROM persons "+
		"WHERE email = $1 AND password_hash = $2;", email, password)

	err := row.Scan(&resId, &class)

	if err != nil {
		fmt.Println("CheckUser: Scan error")
		return resId, class, false
	}

	return resId, class, true
}

func (m DBUserStorage) CreateResume(resumeReg Resume, userId uuid.UUID) (uuid.UUID, bool) {
	id := uuid.New()

	_, err := m.DbConn.Exec("INSERT INTO resumes(id, own_id, first_name, second_name, email, "+
		"city, phone_number, birth_date, sex, citizenship, experience, profession, "+
		"position, wage, education, about)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);",
		id, userId, resumeReg.FirstName, resumeReg.SecondName, resumeReg.Email, resumeReg.City,
		resumeReg.PhoneNumber, resumeReg.BirthDate, resumeReg.Sex, resumeReg.Citizenship, resumeReg.Experience,
		resumeReg.Profession, resumeReg.Position, resumeReg.Wage, resumeReg.Education, resumeReg.About,
	)

	if err != nil {
		fmt.Println("CreateResume: error while creating")
		return id, false
	}

	return id, true
}

func (m DBUserStorage) CreateVacancy(vacancyReg Vacancy, userId uuid.UUID) (uuid.UUID, bool) {
	id := uuid.New()

	_, err := m.DbConn.Exec("INSERT INTO vacancies(id, own_id, experience, profession,"+
		"position, tasks, requirements, conditions, wage, about)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);", id,
		userId, vacancyReg.Experience, vacancyReg.Profession, vacancyReg.Position,
		vacancyReg.Tasks, vacancyReg.Requirements, vacancyReg.Conditions, vacancyReg.Wage, vacancyReg.About,
	)

	if err != nil {
		fmt.Println("CreateVacancy: error while creating")
		return id, false
	}

	return id, true
}

func (m DBUserStorage) GetResume(id uuid.UUID) (Resume, error) {

	row := m.DbConn.QueryRowx("SELECT id, own_id, first_name, second_name, email, "+
		"city, phone_number, birth_date, sex, citizenship, experience, profession, "+
		"position, wage, education, about FROM resumes WHERE id = $1;", id,
	)

	var resume Resume
	_ = row.StructScan(&resume)

	return resume, nil
}

func (m DBUserStorage) DeleteResume(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM resumes WHERE id = $1;", id)

	if err != nil {
		fmt.Println("DeleteResume: error while deleting")
		return err
	}

	return nil
}

func (m DBUserStorage) GetSeeker(id uuid.UUID) (Seeker, error) {
	rows := m.DbConn.QueryRowx("SELECT id, email, first_name, second_name,"+
		" path_to_image FROM persons WHERE id = $1;", id)

	seeker := Seeker{}
	_ = rows.StructScan(&seeker)
	// if err != nil {
	// 	return seekers, err
	// }

	id_rows, err := m.DbConn.Query("SELECT r.id FROM resumes AS r WHERE r.own_id = $1;", seeker.ID)

	if err != nil {
		fmt.Println("GetSeeker: error while query resumes")
		return seeker, err
	}
	defer id_rows.Close()

	resumes := make([]uuid.UUID, 0)

	for id_rows.Next() {
		var id uuid.UUID
		_ = id_rows.Scan(&id)
		// if err != nil {
		// 	return employers, err
		// }
		resumes = append(resumes, id)
	}

	seeker.Resumes = resumes

	return seeker, nil
}

func (m DBUserStorage) GetEmployer(id uuid.UUID) (Employer, error) {

	rows := m.DbConn.QueryRowx("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site,"+
		"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image,"+
		"c.city FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE p.id = $1;", id) //and p.class

	empl := Employer{}
	_ = rows.StructScan(&empl)
	// if err != nil {
	// 	return employers, err
	// }

	id_rows, err := m.DbConn.Query("SELECT v.id FROM vacancies AS v WHERE v.own_id = $1", empl.ID)
	if err != nil {
		return empl, err
	}
	defer id_rows.Close()

	vacancies := make([]uuid.UUID, 0)

	for id_rows.Next() {
		var id uuid.UUID
		_ = id_rows.Scan(&id)
		// if err != nil {
		// 	return employers, err
		// }
		vacancies = append(vacancies, id)
	}

	empl.Vacancies = vacancies

	return empl, nil
}

func (m DBUserStorage) PutSeeker(seekerInput SeekerReg, id uuid.UUID) bool {

	_, err := m.DbConn.Exec(
		"UPDATE persons SET email = $1, first_name = $2, second_name = $3, password_hash = $4"+
			" WHERE id = $5;", seekerInput.Email, seekerInput.FirstName,
		seekerInput.SecondName, seekerInput.Password, id,
	)

	if err != nil {
		fmt.Println("error while changing user")
		return false
	}

	return true
}

func (m DBUserStorage) PutEmployer(employerInput EmployerReg, id uuid.UUID) bool {
	tx, err := m.DbConn.Begin()

	if err != nil {
		fmt.Println("error while changing employer")
		return false
	}

	_, err = m.DbConn.Exec(
		"UPDATE persons SET email = $1, first_name = $2, second_name = $3, password_hash = $4"+
			"WHERE id = $5;", employerInput.Email, employerInput.FirstName,
		employerInput.SecondName, employerInput.Password, id,
	)

	if err != nil {
		fmt.Println("error while changing employer")
		return false
	}

	_, err = m.DbConn.Exec(
		"UPDATE companies SET company_name = $1, site = $2, phone_number = $3, "+
			"extra_phone_number = $4, city = $5, empl_num = $6 WHERE own_id = $7;",
		employerInput.CompanyName, employerInput.Site, employerInput.PhoneNumber,
		employerInput.ExtraPhoneNumber, employerInput.City, employerInput.EmplNum, id,
	)

	if err != nil {
		fmt.Println("error while changing employer")
		tx.Rollback()
		return false
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println("error while changing employer")
		tx.Rollback()
		return false
	}

	return true
}

func (m DBUserStorage) PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool {

	_, err := m.DbConn.Exec("UPDATE resumes SET "+
		"first_name = $1, second_name = $2, email = $3, "+
		"city = $4, phone_number = $5, birth_date = $6, sex = $7, citizenship = $8, "+
		"experience = $9, profession = $10, position = $11, wage = $12, education = $13, about = $14 "+
		"WHERE id = $15 AND own_id = $16;",
		resume.FirstName, resume.SecondName, resume.Email, resume.City, resume.PhoneNumber,
		resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience, resume.Profession,
		resume.Position, resume.Wage, resume.Education, resume.About, resumeId, userId,
	)

	if err != nil {
		fmt.Println("PutResume: error while changing")
		return false
	}

	return true
}

func (m *DBUserStorage) GetEmployers() ([]Employer, error) {
	employers := []Employer{}

	rows, err := m.DbConn.Queryx("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site,"+
		"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image, c.city, "+
		" c.description FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE role = $1;", EmployerStr)
	defer rows.Close()

	if err != nil {
		return employers, err
	}

	for rows.Next() {
		empl := Employer{}
		_ = rows.StructScan(&empl)
		// if err != nil {
		// 	return employers, err
		// }

		id_rows, err := m.DbConn.Query("SELECT v.id FROM vacancies AS v WHERE v.own_id = $1", empl.ID)
		if err != nil {
			return employers, err
		}
		defer id_rows.Close()

		vacancies := make([]uuid.UUID, 0)

		for id_rows.Next() {
			var id uuid.UUID
			_ = id_rows.Scan(&id)
			// if err != nil {
			// 	return employers, err
			// }
			vacancies = append(vacancies, id)
		}

		empl.Vacancies = vacancies
		employers = append(employers, empl)
	}

	return employers, nil
}

func (m DBUserStorage) GetResumes() ([]Resume, error) {

	resumes := []Resume{}

	rows, err := m.DbConn.Queryx("SELECT id, own_id, first_name, second_name, email, " +
		"city, phone_number, birth_date, sex, citizenship, experience, profession, " +
		"position, wage, education, about FROM resumes;",
	)
	if err != nil {
		return resumes, err
	}

	for rows.Next() {
		var resume Resume

		_ = rows.StructScan(&resume)
		if err != nil {
			return resumes, err
		}

		resumes = append(resumes, resume)
	}

	return resumes, nil
}

func (m DBUserStorage) GetVacancy(id uuid.UUID) (Vacancy, error) {
	row := m.DbConn.QueryRowx("SELECT id, own_id, experience, profession,"+
		"position, tasks, requirements, conditions, wage, about "+
		" FROM vacancies WHERE id = $1;", id)

	var vacancy Vacancy
	_ = row.StructScan(&vacancy)
	// fmt.Println(vacancy)
	id_row := m.DbConn.QueryRow("SELECT company_name FROM companies WHERE own_id = $1;", vacancy.OwnerID)
	_ = id_row.Scan(&vacancy.CompanyName)

	return vacancy, nil
}

func (m DBUserStorage) DeleteVacancy(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM vacancies WHERE id = $1;", id)

	if err != nil {
		fmt.Println("DeleteVacancy: error while deleting")
		return err
	}

	return nil
}

func (m DBUserStorage) PutVacancy(vacancy Vacancy, userId uuid.UUID, vacancyId uuid.UUID) bool {

	_, err := m.DbConn.Exec(
		"UPDATE vacancies SET experience = $1, profession = $2, position = $3, tasks = $4, "+
			"requirements = $5, wage = $6, conditions = $7, about = $8 "+
			"WHERE id = $9 AND own_id = $10;", vacancy.Experience, vacancy.Profession,
		vacancy.Position, vacancy.Tasks, vacancy.Requirements, vacancy.Wage, vacancy.Conditions,
		vacancy.About, vacancyId, userId,
	)

	if err != nil {
		fmt.Println("PutVacancy: error while changing")
		return false
	}

	return true
}

func (m DBUserStorage) GetVacancies() ([]Vacancy, error) {
	vacancies := []Vacancy{}

	rows, err := m.DbConn.Queryx("SELECT v.id, v.own_id, v.experience, v.profession," +
		"v.position, v.tasks, v.requirements, v.conditions, v.wage, v.about, c.company_name" +
		" FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id;")

	if err != nil {
		return vacancies, err
	}

	for rows.Next() {
		var vacancy Vacancy

		_ = rows.StructScan(&vacancy)
		// if err != nil {
		// 	return vacancies, err
		// }

		vacancies = append(vacancies, vacancy)
	}

	return vacancies, nil
}

func (m DBUserStorage) SetImage(id uuid.UUID, class string, imageName string) bool {

	_, err := m.DbConn.Exec(
		"UPDATE persons SET path_to_image = $1 WHERE id = $2", imageName, id,
	)

	if err != nil {
		fmt.Println("error while setting image to user")
		return false
	}

	return true
}

func (m DBUserStorage) CreateRespond(respond Respond, userId uuid.UUID) (uuid.UUID, bool) {
	id := uuid.New()
	fmt.Printf("respond.ResumeID: %s\n", respond.ResumeID)

	resume, err := m.GetResume(respond.ResumeID)
	if err != nil {
		fmt.Println("CreateRespond: no such resume error\n")
		return id, false
	}

	if resume.OwnerID != userId {
		fmt.Println(resume)
		fmt.Printf("Userid: %s\n", userId)
		fmt.Printf("resume.OwnerID: %s\n", resume.OwnerID)

		fmt.Println("CreateRespond: forbidden error\n")
		return id, false
	}

	_, err = m.DbConn.Exec("INSERT INTO responds(resume_id, vacancy_id, status)"+
		"VALUES($1, $2, $3);",
		respond.ResumeID, respond.VacancyID, respond.Status,
	)

	if err != nil {
		fmt.Println("CreateRespond: error while creating")
		return id, false
	}

	return id, true
}

func (m DBUserStorage) GetResponds(record AuthStorageValue, params map[string]string) ([]Respond, error) {

	responds := []Respond{}
	log.Println("GetResponds: start")
	var rows *sqlx.Rows

	if params["resumeid"] == "" && params["vacancyid"] == "" {
		if record.Role == EmployerStr {
			rows, _ = m.DbConn.Queryx("SELECT DISTINCT r.resume_id, r.vacancy_id, r.status "+
				"FROM vacancies AS v  JOIN responds AS r ON v.id = r.vacancy_id "+
				"JOIN persons AS p ON v.own_id = $1", record.ID) ///fix join
		} else if record.Role == SeekerStr {
			rows, _ = m.DbConn.Queryx("SELECT DISTINCT r.resume_id, r.vacancy_id, r.status "+
				"FROM resumes AS res JOIN responds AS r ON res.id = r.resume_id "+
				"JOIN persons AS p ON res.own_id = $1", record.ID) ///fix join
		}

	}

	if params["resumeid"] != "" {
		row := m.DbConn.QueryRow("SELECT own_id FROM resumes "+
			"WHERE id = $1 AND own_id = $2;", params["resumeid"], record.ID)
		id := uuid.UUID{}
		_ = row.Scan(&id)
		if id != record.ID {
			return responds, errors.New("Error")
		}
		rows, _ = m.DbConn.Queryx("SELECT resume_id, vacancy_id, status"+
			" FROM responds WHERE resume_id = $1;", params["resumeid"])
	}

	if params["vacancyid"] != "" {
		row := m.DbConn.QueryRow("SELECT own_id FROM vacancies "+
			"WHERE id = $1 AND own_id = $2;", params["vacancyid"], record.ID)
		id := uuid.UUID{}
		_ = row.Scan(&id)
		if id != record.ID {
			return responds, errors.New("Error")
		}
		rows, _ = m.DbConn.Queryx("SELECT resume_id, vacancy_id, status"+
			" FROM responds WHERE vacancy_id = $1;", params["vacancyid"])
	}

	for rows.Next() {
		var respond Respond

		_ = rows.StructScan(&respond)
		// if err != nil {
		// 	return vacancies, err
		// }

		responds = append(responds, respond)
	}

	return responds, nil
}
