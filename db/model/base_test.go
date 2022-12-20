package model

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	var (
		title = "test title"
		id    = 1
	)

	b := &base[Todo]{
		db: gDB,
	}
	mockSQL.MatchExpectationsInOrder(false)
	mockSQL.ExpectBegin()
	mockSQL.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "todos"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	mockSQL.ExpectCommit()
	res, err := b.Create(&Todo{
		Title: title,
		Done:  true,
	})
	require.NoError(t, err)
	assert.Equal(t, title, res.Title)
}

func TestUpdate(t *testing.T) {
	var (
		title = "test title"
		id    = 1
	)

	b := &base[Todo]{
		db: gDB,
	}
	mockSQL.MatchExpectationsInOrder(false)
	mockSQL.ExpectBegin()
	mockSQL.ExpectExec(regexp.QuoteMeta(
		`UPDATE "todos"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockSQL.ExpectCommit()
	res, err := b.Update(&Todo{
		Model: gorm.Model{
			ID: uint(id),
		},
		Title: title,
		Done:  true,
	})
	require.NoError(t, err)
	assert.Equal(t, title, res.Title)
}

func TestDelete(t *testing.T) {
	var (
		id = 1
	)
	b := &base[Todo]{
		db: gDB,
	}
	mockSQL.MatchExpectationsInOrder(false)
	mockSQL.ExpectBegin()
	mockSQL.ExpectExec(regexp.QuoteMeta(
		`UPDATE "todos"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockSQL.ExpectCommit()
	err := b.Delete(&Todo{Model: gorm.Model{ID: uint(id)}})
	require.NoError(t, err)
}

func TestFindByID(t *testing.T) {
	var (
		id    = 1
		title = "test title"
		done  = true
	)
	b := &base[Todo]{
		db: gDB,
	}
	mockSQL.MatchExpectationsInOrder(false)
	mockSQL.ExpectBegin()
	const sqlSelectAll = `SELECT * FROM "todos"`
	mockSQL.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "title", "done"}).
			AddRow(id, title, done))
	mockSQL.ExpectCommit()
	todo, err := b.FindById(id)
	require.NoError(t, err)
	assert.Equal(t, title, todo.Title)
	assert.Equal(t, done, todo.Done)
	assert.Equal(t, uint(id), todo.ID)
}

func TestList(t *testing.T) {
	var (
		id    = 1
		title = "test title"
		done  = true
	)
	b := &base[Todo]{
		db: gDB,
	}
	mockSQL.MatchExpectationsInOrder(false)
	mockSQL.ExpectBegin()
	const sqlSelectAll = `SELECT * FROM "todos"`
	mockSQL.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "done"}).
			AddRow(id, title, done))
	mockSQL.ExpectCommit()
	todos, err := b.FindAllByIds([]any{id})
	require.NoError(t, err)
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, title, todos[0].Title)
	assert.Equal(t, done, todos[0].Done)
	assert.Equal(t, uint(id), todos[0].ID)
}
