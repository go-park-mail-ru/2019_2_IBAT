package repository

import (
	. "2019_2_IBAT/pkg/pkg/models"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDBUserStorage_SetTagsIDs_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer sqlxDB.Close()

	tagIDs := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	authRec := AuthStorageValue{
		ID:   uuid.New(),
		Role: SeekerStr,
	}

	for _, tagID := range tagIDs {
		mock.
			ExpectExec(`INSERT INTO recommendations`).
			WithArgs(
				authRec.ID, tagID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	repo := DBRecommendsStorage{
		DbConn: sqlxDB,
	}

	err = repo.SetTagIDs(authRec, tagIDs)

	if err != nil {
		t.Error("Failed to set tag IDs\n")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStorage_GetTagIDs_Correct(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	rows := sqlmock.
		NewRows([]string{"tag_id"})

	expect := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	authRec := AuthStorageValue{
		ID: uuid.New(),
	}

	mock.
		ExpectQuery(`SELECT tag_id FROM recommendations WHERE`).
		WithArgs(authRec.ID).
		WillReturnRows(rows)

	repo := DBRecommendsStorage{
		DbConn: sqlxDB,
	}

	gotTagIDs, err := repo.GetTagIDs(authRec)

	if err != nil {
		t.Error("Failed to get tag IDs\n")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.Equal(t, expect, gotTagIDs, "The two values should be the same.")
}

func TestDBUserStorage_GetTagIDs_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	rows := sqlmock.
		NewRows([]string{"tag_id"})

	expect := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	authRec := AuthStorageValue{
		ID: uuid.New(),
	}

	mock.
		ExpectQuery(`SELECT tag_id FROM recommendations WHERE`).
		WithArgs(authRec.ID).
		WillReturnError(fmt.Errorf("bad query"))

	repo := DBRecommendsStorage{
		DbConn: sqlxDB,
	}

	_, err = repo.GetTagIDs(authRec)

	if err == nil {
		t.Error("Expected error")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.Equal(t, InternalErrorMsg, err.Error(), "The two values should be the same.")
}

func TestDBUserStorage_GetUsersForTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}

	repo := DBRecommendsStorage{
		DbConn: sqlxDB,
	}

	rows := sqlmock.
		NewRows([]string{"tag_id"})

	expectUserIDs := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	tagIDs := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	for _, item := range expectUserIDs {
		rows = rows.AddRow(item)
	}

	mock.
		ExpectQuery(`SELECT DISTINCT person_id FROM recommendations WHERE`).
		WillReturnRows(rows)

	gotUserIDs, err := repo.GetUsersForTags(tagIDs)

	if err != nil {
		t.Error("Failed to get tag IDs\n")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.Equal(t, expectUserIDs, gotUserIDs, "The two values should be the same.")

	mock.
		ExpectQuery(`SELECT DISTINCT person_id FROM recommendations WHERE`).
		WillReturnError(fmt.Errorf(InternalErrorMsg))

	_, err = repo.GetUsersForTags(tagIDs)

	if err == nil {
		t.Error("Expected error\n")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.Equal(t, InternalErrorMsg, err.Error(), "The two values should be the same.")
}
