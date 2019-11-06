package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DBUserStorage struct {
	DbConn *sqlx.DB
}

func (m *DBUserStorage) DeleteUser(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM persons WHERE id = $1", id)
	if err != nil {
		fmt.Println("DeleteUser: error while deleting")
		return errors.New("DeleteUser: error while deleting")
	}
	return nil
}

func (m *DBUserStorage) CheckUser(email string, password string) (uuid.UUID, string, bool) {

	// password = passwords.CheckPass()
	row := m.DbConn.QueryRow("SELECT id, role, password_hash FROM persons "+
		"WHERE email = $1", email)

	resId := uuid.UUID{}
	var class string
	var password_hash []byte
	err := row.Scan(&resId, &class, &password_hash)

	// if !passwords.CheckPass(password_hash, password) || err != nil {
	// 	return resId, class, false
	// }
	if password != string(password_hash) || err != nil {
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
