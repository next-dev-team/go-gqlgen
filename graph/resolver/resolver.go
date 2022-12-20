package resolver

import (
	"go-graph/db/model"
	"go-graph/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// add on demand services here
	todoSvc *service.ServiceTodo
}

func New() *Resolver {
	return &Resolver{
		// create a new service here
		todoSvc: service.NewServiceTodo(model.NewDefaultTodoRepo()),
	}
}
