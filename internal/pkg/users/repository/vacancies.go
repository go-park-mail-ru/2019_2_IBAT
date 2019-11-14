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
	_, err := m.DbConn.Exec("INSERT INTO vacancies(id, own_id, experience,"+
		"position, tasks, requirements, conditions, wage_from, wage_to, about)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
		vacancyReg.ID, vacancyReg.OwnerID, vacancyReg.Experience, vacancyReg.Position,
		vacancyReg.Tasks, vacancyReg.Requirements, vacancyReg.Conditions, vacancyReg.WageFrom,
		vacancyReg.WageTo, vacancyReg.About,
	)

	if err != nil {
		fmt.Printf("CreateVacancy: %s\n", err)
		return false
	}

	return true
}

func (m *DBUserStorage) GetVacancy(id uuid.UUID) (Vacancy, error) {

	row := m.DbConn.QueryRowx("SELECT v.id, v.own_id, c.company_name, v.experience,"+
		"v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, "+
		"v.region, v.type_of_employment, v.work_schedule "+
		"FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id WHERE id = $1;", id)

	var vacancy Vacancy
	err := row.StructScan(&vacancy)
	if err != nil {
		fmt.Printf("CreateVacancy: %s\n", err)
		return vacancy, errors.New(InvalidIdMsg)
	}

	return vacancy, nil
}

func (m *DBUserStorage) DeleteVacancy(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM vacancies WHERE id = $1;", id) //check fi invalid id or internal error

	if err != nil {
		fmt.Printf("DeleteVacancy: %s\n", err)
		return errors.New(InvalidIdMsg)
	}

	return nil
}

func (m *DBUserStorage) PutVacancy(vacancy Vacancy, userId uuid.UUID, vacancyId uuid.UUID) bool {

	_, err := m.DbConn.Exec(
		"UPDATE vacancies SET experience = $1, position = $2, tasks = $3, "+
			"requirements = $4, wage_from = $5, wage_to = $6, conditions = $7, about = $8 "+
			"WHERE id = $9 AND own_id = $10;", vacancy.Experience,
		vacancy.Position, vacancy.Tasks, vacancy.Requirements, vacancy.Conditions, vacancy.WageFrom,
		vacancy.WageTo, vacancy.About, vacancyId, userId,
	)

	if err != nil {
		fmt.Printf("PutVacancy: %s\n", err)
		return false
	}

	return true
}

func (m *DBUserStorage) GetVacancies(authInfo AuthStorageValue, params map[string]interface{}) ([]Vacancy, error) {
	vacancies := []Vacancy{}
	log.Printf("Params: %s\n\n", params)
	query := paramsToQuery(params)

	var nmst *sqlx.NamedStmt
	var err error

	if query != "" {
		nmst, err = m.DbConn.PrepareNamed("SELECT v.id, v.own_id, c.company_name, v.experience, " +
			"v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, " +
			"v.region, v.type_of_employment, v.work_schedule " +
			"FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id WHERE " + query)

		if err != nil {
			fmt.Printf("GetVacancies: %s\n", err)
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
			"v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, " +
			"v.region, v.type_of_employment, v.work_schedule " +
			" FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id;")
	}

	if err != nil {
		fmt.Printf("GetVacancies: %s\n", err)
		return vacancies, errors.New(InternalErrorMsg)
	}

	defer rows.Close()

	favVacMap := m.queryFavVacIDs(authInfo.ID)

	for rows.Next() {
		var vacancy Vacancy

		err = rows.StructScan(&vacancy)
		if err != nil {
			fmt.Printf("GetVacancies: %s\n", err)
			return vacancies, errors.New(InternalErrorMsg)
		}

		_, ok := favVacMap[vacancy.ID]
		if ok {
			vacancy.Favorite = true
		}

		vacancies = append(vacancies, vacancy)
	}

	return vacancies, nil
}

func paramsToQuery(params map[string]interface{}) string {
	var query []string

	if params["position"] != nil {
		params["position"] = "%" + params["position"].(string) + "%"
		query = append(query, "position LIKE :position")
	}

	if params["region"] != nil {
		query = append(query, "region = :region")
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

func (m *DBUserStorage) queryFavVacIDs(id uuid.UUID) map[uuid.UUID]bool {
	favVacRows, err := m.DbConn.Queryx("SELECT vacancy_id FROM favorite_vacancies WHERE "+ //err
		"person_id = $1", id)
	if err == nil {
		defer favVacRows.Close()
	}

	favVacMap := map[uuid.UUID]bool{}
	for favVacRows.Next() {
		var id uuid.UUID
		err = favVacRows.Scan(&id)
		if err == nil {
			log.Printf("GetVacancies: %s\n", err)
			favVacMap[id] = true
		}
	}
	return favVacMap
}
