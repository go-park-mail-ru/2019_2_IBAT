package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m *DBUserStorage) GetTags() ([]Tag, error) {
	tags := []Tag{}

	rows, err := m.DbConn.Queryx("SELECT parent_tag, child_tag FROM tags;")

	if err != nil {
		fmt.Printf("GetTags: %s\n", err)
		return tags, errors.New(InternalErrorMsg)
	}
	defer rows.Close()

	for rows.Next() {
		tag := Tag{}
		err = rows.StructScan(&tag)
		if err != nil {
			fmt.Printf("GetTags: %s\n", err)
			return tags, errors.New(InternalErrorMsg)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (m *DBUserStorage) GetVacancyTagIDs(vacancyId uuid.UUID) ([]uuid.UUID, error) {
	var tagIDs []uuid.UUID

	rows, err := m.DbConn.Queryx("SELECT tag_id FROM vac_tag_relations AS vr "+
		"WHERE vr.vacancy_id = $1;", vacancyId,
	)

	if err != nil {
		log.Printf("GetVacancyTagIDs: %s", err)
		return tagIDs, err
	}

	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			fmt.Printf("GetTags: %s\n", err)
			// return tags, errors.New(InternalErrorMsg)
		}

		tagIDs = append(tagIDs, id)
	}

	return tagIDs, nil
}
