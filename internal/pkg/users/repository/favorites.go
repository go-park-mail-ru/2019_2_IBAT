package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func (m *DBUserStorage) CreateFavorite(favVac FavoriteVacancy) bool {

	_, err := m.DbConn.Exec("INSERT INTO favorite_vacancies(person_id, vacancy_id)"+
		"VALUES($1, $2);",
		favVac.PersonID, favVac.VacancyID,
	)

	if err != nil {
		fmt.Println("CreateFavorite: error while creating")
		return false
	}

	return true
}

func (m *DBUserStorage) GetFavoriteVacancies(record AuthStorageValue) ([]Vacancy, error) {
	vacancies := []Vacancy{}

	rows, err := m.DbConn.Queryx("SELECT v.id, v.own_id, c.company_name, v.experience,"+
		"v.profession, v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, "+
		"v.region, v.type_of_employment, v.work_schedule "+

		"FROM favorite_vacancies AS fv "+
		"JOIN vacancies AS v ON (fv.vacancy_id = v.id) "+
		"JOIN companies AS c ON v.own_id = c.own_id WHERE fv.person_id = $1;", record.ID) //fux query
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