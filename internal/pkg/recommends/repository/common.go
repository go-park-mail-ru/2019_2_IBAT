package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBRecommendsStorage struct {
	DbConn *sqlx.DB
}

func (m DBRecommendsStorage) SetTagIDs(AuthRec AuthStorageValue, tagIDs []uuid.UUID) error {

	for _, id := range tagIDs {
		_, _ = m.DbConn.Exec(
			"INSERT INTO recommendations(person_id, tag_id)VALUES"+
				"($1, $2);", AuthRec.ID, id,
		)
	} //make by one insert

	return nil
}

func (m DBRecommendsStorage) GetTagIDs(AuthRec AuthStorageValue) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	rows, err := m.DbConn.Query("SELECT tag_id FROM recommendations WHERE "+
		"person_id = $1;", AuthRec.ID,
	)

	if err != nil {
		fmt.Printf("GetTagIDs: %s\n", err)
		return ids, errors.New(InternalErrorMsg)
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID

		err = rows.Scan(&id)

		if err != nil {
			fmt.Printf("GetTagIDs: %s\n", err)
			return ids, errors.New(InternalErrorMsg)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
