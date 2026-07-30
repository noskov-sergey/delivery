package http

import (
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) GetCouriers(w http.ResponseWriter, r *http.Request) {
	query, err := queries.NewGetAllQueriersQuery()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(WrapError(http.StatusBadRequest, err))
		return
	}

	couriers, err := s.getAllQueriersQueryHandler.Handle(r.Context(), query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(WrapError(http.StatusInternalServerError, err))
		return
	}

	bytes, err := dtoToResponse(couriers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(WrapError(http.StatusInternalServerError, err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

func dtoToResponse(response *queries.GetAllCouriersResponse) ([]byte, error) {
	couriers := make([]servers.Courier, len(response.Couriers))
	for i, c := range response.Couriers {
		couriers[i] = servers.Courier{
			Id:   c.ID,
			Name: c.Name,
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
