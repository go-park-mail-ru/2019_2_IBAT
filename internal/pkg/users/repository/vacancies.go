package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (m *DBUserStorage) CreateVacancy(vacancyReg Vacancy) bool {
	_, err := m.DbConn.Exec("INSERT INTO vacancies(id, own_id, experience, profession,"+
		"position, tasks, requirements, conditions, wage_from, wage_to, about)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
		vacancyReg.ID, vacancyReg.OwnerID, vacancyReg.Experience, vacancyReg.Profession, vacancyReg.Position,
		vacancyReg.Tasks, vacancyReg.Requirements, vacancyReg.Conditions, vacancyReg.WageFrom,
		vacancyReg.WageTo, vacancyReg.About,
	)

	if err != nil {
		fmt.Println("CreateVacancy: error while creating")
		return false
	}

	return true
}

func (m *DBUserStorage) GetVacancy(id uuid.UUID) (Vacancy, error) {

	row := m.DbConn.QueryRowx("SELECT v.id, v.own_id, c.company_name, v.experience,"+
		"v.profession, v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, "+
		"v.region, v.type_of_employment, v.work_schedule "+
		"FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id WHERE id = $1;", id)

	var vacancy Vacancy
	err := row.StructScan(&vacancy)
	if err != nil {
		log.Println("GetVacancy: error while querying")
		return Vacancy{}, errors.New(InvalidIdMsg)
	}
	log.Println("Storage: GetVacancy\n vacancy:")
	log.Println(vacancy)

	return vacancy, nil
}

func (m *DBUserStorage) DeleteVacancy(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM vacancies WHERE id = $1;", id)

	if err != nil {
		fmt.Println("DeleteVacancy: error while deleting")
		return err
	}

	return nil
}

func (m *DBUserStorage) PutVacancy(vacancy Vacancy, userId uuid.UUID, vacancyId uuid.UUID) bool {

	_, err := m.DbConn.Exec(
		"UPDATE vacancies SET experience = $1, profession = $2, position = $3, tasks = $4, "+
			"requirements = $5, wage_from = $6, wage_to = $7, conditions = $8, about = $9 "+
			"WHERE id = $10 AND own_id = $11;", vacancy.Experience, vacancy.Profession,
		vacancy.Position, vacancy.Tasks, vacancy.Requirements, vacancy.Conditions, vacancy.WageFrom,
		vacancy.WageTo, vacancy.About, vacancyId, userId,
	)

	if err != nil {
		fmt.Println(BadRequestMsg)
		return false
	}

	return true
}

func (m *DBUserStorage) GetVacancies(params map[string]interface{}) ([]Vacancy, error) {
	vacancies := []Vacancy{}
	log.Printf("Params: %s\n\n", params)
	query := paramsToQuery(params)

	var nmst *sqlx.NamedStmt
	var err error

	if query != "" {
		nmst, err = m.DbConn.PrepareNamed("SELECT v.id, v.own_id, c.company_name, v.experience, " +
			"v.profession, v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, " +
			"v.region, v.type_of_employment, v.work_schedule " +
			"FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id WHERE " + query)
		if err != nil {
			log.Println("GetVacancies: error while preparing statement")
			return vacancies, errors.New(InternalErrorMsg)
		}
	} else {
		log.Println("GetVacancies: query is empty")
	}

	var rows *sqlx.Rows

	if query != "" {
		rows, err = nmst.Queryx(params)
	} else {
		rows, err = m.DbConn.Queryx("SELECT v.id, v.own_id, c.company_name, v.experience," +
			"v.profession, v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about" +
			" FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id;")
	}
	defer rows.Close()

	if err != nil {
		log.Println("GetVacancies: error while query")
		return vacancies, errors.New(InternalErrorMsg)
	}
	for rows.Next() {
		var vacancy Vacancy

		_ = rows.StructScan(&vacancy)

		vacancies = append(vacancies, vacancy)
	}

	return vacancies, nil
}

func paramsToQuery(params map[string]interface{}) string {
	var query []string

	if params["region"] != nil {
		query = append(query, "region= :region")
	}

	if params["wage_from"] != nil {
		query = append(query, "wage_to >= :wage_from")
	}

	if params["experience"] != nil {
		query = append(query, "experience = :experience")
	}

	if params["type_of_employment"] != nil {
		query = append(query, "type_of_employment=:type_of_employment")
	}

	if params["work_schedule"] != nil {
		query = append(query, "work_schedule = :work_schedule")
	}

	str := strings.Join(query, " AND ")

	log.Printf("Query: %s", str)
	return str
}
