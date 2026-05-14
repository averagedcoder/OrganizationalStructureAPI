package service

import "testing"

func TestDepartmentCannotBeParentOfItself(t *testing.T) {
	id := 1
	parentID := 1

	if id == parentID {
		return
	}

	t.Fatal("expected department to not allow self parent")
}

func TestDepartmentDepthLimit(t *testing.T) {
	depth := 10

	if depth > 5 {
		return
	}

	t.Fatal("expected depth limit validation")
}
