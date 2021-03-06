package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	. "2019_2_IBAT/pkg/pkg/models"
)

// go test ./... -coverprofile cover.out; go tool cover -func cover.out
// cat cover.out | fgrep -v _easyjson.go > cover.tmp
// go tool cover -func cover.tmp

type DBRecommendsStorage struct {
	DbConn *sqlx.DB
}

func (m DBRecommendsStorage) SetTagIDs(AuthRec AuthStorageValue, tagIDs []string) error {

	for _, id := range tagIDs {
		_, err := m.DbConn.Exec(
			"INSERT INTO recommendations(person_id, tag_id)VALUES"+
				"($1, $2);", AuthRec.ID, id,
		)
		log.Printf("SetTagIDs %s", err)
	} //make by one insert

	return nil
}

func (m DBRecommendsStorage) GetTagIDs(AuthRec AuthStorageValue) ([]string, error) {
	var ids []string
	rows, err := m.DbConn.Query("SELECT tag_id FROM recommendations WHERE "+
		"person_id = $1;", AuthRec.ID,
	)

	if err != nil {
		fmt.Printf("GetTagIDs: %s\n", err)
		return ids, errors.New(InternalErrorMsg)
	}
	defer rows.Close()

	for rows.Next() {
		var id string

		err = rows.Scan(&id)

		if err != nil {
			fmt.Printf("GetTagIDs: %s\n", err)
			return ids, errors.New(InternalErrorMsg)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (m DBRecommendsStorage) GetUsersForTags(tagIDs []string) ([]string, error) {
	var ids []string

	params := make(map[string]interface{})
	params["ids"] = tagIDs

	query, args, err := sqlx.Named("SELECT DISTINCT person_id FROM recommendations WHERE "+
		"tag_id IN (:ids);", params,
	)

	if err != nil {
		fmt.Printf("GetUsersForTags: %s\n", err)
		return ids, errors.New(InternalErrorMsg)
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		fmt.Printf("GetUsersForTags: %s\n", err)
		return ids, errors.New(InternalErrorMsg)
	}

	query = m.DbConn.Rebind(query)
	rows, err := m.DbConn.Queryx(query, args...)
	if err != nil {
		fmt.Printf("GetUsersForTags: %s\n", err)
		return ids, errors.New(InternalErrorMsg)
	}

	defer rows.Close()

	for rows.Next() {
		var id string

		err = rows.Scan(&id)

		if err != nil {
			fmt.Printf("GetUsersForTags: %s\n", err)
			return ids, errors.New(InternalErrorMsg)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
