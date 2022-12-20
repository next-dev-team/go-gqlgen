package model

import (
	"go-graph/db"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TodoRepo interface {
	Base[Todo]
}

type todoRepo struct {
	base[Todo]
}

func NewDefaultTodoRepo() TodoRepo {
	return NewTodoRepo(db.GetConnection())
}

func NewTodoRepo(db *gorm.DB) TodoRepo {
	return &todoRepo{base: base[Todo]{db: db}}
}
