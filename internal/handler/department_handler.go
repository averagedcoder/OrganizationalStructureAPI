package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/averagedcoder/OrganizationalStructureAPI/internal/service"
)

type DepartmentHandler struct {
	service *service.DepartmentService
}

func NewDepartmentHandler(service *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/departments/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// параметры
	depth := 1
	includeEmployees := true

	if d := r.URL.Query().Get("depth"); d != "" {
		if val, err := strconv.Atoi(d); err == nil {
			if val > 5 {
				val = 5
			}
			depth = val
		}
	}

	if ie := r.URL.Query().Get("include_employees"); ie == "false" {
		includeEmployees = false
	}

	tree, err := h.service.GetTree(id, depth, includeEmployees)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, http.StatusOK, tree)
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req createDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	dep, err := h.service.Create(req.Name, req.ParentID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dep)
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/departments/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req updateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	dep, err := h.service.Update(id, req.Name, req.ParentID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dep)
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/departments/"):]
	id, _ := strconv.Atoi(idStr)

	mode := r.URL.Query().Get("mode")
	reassignIDStr := r.URL.Query().Get("reassign_to_department_id")

	var reassignID *int
	if reassignIDStr != "" {
		val, _ := strconv.Atoi(reassignIDStr)
		reassignID = &val
	}

	err := h.service.Delete(id, mode, reassignID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
