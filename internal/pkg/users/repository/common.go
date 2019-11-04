package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBUserStorage struct {
	DbConn *sqlx.DB
}

func (m *DBUserStorage) DeleteUser(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM persons WHERE id = $1", id)
	if err != nil {
		fmt.Println("DeleteUser: error while deleting")
		return err
	}
	return nil
}

func (m *DBUserStorage) CheckUser(email string, password string) (uuid.UUID, string, bool) {

	row := m.DbConn.QueryRow("SELECT id, role FROM persons "+
		"WHERE email = $1 AND password_hash = $2;", email, password)

	resId := uuid.UUID{}
	var class string
	err := row.Scan(&resId, &class)

	if err != nil {
		fmt.Println("CheckUser: Scan error")
		fmt.Printf("CheckUser: resId - %s\n", resId)
		fmt.Printf("CheckUser: class - %s\n", class)
		return resId, class, false
	}

	return resId, class, true
}

func (m *DBUserStorage) SetImage(id uuid.UUID, class string, imageName string) bool {

	_, err := m.DbConn.Exec(
		"UPDATE persons SET path_to_image = $1 WHERE id = $2", imageName, id,
	)

	if err != nil {
		fmt.Println("error while setting image to user")
		return false
	}

	return true
}
