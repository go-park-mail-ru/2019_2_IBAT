package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	. "2019_2_IBAT/pkg/pkg/models"
)

func (m *DBUserStorage) CreateFavorite(favVac FavoriteVacancy) bool {

	_, err := m.DbConn.Exec("INSERT INTO favorite_vacancies(person_id, vacancy_id)"+
		"VALUES($1, $2);",
		favVac.PersonID, favVac.VacancyID,
	)

	if err != nil {
		fmt.Printf("CreateFavorite: %s \n", err)
		return false
	}

	return true
}

func (m *DBUserStorage) GetFavoriteVacancies(record AuthStorageValue) ([]Vacancy, error) {
	vacancies := []Vacancy{}

	rows, err := m.DbConn.Queryx("SELECT v.id, v.own_id, c.company_name, v.experience, "+
		"v.position, v.tasks, v.requirements, v.wage_from, v.wage_to, v.conditions, v.about, "+
		"v.region, v.type_of_employment, v.work_schedule "+

		"FROM favorite_vacancies AS fv "+
		"JOIN vacancies AS v ON (fv.vacancy_id = v.id) "+
		"JOIN companies AS c ON v.own_id = c.own_id WHERE fv.person_id = $1;", record.ID) //fux query

	if err != nil {
		fmt.Printf("GetFavoriteVacancies: %s\n", err)
		return vacancies, errors.New(InternalErrorMsg)
	}

	defer rows.Close()

	for rows.Next() {
		var vacancy Vacancy

		err = rows.StructScan(&vacancy)
		if err != nil {
			fmt.Printf("GetFavoriteVacancies: %s\n", err)
			return vacancies, errors.New(InternalErrorMsg)
		}
		vacancy.Favorite = true
		vacancies = append(vacancies, vacancy)
	}

	return vacancies, nil
}

func (m *DBUserStorage) DeleteFavoriteVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) error {
	_, err := m.DbConn.Exec("DELETE FROM favorite_vacancies WHERE vacancy_id = $1 AND person_id = $2;",
		vacancyId, authInfo.ID,
	) //check fi invalid id or internal error

	if err != nil {
		fmt.Printf("DeleteVacancy: %s\n", err)
		return errors.New(InvalidIdMsg) //dif errors
	}

	return nil
}
