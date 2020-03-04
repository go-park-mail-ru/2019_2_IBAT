package repository

import (
	. "2019_2_IBAT/pkg/pkg/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (m *DBUserStorage) CreateRespond(respond Respond, userId uuid.UUID) bool {
	fmt.Printf("respond.ResumeID: %s\n", respond.ResumeID)

	resume, err := m.GetResume(respond.ResumeID)
	if err != nil {
		fmt.Printf("CreateRespond: %s\n", err)
		return false
	}

	if resume.OwnerID != userId {
		// fmt.Println(resume)
		// fmt.Printf("Userid: %s\n", userId)
		// fmt.Printf("resume.OwnerID: %s\n", resume.OwnerID)

		fmt.Println("CreateRespond: forbidden error")
		return false
	}

	_, err = m.DbConn.Exec("INSERT INTO responds(resume_id, vacancy_id, status)"+
		"VALUES($1, $2, $3);",
		respond.ResumeID, respond.VacancyID, respond.Status,
	)

	if err != nil {
		fmt.Printf("CreateRespond: %s\n", err)
		return false
	}

	return true
}

func (m *DBUserStorage) GetResponds(record AuthStorageValue, params map[string]string) ([]Respond, error) {

	responds := []Respond{}
	var rows *sqlx.Rows
	var err error

	if params["resume_id"] == "" && params["vacancy_id"] == "" {
		if record.Role == EmployerStr {
			rows, err = m.DbConn.Queryx("SELECT DISTINCT r.resume_id, r.vacancy_id, r.status "+
				"FROM vacancies AS v  JOIN responds AS r ON v.id = r.vacancy_id "+
				"JOIN persons AS p ON v.own_id = $1;", record.ID) ///fix join
		} else if record.Role == SeekerStr {
			rows, err = m.DbConn.Queryx("SELECT DISTINCT r.resume_id, r.vacancy_id, r.status "+
				"FROM resumes AS res JOIN responds AS r ON res.id = r.resume_id "+
				"JOIN persons AS p ON res.own_id = $1;", record.ID) ///fix join
		}

	}

	if err != nil {
		fmt.Printf("GetResponds: %s\n", err)
		return responds, errors.New(InternalErrorMsg)
	}

	if params["resume_id"] != "" {
		row := m.DbConn.QueryRow("SELECT own_id FROM resumes "+
			"WHERE id = $1 AND own_id = $2;", params["resume_id"], record.ID)
		id := uuid.UUID{}
		err = row.Scan(&id)
		if err != nil {
			fmt.Printf("GetResponds: %s\n", err)
			return responds, errors.New(ForbiddenMsg)
		}

		rows, err = m.DbConn.Queryx("SELECT resume_id, vacancy_id, status"+
			" FROM responds WHERE resume_id = $1;", params["resume_id"])
		if err != nil {
			fmt.Printf("GetResponds: %s\n", err)
			return responds, errors.New(InternalErrorMsg)
		}
	}

	if params["vacancy_id"] != "" {
		row := m.DbConn.QueryRow("SELECT own_id FROM vacancies "+
			"WHERE id = $1 AND own_id = $2;", params["vacancy_id"], record.ID)
		id := uuid.UUID{}

		err = row.Scan(&id)
		if err != nil {
			fmt.Printf("GetResponds: %s\n", err)
			return responds, errors.New(ForbiddenMsg)
		}

		rows, err = m.DbConn.Queryx("SELECT resume_id, vacancy_id, status"+
			" FROM responds WHERE vacancy_id = $1;", params["vacancy_id"])
		if err != nil {
			fmt.Printf("GetResponds: %s\n", err)
			return responds, errors.New(InternalErrorMsg)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var respond Respond

		err = rows.StructScan(&respond)
		if err != nil {
			fmt.Printf("GetResponds: %s\n", err)
			return responds, errors.New(InternalErrorMsg)
		}

		responds = append(responds, respond)
	}

	return responds, nil
}
