package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/averagedcoder/OrganizationalStructureAPI/internal/apperrors"
)

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}
type createDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}
type updateDepartmentRequest struct {
	Name     *string `json:"name"`
	ParentID *int    `json:"parent_id"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Response{
		Data: data,
	})
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var status int

	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		status = http.StatusNotFound

	case errors.Is(err, apperrors.ErrConflict):
		status = http.StatusConflict

	case errors.Is(err, apperrors.ErrBadRequest):
		status = http.StatusBadRequest

	default:
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(Response{
		Error: err.Error(),
	})
}
