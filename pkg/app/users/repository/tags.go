package repository

import (
	. "2019_2_IBAT/pkg/pkg/models"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func buildSpheresQuery(spheres []Pair) (string, map[string]interface{}) {

	sphMap := make(map[string]interface{})
	var sphQueryArr []string

	for i, item := range spheres {
		parent_tag := "parent_tag" + strconv.Itoa(i)
		child_tag := "child_tag" + strconv.Itoa(i)
		sphMap[parent_tag] = item.First
		sphMap[child_tag] = item.Second
		sphQueryArr = append(sphQueryArr, "(parent_tag = :"+parent_tag+" AND child_tag = :"+child_tag+" )")

	}

	sphQuery := strings.Join(sphQueryArr, " OR ")

	return sphQuery, sphMap
}

func (m *DBUserStorage) GetTagIDs(spheres []Pair) ([]uuid.UUID, error) {
	var tagIds []uuid.UUID

	if !(len(spheres) > 0) {
		return tagIds, nil
	}

	var nmstTags *sqlx.NamedStmt
	var err error

	sphQuery, sphMap := buildSpheresQuery(spheres)

	nmstTags, err = m.DbConn.PrepareNamed("SELECT id FROM tags WHERE " + sphQuery)
	if err != nil {
		log.Printf("GetTagIDs: %s\n", err)
		return tagIds, fmt.Errorf(InternalErrorMsg)
	} //real error message

	tagRows, err := nmstTags.Queryx(sphMap)
	if err != nil {
		log.Printf("GetTagIDs: %s\n", err)
		return tagIds, fmt.Errorf(InternalErrorMsg)
	} //real error message

	if err == nil && sphQuery != "" {
		defer tagRows.Close()
		for tagRows.Next() {
			var tagId uuid.UUID

			err = tagRows.Scan(&tagId)
			if err != nil {
				log.Printf("GetTagIDs: %s\n", err)
			}

			tagIds = append(tagIds, tagId)
		}
	}

	return tagIds, nil
}
