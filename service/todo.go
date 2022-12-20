package service

import (
	"context"
	"go-graph/db/model"
	"go-graph/graph/modelgen"
)

type ServiceTodo struct {
	repo model.TodoRepo
}

func NewServiceTodo(repo model.TodoRepo) *ServiceTodo {
	return &ServiceTodo{
		repo: repo,
	}
}

func (s *ServiceTodo) NewTodo(ctx context.Context, input *modelgen.NewTodo) (*modelgen.Todo, error) {
	res, err := s.repo.Create(&model.Todo{
		Title: input.Text,
	})
	if err != nil {
		return nil, err
	}
	return &modelgen.Todo{ID: int(res.ID), Text: res.Title, Done: res.Done}, nil
}

func (s *ServiceTodo) GetTodo(ctx context.Context, id string) (*modelgen.Todo, error) {
	res, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return &modelgen.Todo{ID: int(res.ID), Text: res.Title, Done: res.Done}, nil
}
func (s *ServiceTodo) GetTodos(ctx context.Context) ([]*modelgen.Todo, error) {
	res, err := s.repo.FindAllByIds([]any{})
	if err != nil {
		return nil, err
	}
	todos := make([]*modelgen.Todo, len(res))
	for i, v := range res {
		todos[i] = &modelgen.Todo{ID: int(v.ID), Text: v.Title, Done: v.Done}
	}
	return todos, nil
}
