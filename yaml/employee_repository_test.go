package yaml

import (
	"testing"

	"github.com/jeffscottbrown/satchel/model"
	"github.com/stretchr/testify/assert"
)

func TestGetEmployees(t *testing.T) {

	repo := &YamlEmployeeRepository{}
	employees, err := repo.GetEmployees()
	assert.NoError(t, err, "expected no error when getting employees")
	assert.Len(t, employees, 3, "expected 3 employees")

	var henry *model.Employee
	for _, e := range employees {
		if e.Name == "Henry David Thoreau" {
			henry = &e
			break
		}
	}
	assert.NotNil(t, henry, "expected to find employee Henry Thoreau")
	assert.Equal(t, 10, henry.Scores["Contemplative"], "expected Henry Thoreau to have a Contemplative score of 10")
}

func TestGetEmployeeByName_Found(t *testing.T) {
	repo := &YamlEmployeeRepository{}
	employee, err := repo.GetEmployeeByName("Henry David Thoreau")
	assert.NoError(t, err, "expected no error when getting employee by name")
	assert.NotNil(t, employee, "expected to find employee Henry Thoreau")
	assert.Equal(t, "Henry David Thoreau", employee.Name, "expected employee name to be Henry David Thoreau")
}
func TestGetEmployeeByName_NotFound(t *testing.T) {
	repo := &YamlEmployeeRepository{}
	employee, err := repo.GetEmployeeByName("Nonexistent Employee")
	assert.Error(t, err, "expected error when getting nonexistent employee")
	assert.Equal(t, model.Employee{}, employee, "expected employee to be empty struct")
	assert.Equal(t, "employee not found", err.Error(), "expected error message to be 'employee not found'")
}
