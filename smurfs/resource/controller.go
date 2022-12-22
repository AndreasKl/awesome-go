package resource

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"awesome/apperrors"
	"awesome/smurfs"
)

type Controller struct {
	repository smurfs.Repository
}

func NewController(repository smurfs.Repository) *Controller {
	return &Controller{repository: repository}
}

func (c *Controller) List(res http.ResponseWriter, req *http.Request) {
	response := mapSmurfsToResponse(c.repository.List(req.Context()))
	body, err := json.Marshal(response)
	if err != nil {
		// TODO: Encapsulate error handling, should be a json body with a nice secure error message.
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	_, _ = res.Write(body)
}

func (c *Controller) Get(res http.ResponseWriter, req *http.Request) {
	param := chi.URLParam(req, "id")
	id, err := uuid.Parse(param)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	smurf, err := c.repository.Get(req.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
	}

	response := mapSmurfToResponse(smurf)
	body, err := json.Marshal(response)
	if err != nil {
		// TODO: Encapsulate error handling, should be a json body with a nice secure error message.
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	_, _ = res.Write(body)
}

type SmurfsResponse []SmurfResponse

type SmurfResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Height uint   `json:"height"`
}

func mapSmurfsToResponse(smurfs []smurfs.Smurf) SmurfsResponse {
	response := make(SmurfsResponse, 0, len(smurfs))
	for _, smurf := range smurfs {
		response = append(response, mapSmurfToResponse(smurf))
	}
	return response
}

func mapSmurfToResponse(smurf smurfs.Smurf) SmurfResponse {
	return SmurfResponse{
		ID:     smurf.ID.String(),
		Name:   string(smurf.Name),
		Height: smurf.Height,
	}
}
