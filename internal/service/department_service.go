package service

import (
	"errors"
	"strings"

	"github.com/averagedcoder/OrganizationalStructureAPI/internal/apperrors"
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/dto"
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/model"
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/repository"
)

type DepartmentService struct {
	repo *repository.DepartmentRepository
}

func NewDepartmentService(repo *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) Create(name string, parentID *int) (*model.Department, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, apperrors.ErrBadRequest
	}

	// получаем всех детей этого parent и проверяем дубликаты
	if parentID != nil {
		children, err := s.repo.GetChildren(*parentID)
		if err != nil {
			return nil, apperrors.ErrConflict
		}

		for _, c := range children {
			if strings.EqualFold(c.Name, name) {
				return nil, apperrors.ErrConflict
			}
		}
	}

	dep := &model.Department{
		Name:     name,
		ParentID: parentID,
	}

	err := s.repo.Create(dep)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	return dep, nil
}

func (s *DepartmentService) GetTree(id int, depth int, includeEmployees bool) (*dto.DepartmentDTO, error) {
	if depth < 0 {
		return nil, nil
	}

	dep, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	result := &dto.DepartmentDTO{
		ID:        dep.ID,
		Name:      dep.Name,
		ParentID:  dep.ParentID,
		CreatedAt: dep.CreatedAt,
	}

	// сотрудники
	if includeEmployees {
		employees, err := s.repo.GetEmployeesByDepartment(id)
		if err == nil {
			for _, e := range employees {
				result.Employees = append(result.Employees, dto.EmployeeDTO{
					ID:        e.ID,
					FullName:  e.FullName,
					Position:  e.Position,
					CreatedAt: e.CreatedAt,
				})
			}
		}
	}

	// если depth = 0 — дальше не идём
	if depth == 0 {
		return result, nil
	}

	children, err := s.repo.GetChildren(id)
	if err != nil {
		return result, apperrors.ErrNotFound
	}

	for _, child := range children {
		childTree, err := s.GetTree(child.ID, depth-1, includeEmployees)
		if err == nil {
			result.Children = append(result.Children, *childTree)
		}
	}

	return result, nil
}

func (s *DepartmentService) Update(id int, name *string, parentID *int) (*model.Department, error) {
	dep, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	// изменение имени
	if name != nil {
		trim := strings.TrimSpace(*name)
		if trim == "" {
			return nil, apperrors.ErrConflict
		}
		dep.Name = trim
	}

	// смена parent
	if parentID != nil {
		if *parentID == id {
			return nil, apperrors.ErrConflict
		}

		// проверка цикла
		children, _ := s.repo.GetChildren(id)
		for _, c := range children {
			if c.ID == *parentID {
				return nil, apperrors.ErrConflict
			}
		}

		dep.ParentID = parentID
	}

	err = s.repo.Create(dep)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	return dep, nil
}

func (s *DepartmentService) Delete(id int, mode string, reassignID *int) error {
	if mode == "reassign" {
		if reassignID == nil {
			return errors.New("reassign_to_department_id is required")
		}

		// перенос сотрудников
		err := s.repo.ReassignEmployees(id, *reassignID)
		if err != nil {
			return apperrors.ErrInternal
		}

		// удалить департамент
		return s.repo.Delete(id)
	}

	// cascade (по умолчанию)
	return s.repo.Delete(id)
}
