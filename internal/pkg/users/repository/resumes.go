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

func (m *DBUserStorage) CreateResume(resumeReg Resume) bool {
	_, err := m.DbConn.Exec("INSERT INTO resumes(id, own_id, first_name, second_name, email, "+
		"region, phone_number, birth_date, sex, citizenship, experience, "+
		"position, wage, education, about)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);",
		resumeReg.ID, resumeReg.OwnerID, resumeReg.FirstName, resumeReg.SecondName, resumeReg.Email, resumeReg.Region,
		resumeReg.PhoneNumber, resumeReg.BirthDate, resumeReg.Sex, resumeReg.Citizenship, resumeReg.Experience,
		resumeReg.Position, resumeReg.Wage, resumeReg.Education, resumeReg.About,
	)

	if err != nil {
		fmt.Printf("GetResponds: %s\n", err)
		return false
	}

	return true
}

func (m *DBUserStorage) GetResume(id uuid.UUID) (Resume, error) {

	row := m.DbConn.QueryRowx("SELECT id, own_id, first_name, second_name, email, "+
		"region, phone_number, birth_date, sex, citizenship, experience, "+
		"position, wage, education, about, work_schedule, type_of_employment FROM resumes WHERE id = $1;", id,
	)

	var resume Resume
	err := row.StructScan(&resume)
	if err != nil {
		fmt.Printf("GetResume: %s\n", err)
		return resume, errors.New(InvalidIdMsg)
	}

	return resume, nil
}

func (m *DBUserStorage) DeleteResume(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM resumes WHERE id = $1;", id)

	if err != nil {
		fmt.Printf("DeleteResume: %s\n", err)
		return errors.New(InternalErrorMsg)
	}

	return nil
}

func (m *DBUserStorage) PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool {

	_, err := m.DbConn.Exec("UPDATE resumes SET "+
		"first_name = $1, second_name = $2, email = $3, "+
		"region = $4, phone_number = $5, birth_date = $6, sex = $7, citizenship = $8, "+
		"experience = $9, position = $10, wage = $11, education = $12, about = $13 "+
		"WHERE id = $15 AND own_id = $16;",
		resume.FirstName, resume.SecondName, resume.Email, resume.Region, resume.PhoneNumber,
		resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience,
		resume.Position, resume.Wage, resume.Education, resume.About, resumeId, userId,
	)

	if err != nil {
		fmt.Printf("PutResume: %s\n", err)
		return false
	}

	return true
}

func (m *DBUserStorage) GetResumes(authInfo AuthStorageValue, params map[string]interface{}) ([]Resume, error) {
	resumes := []Resume{}

	query := paramsToResumesQuery(params)

	var nmst *sqlx.NamedStmt
	var err error

	if query != "" && params["own"] == nil {
		nmst, err = m.DbConn.PrepareNamed("SELECT id, own_id, first_name, second_name, email, " +
			"region, phone_number, birth_date, sex, citizenship, experience," +
			"position, wage, education, about, work_schedule, type_of_employment FROM resumes WHERE " + query)

		if err != nil {
			fmt.Printf("GetResumes: %s\n", err)
			return resumes, errors.New(InternalErrorMsg)
		}
	} else {
		log.Println("GetResumes: query is empty")
	}

	var rows *sqlx.Rows
	if query != "" && params["own"] == nil {
		rows, err = nmst.Queryx(params)
	} else if params["own"] != nil {
		rows, err = m.DbConn.Queryx("SELECT id, own_id, first_name, second_name, email, "+
			"region, phone_number, birth_date, sex, citizenship, experience, "+
			"position, wage, education, about, work_schedule, "+
			"type_of_employment FROM resumes WHERE own_id = $1;", authInfo.ID,
		)
	} else {
		rows, err = m.DbConn.Queryx("SELECT id, own_id, first_name, second_name, email, " +
			"region, phone_number, birth_date, sex, citizenship, experience, " +
			"position, wage, education, about, work_schedule, type_of_employment FROM resumes;",
		)
	}

	if err != nil {
		fmt.Printf("GetResumes: %s\n", err)
		return resumes, errors.New(InternalErrorMsg)
	}
	defer rows.Close()

	for rows.Next() {
		var resume Resume

		err = rows.StructScan(&resume)

		if err != nil {
			fmt.Printf("GetResumes: %s\n", err)
			return resumes, errors.New(InternalErrorMsg)
		}
		resumes = append(resumes, resume)
	}

	return resumes, nil
}

func paramsToResumesQuery(params map[string]interface{}) string {
	var query []string

	if params["position"] != nil {
		params["position"] = "%" + params["position"].(string) + "%"
		query = append(query, "position LIKE :position")
	}

	if params["region"] != nil {
		query = append(query, "region= :region")
	}

	if params["wage_from"] != nil {
		query = append(query, "wage >= :wage_from")
	}

	if params["wage_to"] != nil {
		query = append(query, "wage <= :wage_to")
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
