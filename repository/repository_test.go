package repository

import (
	"testing"

	"github.com/jeffscottbrown/satchel/model"
	"github.com/stretchr/testify/assert"
)

func TestGetEmployees_RepositoryNotInitialized(t *testing.T) {
	// Ensure repository is not initialized
	ConfigureRepositoryForTest(t, nil)

	employees, err := GetEmployees()
	assert.Nil(t, employees)
	assert.Error(t, err)
	assert.EqualError(t, err, "repository has not been initialized")
}

func TestGetEmployees_RepositoryInitialized(t *testing.T) {
	_, err := GetEmployees()
	assert.NoError(t, err)
}

func TestGetEmployees_RepositoryReturnsError(t *testing.T) {
	t.Skip()

	employees, err := GetEmployees()
	assert.Nil(t, employees)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
}

func TestGetEmployeeByName_Found(t *testing.T) {
	SaveEmployee(&model.Employee{
		Email: "bob@somewhere.com",
		Name:  "Bob",
	})

	employee, err := GetEmployeeByEmail("bob@somewhere.com")
	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, "Bob", employee.Name)
}

func TestGetEmployeeByName_NotFound(t *testing.T) {
	_, err := GetEmployeeByEmail("charlie@somewhere.com")
	assert.Error(t, err)
	assert.EqualError(t, err, "record not found")
}

func TestGetEmployeeByName_RepositoryNotInitialized(t *testing.T) {
	ConfigureRepositoryForTest(t, nil)

	employee, err := GetEmployeeByEmail("alice@somewhere.com")
	assert.Nil(t, employee)
	assert.Error(t, err)
	assert.EqualError(t, err, "repository has not been initialized")
}

func TestGetEmployeeByName_RepositoryReturnsError(t *testing.T) {
	employee, err := GetEmployeeByEmail("alice@somewhere.com")
	assert.Nil(t, employee)
	assert.Error(t, err)
	assert.EqualError(t, err, "record not found")
}

func TestMain(m *testing.M) {
	RunTestsWithTestContainer(m)
}
