package http

import (
	"delivery/internal/core/application/usecases/commands"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *Server) CreateCourier(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(WrapError(http.StatusBadRequest, err))
		return
	}
	defer r.Body.Close()

	var req struct {
		Name  string `json:"name"`
		Speed int    `json:"speed"`
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(WrapError(http.StatusBadRequest, err))
		return
	}

	command, err := commands.NewCreateCourierCommand(req.Name, req.Speed)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(WrapError(http.StatusBadRequest, err))
		return
	}

	err = s.createCourierCommandHandler.Handle(r.Context(), command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(WrapError(http.StatusInternalServerError, err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}
