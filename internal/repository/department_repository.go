package repository

import (
	"github.com/averagedcoder/OrganizationalStructureAPI/internal/model"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(dep *model.Department) error {
	return r.db.Create(dep).Error
}

func (r *DepartmentRepository) GetByID(id int) (*model.Department, error) {
	var dep model.Department

	err := r.db.First(&dep, id).Error
	if err != nil {
		return nil, err
	}

	return &dep, nil
}

func (r *DepartmentRepository) GetChildren(parentID int) ([]model.Department, error) {
	var deps []model.Department

	err := r.db.Where("parent_id = ?", parentID).Find(&deps).Error
	return deps, err
}

func (r *DepartmentRepository) UpdateParent(id int, parentID *int) error {
	return r.db.Model(&model.Department{}).
		Where("id = ?", id).
		Update("parent_id", parentID).Error
}

func (r *DepartmentRepository) Delete(id int) error {
	return r.db.Delete(&model.Department{}, id).Error
}

func (r *DepartmentRepository) GetEmployeesByDepartment(departmentID int) ([]model.Employee, error) {
	var employees []model.Employee

	err := r.db.Where("department_id = ?", departmentID).Find(&employees).Error
	return employees, err
}

func (r *DepartmentRepository) ReassignEmployees(fromID, toID int) error {
	return r.db.Model(&model.Employee{}).
		Where("department_id = ?", fromID).
		Update("department_id", toID).Error
}
