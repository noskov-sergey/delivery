package http

import (
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"errors"
)

var _ servers.ServerInterface = (*Server)(nil)

type Server struct {
	createOrderCommandHandler           commands.CreateOrderCommandHandler
	getAllUncompletedOrdersQueryHandler queries.GetAllUncompletedOrdersQueryHandler
	createCourierCommandHandler         commands.CreateCourierCommandHandler
	getAllQueriersQueryHandler          queries.GetAllQueriersQueryHandler
}

func NewServer(createOrderCommandHandler commands.CreateOrderCommandHandler,
	getAllUncompletedOrdersQueryHandler queries.GetAllUncompletedOrdersQueryHandler,
	createCourierCommandHandler commands.CreateCourierCommandHandler,
	getAllQueriersQueryHandler queries.GetAllQueriersQueryHandler,
) (*Server, error) {
	if createOrderCommandHandler == nil {
		return nil, errors.New("createOrderCommandHandler is required")
	}
	if getAllUncompletedOrdersQueryHandler == nil {
		return nil, errors.New("getAllUncompletedOrdersQueryHandler is required")
	}
	if createCourierCommandHandler == nil {
		return nil, errors.New("createCourierCommandHandler is required")
	}
	if getAllQueriersQueryHandler == nil {
		return nil, errors.New("getAllQueriersQueryHandler is required")
	}
	return &Server{
		createOrderCommandHandler:           createOrderCommandHandler,
		getAllUncompletedOrdersQueryHandler: getAllUncompletedOrdersQueryHandler,
		createCourierCommandHandler:         createCourierCommandHandler,
		getAllQueriersQueryHandler:          getAllQueriersQueryHandler,
	}, nil
}
