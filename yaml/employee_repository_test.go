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
	assert.NoError(t, err, "expected no error when getting employee by name")
	assert.Nil(t, employee, "expected to not find employee Nonexistent Employee")
}
