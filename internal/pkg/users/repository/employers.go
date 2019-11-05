package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"

	"github.com/google/uuid"
)

func (m *DBUserStorage) CreateEmployer(employerInput EmployerReg) (uuid.UUID, bool) {
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

func (m *DBUserStorage) GetEmployer(id uuid.UUID) (Employer, error) {

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

func (m *DBUserStorage) PutEmployer(employerInput EmployerReg, id uuid.UUID) bool {
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
