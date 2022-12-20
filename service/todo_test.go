package service

import (
	"context"
	"go-graph/db/model"
	"go-graph/graph/modelgen"
	testutil "go-graph/test"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupServiceTodo(mockRepo *testutil.MockRepo[model.Todo]) *ServiceTodo {
	return NewServiceTodo(mockRepo)
}

func TestNewTodo(t *testing.T) {
	var (
		text   = "task 1"
		userId = "user-1"
		id     = 1
	)
	mockRepo := &testutil.MockRepo[model.Todo]{
		Model: &model.Todo{
			Model: gorm.Model{
				ID: uint(id),
			},
			Title: text,
			Done:  false,
		},
	}
	s := setupServiceTodo(mockRepo)
	res, err := s.NewTodo(context.Background(), &modelgen.NewTodo{
		Text:   text,
		UserID: userId,
	})
	require.NoError(t, err)
	assert.Equal(t, text, res.Text)
	assert.Equal(t, false, res.Done)
	assert.Equal(t, id, res.ID)
}

func TestGetTodo(t *testing.T) {
	var (
		text = "task 1"
		id   = 1
		done = true
	)
	mockRepo := &testutil.MockRepo[model.Todo]{
		Model: &model.Todo{
			Model: gorm.Model{
				ID: uint(id),
			},
			Title: text,
			Done:  done,
		},
	}
	s := setupServiceTodo(mockRepo)
	res, err := s.GetTodo(context.Background(), strconv.Itoa(id))
	require.NoError(t, err)
	assert.Equal(t, text, res.Text)
	assert.Equal(t, done, res.Done)
	assert.Equal(t, id, res.ID)
}

func TestGetTodos(t *testing.T) {
	var (
		text = "task 1"
		id   = 1
		done = true
	)
	mockRepo := &testutil.MockRepo[model.Todo]{
		Models: []*model.Todo{
			{
				Model: gorm.Model{
					ID: uint(id),
				},
				Title: text,
				Done:  done,
			},
		},
	}
	s := setupServiceTodo(mockRepo)
	res, err := s.GetTodos(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, text, res[0].Text)
	assert.Equal(t, done, res[0].Done)
	assert.Equal(t, id, res[0].ID)
}
