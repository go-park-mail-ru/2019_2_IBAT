package repository

import (
	. "2019_2_IBAT/pkg/pkg/models"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const DoesNotMatterString = "Не имеет значения"

func (m *DBUserStorage) CreateVacancy(vacancyReg Vacancy) bool {
	_, err := m.DbConn.Exec("INSERT INTO vacancies(id, own_id, experience,"+
		"position, tasks, requirements, conditions, wage_from, wage_to, about)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
		vacancyReg.ID, vacancyReg.OwnerID, vacancyReg.Experience, vacancyReg.Position,
		vacancyReg.Tasks, vacancyReg.Requirements, vacancyReg.Conditions, vacancyReg.WageFrom,
		vacancyReg.WageTo, vacancyReg.About,
	)

	if err != nil {
		log.Printf("CreateVacancy: %s\n", err)
		return false
	}

	for _, item := range vacancyReg.Spheres {
		_, err := m.DbConn.Exec("INSERT INTO vac_tag_relations(tag_id, vacancy_id)VALUES"+
			"((SELECT id from tags WHERE parent_tag = $1 AND child_tag = $2), $3);",
			item.First, item.Second, vacancyReg.ID,
		)
		if err != nil {
			log.Printf("CreateVacancy: %s\n", err)
		}
	}

	return true
}

func (m *DBUserStorage) GetVacancy(id uuid.UUID, userId uuid.UUID) (Vacancy, error) {

	row := m.DbConn.QueryRowx("SELECT v.id, v.own_id, c.company_name, v.experience,"+
		"v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, "+
		"v.region, v.type_of_employment, v.work_schedule "+
		"FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id WHERE id = $1;", id)

	var vacancy Vacancy
	err := row.StructScan(&vacancy)
	if err != nil {
		log.Printf("GetVacancy: %s\n", err)
		return vacancy, errors.New(InvalidIdMsg)
	}

	favVacRow := m.DbConn.QueryRowx("SELECT vacancy_id FROM favorite_vacancies WHERE "+
		"person_id = $1 AND vacancy_id = $2;", userId, id)

	var resId uuid.UUID
	err = favVacRow.Scan(&resId)
	if err == nil {
		vacancy.Favorite = true
	}

	isRepondedRows, err := m.DbConn.Queryx("SELECT rps.vacancy_id FROM resumes AS rms "+
		"JOIN responds AS rps ON (rms.id = rps.resume_id) WHERE "+
		"rms.own_id = $1 AND rps.vacancy_id = $2;", userId, id)
	if err != nil {
		log.Printf("GetVacancy: %s\n", err)
	} else {
		defer isRepondedRows.Close()

		var vacId uuid.UUID
		isRepondedRows.Next()
		err = isRepondedRows.Scan(&vacId)
		log.Printf("GetVacancy: %s\n", err)
		if err == nil {
			vacancy.IsResponded = true
		}
	}
	return vacancy, nil
}

func (m *DBUserStorage) DeleteVacancy(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM vacancies WHERE id = $1;", id) //check fi invalid id or internal error

	if err != nil {
		log.Printf("DeleteVacancy: %s\n", err)
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
		log.Printf("PutVacancy: %s\n", err)
		return false
	}

	return true
}

func (m *DBUserStorage) queryFavVacIDs(id uuid.UUID) map[uuid.UUID]bool {
	favVacMap := map[uuid.UUID]bool{}
	favVacRows, err := m.DbConn.Queryx("SELECT vacancy_id FROM favorite_vacancies WHERE "+ //err
		"person_id = $1", id)
	if err == nil {
		defer favVacRows.Close()
	} else {
		return favVacMap
	}

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

func (m *DBUserStorage) GetVacancies(authInfo AuthStorageValue, params map[string]interface{}) ([]Vacancy, error) {
	vacancies := []Vacancy{}
	log.Printf("Params: %s\n\n", params)

	paramsStr := paramsToQuery(params)

	var tagIDsLength int
	if params["tag_ids"] != nil {
		tagIDsLength = len(params["tag_ids"].([]uuid.UUID))
	}

	if tagIDsLength > 0 && paramsStr != "" {
		paramsStr = " AND " + paramsStr
	}

	var rows *sqlx.Rows

	var query string
	var args []interface{}
	var err error
	if tagIDsLength > 0 {
		query, args, err = sqlx.Named("SELECT v.id, v.own_id, c.company_name, v.experience,"+
			"v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, "+
			"v.region, v.type_of_employment, v.work_schedule "+
			" FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id "+
			" INNER JOIN vac_tag_relations AS vt ON v.id = vt.vacancy_id WHERE vt.tag_id IN (:tag_ids)"+paramsStr, params)
	} else if paramsStr != "" {
		query, args, err = sqlx.Named("SELECT v.id, v.own_id, c.company_name, v.experience,"+
			"v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, "+
			"v.region, v.type_of_employment, v.work_schedule "+
			" FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id "+
			" WHERE "+paramsStr, params)
	}

	if err != nil {
		log.Printf("GetVacancies: %s\n", err)
		return vacancies, errors.New(InternalErrorMsg)
	}

	if query != "" || tagIDsLength > 0 {
		query, args, err = sqlx.In(query, args...)
		if err != nil {
			log.Printf("GetVacancies: %s\n", err)
			return vacancies, errors.New(InternalErrorMsg)
		}

		query = m.DbConn.Rebind(query)

		rows, err = m.DbConn.Queryx(query, args...)
	} else {
		rows, err = m.DbConn.Queryx("SELECT v.id, v.own_id, c.company_name, v.experience," +
			"v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, " +
			"v.region, v.type_of_employment, v.work_schedule " +
			" FROM vacancies AS v JOIN companies AS c ON v.own_id = c.own_id;")
	}

	if err != nil {
		log.Printf("GetVacancies: %s\n", err)
		return vacancies, errors.New(InternalErrorMsg)
	}

	defer rows.Close()

	favVacMap := m.queryFavVacIDs(authInfo.ID)

	for rows.Next() {
		var vacancy Vacancy

		err = rows.StructScan(&vacancy)
		if err != nil {
			log.Printf("GetVacancies: %s\n", err)
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
		query = append(query, "position ILIKE :position")
	}

	if params["region"] != nil {
		query = append(query, "v.region = :region")
	}

	if params["wage_from"] != nil {
		query = append(query, "wage_to >= :wage_from")
	}

	if params["experience"] != nil {
		if params["experience"].(string) != DoesNotMatterString {
			query = append(query, "experience = :experience")
		}
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
