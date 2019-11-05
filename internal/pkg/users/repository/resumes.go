package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m *DBUserStorage) CreateResume(resumeReg Resume) bool {
	_, err := m.DbConn.Exec("INSERT INTO resumes(id, own_id, first_name, second_name, email, "+
		"city, phone_number, birth_date, sex, citizenship, experience, profession, "+
		"position, wage, education, about)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);",
		resumeReg.ID, resumeReg.OwnerID, resumeReg.FirstName, resumeReg.SecondName, resumeReg.Email, resumeReg.City,
		resumeReg.PhoneNumber, resumeReg.BirthDate, resumeReg.Sex, resumeReg.Citizenship, resumeReg.Experience,
		resumeReg.Profession, resumeReg.Position, resumeReg.Wage, resumeReg.Education, resumeReg.About,
	)

	if err != nil {
		fmt.Println("CreateResume: error while creating")
		return false
	}

	return true
}

func (m *DBUserStorage) GetResume(id uuid.UUID) (Resume, error) {

	row := m.DbConn.QueryRowx("SELECT id, own_id, first_name, second_name, email, "+
		"city, phone_number, birth_date, sex, citizenship, experience, profession, "+
		"position, wage, education, about FROM resumes WHERE id = $1;", id,
	)

	var resume Resume
	err := row.StructScan(&resume)
	if err != nil {
		log.Println("GetResume: error while querying")
		return Resume{}, errors.New("GetResume: error while querying")
	}
	log.Println("Storage: GetResume\n Resume:")
	log.Println(resume)

	return resume, nil
}

func (m *DBUserStorage) DeleteResume(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM resumes WHERE id = $1;", id)

	if err != nil {
		fmt.Println("DeleteResume: error while deleting")
		return err
	}

	return nil
}

func (m *DBUserStorage) PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool {

	_, err := m.DbConn.Exec("UPDATE resumes SET "+
		"first_name = $1, second_name = $2, email = $3, "+
		"city = $4, phone_number = $5, birth_date = $6, sex = $7, citizenship = $8, "+
		"experience = $9, profession = $10, position = $11, wage = $12, education = $13, about = $14 "+
		"WHERE id = $15 AND own_id = $16;",
		resume.FirstName, resume.SecondName, resume.Email, resume.City, resume.PhoneNumber,
		resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience, resume.Profession,
		resume.Position, resume.Wage, resume.Education, resume.About, resumeId, userId,
	)

	if err != nil {
		fmt.Println("PutResume: error while changing")
		return false
	}

	return true
}

func (m *DBUserStorage) GetResumes() ([]Resume, error) {

	resumes := []Resume{}

	rows, err := m.DbConn.Queryx("SELECT id, own_id, first_name, second_name, email, " +
		"city, phone_number, birth_date, sex, citizenship, experience, profession, " +
		"position, wage, education, about FROM resumes;",
	)
	if err != nil {
		log.Println("GetResumes: error while querying")
		return resumes, errors.New("GetResumes: error while querying")
	}

	for rows.Next() {
		var resume Resume

		_ = rows.StructScan(&resume)

		resumes = append(resumes, resume)
	}

	return resumes, nil
}
