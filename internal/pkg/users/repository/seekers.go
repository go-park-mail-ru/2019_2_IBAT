package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"

	"github.com/google/uuid"
)

func (m *DBUserStorage) CreateSeeker(seekerInput SeekerReg) (uuid.UUID, bool) {
	id := uuid.New()

	_, err := m.DbConn.Exec(
		"INSERT INTO persons(id, email, first_name, second_name, password_hash, role)"+
			"VALUES($1, $2, $3, $4, $5, $6)", id, seekerInput.Email, seekerInput.FirstName,
		seekerInput.SecondName, seekerInput.Password, SeekerStr,
	)

	if err != nil {
		fmt.Println("error while creating user")
		return uuid.UUID{}, false
	}

	return id, true
}

func (m *DBUserStorage) GetSeekers() ([]Seeker, error) { //not tested
	seekers := []Seeker{}

	rows, err := m.DbConn.Queryx("SELECT id, email, first_name, second_name,"+
		"path_to_image FROM persons WHERE role = $1;", SeekerStr)
	if err != nil {
		fmt.Println("GetSeeker: error while query seekers")
		return seekers, err
	}
	defer rows.Close()

	for rows.Next() {
		seek := Seeker{}
		_ = rows.StructScan(&seek)
		// if err != nil {
		// 	return seekers, err
		// }

		id_rows, err := m.DbConn.Query("SELECT r.id FROM resumes AS r WHERE r.own_id = $1;", seek.ID)
		if err != nil {
			fmt.Println("GetSeeker: error while query resumes")
			return seekers, err
		}
		defer id_rows.Close()

		resumes := make([]uuid.UUID, 0)

		for id_rows.Next() {
			var id uuid.UUID
			_ = id_rows.Scan(&id)
			// if err != nil {
			// 	return employers, err
			// }
			resumes = append(resumes, id)
		}

		seek.Resumes = resumes
		seekers = append(seekers, seek)
	}

	return seekers, nil
}

func (m *DBUserStorage) GetSeeker(id uuid.UUID) (Seeker, error) {

	rows := m.DbConn.QueryRowx("SELECT id, email, first_name, second_name,"+
		" path_to_image FROM persons WHERE id = $1;", id)

	seeker := Seeker{}
	_ = rows.StructScan(&seeker)
	// if err != nil {
	// 	return seekers, err
	// }

	id_rows, err := m.DbConn.Query("SELECT r.id FROM resumes AS r WHERE r.own_id = $1;", seeker.ID)

	if err != nil {
		fmt.Println("GetSeeker: error while query resumes")
		return seeker, err
	}
	defer id_rows.Close()

	resumes := make([]uuid.UUID, 0)

	for id_rows.Next() {
		var id uuid.UUID
		_ = id_rows.Scan(&id)
		// if err != nil {
		// 	return employers, err
		// }
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
		fmt.Println("error while changing user")
		return false
	}

	return true
}
