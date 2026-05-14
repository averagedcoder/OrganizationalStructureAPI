package dto

import "time"

type EmployeeDTO struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

type DepartmentDTO struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`

	Employees []EmployeeDTO   `json:"employees,omitempty"`
	Children  []DepartmentDTO `json:"children,omitempty"`
}
