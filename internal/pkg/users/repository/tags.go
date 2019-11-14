package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"

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
