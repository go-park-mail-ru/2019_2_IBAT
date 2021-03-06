package repository

import (
	. "2019_2_IBAT/pkg/pkg/models"

	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m *DBUserStorage) CreateSeeker(seekerInput Seeker) bool {

	// salt := make([]byte, 8)
	// rand.Read(salt)
	// seekerInput.Password = string(passwords.HashPass(salt, seekerInput.Password))

	_, err := m.DbConn.Exec(
		"INSERT INTO persons(id, email, first_name, second_name, password_hash, role, path_to_image)"+
			"VALUES($1, $2, $3, $4, $5, $6, $7);", seekerInput.ID, seekerInput.Email, seekerInput.FirstName,
		seekerInput.SecondName, seekerInput.Password, SeekerStr, seekerInput.PathToImg,
	)

	if err != nil {
		fmt.Printf("CreateSeeker: %s\n", err)
		return false
	}

	return true
}

func (m *DBUserStorage) GetSeekers() ([]Seeker, error) { //not tested
	seekers := []Seeker{}

	rows, err := m.DbConn.Queryx("SELECT id, email, first_name, second_name,"+
		"path_to_image FROM persons WHERE role = $1;", SeekerStr)
	if err != nil {
		fmt.Printf("GetSeekers: %s\n", err)
		return seekers, errors.New(InternalErrorMsg)
	}
	defer rows.Close()

	for rows.Next() {
		seek := Seeker{}
		err = rows.StructScan(&seek)
		if err != nil {
			fmt.Printf("GetSeekers: %s\n", err)
			return seekers, errors.New(InternalErrorMsg)
		}

		id_rows, err := m.DbConn.Query("SELECT r.id FROM resumes AS r WHERE r.own_id = $1;", seek.ID)
		if err != nil {
			fmt.Printf("GetSeekers: %s\n", err)
			return seekers, errors.New(InternalErrorMsg)
		}
		defer id_rows.Close()

		resumes := make([]uuid.UUID, 0)

		for id_rows.Next() {
			var id uuid.UUID
			err = id_rows.Scan(&id)
			if err != nil {
				fmt.Printf("GetSeekers: %s\n", err)
				return seekers, errors.New(InternalErrorMsg)
			}
			resumes = append(resumes, id)
		}

		seek.Resumes = resumes
		seekers = append(seekers, seek)
	}

	return seekers, nil
}

func (m *DBUserStorage) GetSeeker(id uuid.UUID) (Seeker, error) {

	row := m.DbConn.QueryRowx("SELECT id, email, first_name, second_name,"+
		" path_to_image FROM persons WHERE id = $1;", id)

	seeker := Seeker{}
	err := row.StructScan(&seeker)
	if err != nil {
		fmt.Printf("GetSeeker: %s\n", err)
		return seeker, errors.New(InternalErrorMsg)
	}

	id_rows, err := m.DbConn.Query("SELECT r.id FROM resumes AS r WHERE r.own_id = $1;", seeker.ID)

	if err != nil {
		fmt.Printf("GetSeeker: %s\n", err)
		return seeker, errors.New(InternalErrorMsg)
	}
	defer id_rows.Close()

	resumes := make([]uuid.UUID, 0)

	for id_rows.Next() {
		var id uuid.UUID
		err = id_rows.Scan(&id)

		if err != nil {
			fmt.Printf("GetSeeker: %s\n", err)
			return seeker, errors.New(InternalErrorMsg)
		}

		resumes = append(resumes, id)
	}

	seeker.Resumes = resumes

	return seeker, nil
}

func (m *DBUserStorage) PutSeeker(seekerInput SeekerReg, id uuid.UUID) bool {

	_, err := m.DbConn.Exec(
		"UPDATE persons SET email = $1, first_name = $2, second_name = $3, password_hash = $4"+
			" WHERE id = $5;", seekerInput.Email, seekerInput.FirstName,
		seekerInput.SecondName, seekerInput.Password, id,
	)

	if err != nil {
		fmt.Printf("PutSeeker: %s\n", err)
		return false
	}

	return true
}
