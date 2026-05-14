package model

import "time"

type Employee struct {
	ID           int
	DepartmentID int
	FullName     string
	Position     string
	CreatedAt    time.Time
}
