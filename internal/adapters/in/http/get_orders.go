package http

import (
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	query, err := queries.NewGetAllUncompletedOrdersQuery()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(WrapError(http.StatusBadRequest, err))
		return
	}

	orders, err := s.getAllUncompletedOrdersQueryHandler.Handle(r.Context(), query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(WrapError(http.StatusInternalServerError, err))
		return
	}

	bytes, err := dtoOrderToResponse(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(WrapError(http.StatusInternalServerError, err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

func dtoOrderToResponse(response *queries.GetAllUncompletedOrdersResponse) ([]byte, error) {
	couriers := make([]servers.Order, len(response.Orders))
	for i, c := range response.Orders {
		couriers[i] = servers.Order{
			Id: c.ID,
			Location: servers.Location{
				X: c.Location.X(),
				Y: c.Location.Y(),
			},
		}
	}

	res, err := json.MarshalIndent(couriers, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response")
	}

	return res, nil
}
