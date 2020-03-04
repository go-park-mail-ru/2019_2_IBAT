package repository

import (
	. "2019_2_IBAT/pkg/pkg/models"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDBUserStorage_GetTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	rows := sqlmock.
		NewRows([]string{"parent_tag", "child_tag"})

	expect := []Tag{
		{
			ParentTag: "Закупки",
			ChildTag:  "управление закупками",
		},
		{
			ParentTag: "Закупки",
			ChildTag:  "металлопрокат",
		},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ParentTag, item.ChildTag)
	}

	mock.
		ExpectQuery("SELECT parent_tag, child_tag FROM tags;").
		WithArgs().
		WillReturnRows(rows)

	tags, err := repo.GetTags()

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(tags, expect) {
		t.Errorf("results not match,\n want\n%v,\n have\n %v\n", expect, tags)
		return
	}

	mock.
		ExpectQuery("SELECT parent_tag, child_tag FROM tags;").
		WithArgs().
		WillReturnError(fmt.Errorf("error"))

	tags, err = repo.GetTags()

	if err == nil {
		t.Errorf("Expected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestDBUserStorage_GetVacancyTagIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	vacancyId := uuid.New()
	rows := sqlmock.
		NewRows([]string{"tag_id"})

	expect := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.String())
	}

	mock.
		ExpectQuery("SELECT tag_id FROM vac_tag_relations AS vr WHERE ").
		WithArgs(vacancyId).
		WillReturnRows(rows)

	tags, err := repo.GetVacancyTagIDs(vacancyId)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(tags, expect) {
		t.Errorf("results not match,\n want\n%v,\n have\n %v\n", expect, tags)
		return
	}

	mock.
		ExpectQuery("SELECT tag_id FROM vac_tag_relations AS vr WHERE ").
		WithArgs(vacancyId).
		WillReturnError(fmt.Errorf("error"))

	tags, err = repo.GetVacancyTagIDs(vacancyId)

	if err == nil {
		t.Errorf("Expected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestDBUserStorage_GetTagIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	repo := DBUserStorage{
		DbConn: sqlxDB,
	}

	rows := sqlmock.
		NewRows([]string{"id"})

	pairs := []Pair{
		{
			First:  "Закупки",
			Second: "управление закупками",
		},
		{
			First:  "Закупки",
			Second: "металлопрокат",
		},
	}

	expect := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	for _, item := range expect {
		rows = rows.AddRow(item.String())
	}

	mock.
		ExpectPrepare("SELECT id FROM tags WHERE ")

	mock.
		ExpectQuery("SELECT id FROM tags WHERE ").
		WillReturnRows(rows)

	tags, err := repo.GetTagIDs(pairs)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(tags, expect) {
		t.Errorf("results not match,\n want\n%v,\n have\n %v\n", expect, tags)
		return
	}

	mock.
		ExpectPrepare("SELECT id FROM tags WHERE ")
	mock.
		ExpectQuery("SELECT id FROM tags WHERE ").
		WithArgs().
		WillReturnError(fmt.Errorf("error"))

	tags, err = repo.GetTagIDs(pairs)

	if err == nil {
		t.Errorf("Expected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.
		ExpectPrepare("SELECT id FROM tags WHERE ").
		WillReturnError(fmt.Errorf("error"))

	tags, err = repo.GetTagIDs(pairs)

	if err == nil {
		t.Errorf("Expected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

/*func (m *DBUserStorage) GetTagIDs(spheres []Pair) []uuid.UUID {
	var tagIds []uuid.UUID

	if !(len(spheres) > 0) {
		return tagIds
	}

	var nmstTags *sqlx.NamedStmt
	var err error

	sphQuery, sphMap := buildSpheresQuery(spheres)

	nmstTags, err = m.DbConn.PrepareNamed("SELECT id FROM tags WHERE " + sphQuery)
	if err != nil {
		return tagIds
	} //real error message

	tagRows, err := nmstTags.Queryx(sphMap)
	if err != nil {
		return tagIds
	} //real error message

	if err == nil && sphQuery != "" {
		defer tagRows.Close()
		for tagRows.Next() {
			var tagId uuid.UUID

			err = tagRows.Scan(&tagId)
			if err != nil {
				log.Printf("GetVacancies: %s\n", err)
			}

			tagIds = append(tagIds, tagId)
		}
	}

	return tagIds
}*/
