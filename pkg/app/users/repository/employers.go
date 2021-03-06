package repository

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	. "2019_2_IBAT/pkg/pkg/models"
)

func (m *DBUserStorage) CreateEmployer(employerInput Employer) bool {
	tx, err := m.DbConn.Begin()
	if err != nil {
		fmt.Printf("CreateEmployer: %s\n", err)
		return false
	}
	// id := uuid.New()

	// salt := make([]byte, 8)
	// rand.Read(salt)
	// employerInput.Password = string(passwords.HashPass(salt, employerInput.Password))
	// fmt.Printf("Password hash: %s\n", employerInput.Password)

	_, err = m.DbConn.Exec(
		"INSERT INTO persons(id, email, first_name, second_name, password_hash, role, path_to_image)"+
			"VALUES($1, $2, $3, $4, $5, $6, $7);", employerInput.ID, employerInput.Email, employerInput.FirstName,
		employerInput.SecondName, employerInput.Password, EmployerStr, employerInput.PathToImg,
	)

	if err != nil {
		fmt.Printf("CreateEmployer: %s\n", err)
		tx.Rollback()
		return false
	}

	_, err = m.DbConn.Exec(
		"INSERT INTO companies(own_id, company_name, site, phone_number, "+
			"extra_phone_number, region, empl_num)VALUES($1, $2, $3, $4, $5, $6, $7);",
		employerInput.ID, employerInput.CompanyName, employerInput.Site, employerInput.PhoneNumber,
		employerInput.ExtraPhoneNumber, employerInput.Region, employerInput.EmplNum,
	)

	if err != nil {
		fmt.Printf("CreateEmployer: %s\n", err)
		tx.Rollback()
		return false
	}

	err = tx.Commit()

	if err != nil {
		fmt.Printf("CreateEmployer: %s\n", err)
		tx.Rollback()
		return false
	}

	return true
}

func (m *DBUserStorage) GetEmployer(id uuid.UUID) (Employer, error) {
	log.Println("GetEmployer Repository Start")
	row := m.DbConn.QueryRowx("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site,"+
		"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image, "+
		"c.region, c.description "+
		"FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE p.id = $1;", id) //and p.class

	empl := Employer{}
	err := row.StructScan(&empl)
	if err != nil {
		fmt.Printf("GetEmployer: %s\n", err)
		return empl, err
	}

	id_rows, err := m.DbConn.Query("SELECT v.id FROM vacancies AS v WHERE v.own_id = $1;", empl.ID)
	if err != nil {
		fmt.Printf("GetEmployer: %s\n", err)
		return empl, errors.New(InternalErrorMsg)
	}
	defer id_rows.Close()

	vacancies := make([]uuid.UUID, 0)

	for id_rows.Next() {
		var id uuid.UUID
		err = id_rows.Scan(&id)
		if err != nil {
			fmt.Printf("GetEmployer: %s\n", err)
			return empl, err
		}
		vacancies = append(vacancies, id)
	}

	empl.Vacancies = vacancies

	return empl, nil
}

func (m *DBUserStorage) PutEmployer(employerInput EmployerReg, id uuid.UUID) bool {
	tx, err := m.DbConn.Begin()

	if err != nil {
		fmt.Printf("PutEmployer: %s\n", err)
		return false
	}

	_, err = m.DbConn.Exec(
		"UPDATE persons SET email = $1, first_name = $2, second_name = $3, password_hash = $4"+
			"WHERE id = $5;", employerInput.Email, employerInput.FirstName,
		employerInput.SecondName, employerInput.Password, id,
	)

	if err != nil {
		fmt.Printf("PutEmployer: %s\n", err)
		tx.Rollback()
		return false
	}

	_, err = m.DbConn.Exec(
		"UPDATE companies SET company_name = $1, site = $2, phone_number = $3, "+
			"extra_phone_number = $4, region = $5, empl_num = $6 WHERE own_id = $7;",
		employerInput.CompanyName, employerInput.Site, employerInput.PhoneNumber,
		employerInput.ExtraPhoneNumber, employerInput.Region, employerInput.EmplNum, id,
	)

	if err != nil {
		fmt.Printf("PutEmployer: %s\n", err)
		tx.Rollback()
		return false
	}

	err = tx.Commit()

	if err != nil {
		fmt.Printf("PutEmployer: %s\n", err)
		tx.Rollback()
		return false
	}

	return true
}

func (m *DBUserStorage) GetEmployers(params map[string]interface{}) ([]Employer, error) {
	employers := []Employer{}

	query := paramsEmplToQuery(params)

	var nmst *sqlx.NamedStmt
	var err error

	if query != "" {
		nmst, err = m.DbConn.PrepareNamed("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site," +
			"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image, c.region, " +
			" c.description FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE p.role = 'employer' AND " + query + ";")

		if err != nil {
			log.Printf("GetEmployers: error while preparing statement - %s", err)
			return employers, errors.New(InternalErrorMsg)
		}
	} else {
		log.Println("GetEmployers: query is empty")
	}

	var rows *sqlx.Rows

	if query != "" {
		rows, err = nmst.Queryx(params)
	} else {
		rows, err = m.DbConn.Queryx("SELECT p.id, p.email, c.company_name, p.first_name, p.second_name, c.site,"+
			"c.empl_num, c.phone_number, c.extra_phone_number, c.spheres_of_work, p.path_to_image, c.region, "+
			" c.description FROM persons as p JOIN companies as c ON p.id = c.own_id WHERE role = $1;", EmployerStr)
	}

	if err != nil {
		log.Printf("GetEmployers: %s", err)
		return employers, errors.New(InternalErrorMsg)
	}
	defer rows.Close()

	for rows.Next() {
		empl := Employer{}
		err = rows.StructScan(&empl)
		if err != nil {
			fmt.Printf("GetEmployers: %s\n", err)
			return employers, err
		}

		id_rows, err := m.DbConn.Query("SELECT v.id FROM vacancies AS v WHERE v.own_id = $1;", empl.ID)
		if err != nil {
			fmt.Printf("GetEmployers: %s\n", err)
			return employers, errors.New(InternalErrorMsg)
		}
		defer id_rows.Close()

		vacancies := make([]uuid.UUID, 0)

		for id_rows.Next() {
			var id uuid.UUID
			err = id_rows.Scan(&id)
			if err != nil {
				fmt.Printf("GetEmployers: %s\n", err)
				return employers, err
			}
			vacancies = append(vacancies, id)
		}

		empl.Vacancies = vacancies
		employers = append(employers, empl)
	}

	return employers, nil
}

func paramsEmplToQuery(params map[string]interface{}) string {
	var query []string

	if params["company_name"] != nil {
		query = append(query, "company_name = :company_name")
	} else {

		if params["empl_num"] != nil {
			query = append(query, "empl_num >= :empl_num")
		}

		if params["region"] != nil {
			query = append(query, "region = :region")
		}
	}

	str := strings.Join(query, " AND ")

	log.Printf("Query: %s", str)
	return str
}
