package http

import (
	"delivery/internal/core/application/usecases/commands"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	street := "Тестовая улица, д." + strconv.Itoa(rand.Intn(15)+1)

	command, err := commands.NewCreateOrderCommand(uuid.New(), street, 5)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(WrapError(http.StatusBadRequest, err))
		return
	}

	err = s.createOrderCommandHandler.Handle(r.Context(), command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(WrapError(http.StatusInternalServerError, err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}
