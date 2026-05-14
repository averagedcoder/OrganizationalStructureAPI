package model

import "time"

type Department struct {
	ID        int
	Name      string
	ParentID  *int
	CreatedAt time.Time
}
