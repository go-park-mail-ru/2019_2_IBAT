package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (m *DBUserStorage) CreateRespond(respond Respond, userId uuid.UUID) (uuid.UUID, bool) {
	id := uuid.New()
	fmt.Printf("respond.ResumeID: %s\n", respond.ResumeID)

	resume, err := m.GetResume(respond.ResumeID)
	if err != nil {
		fmt.Println("CreateRespond: no such resume error")
		return id, false
	}

	if resume.OwnerID != userId {
		fmt.Println(resume)
		fmt.Printf("Userid: %s\n", userId)
		fmt.Printf("resume.OwnerID: %s\n", resume.OwnerID)

		fmt.Println("CreateRespond: forbidden error")
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

func (m *DBUserStorage) GetResponds(record AuthStorageValue, params map[string]string) ([]Respond, error) {

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
