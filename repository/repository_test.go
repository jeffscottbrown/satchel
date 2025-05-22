package repository

import (
	"errors"
	"testing"

	"github.com/jeffscottbrown/satchel/model"
	"github.com/stretchr/testify/assert"
)

// mock implementation of EmployeeRepository
type mockEmployeeRepository struct {
	employees []model.Employee
	err       error
}

func (m *mockEmployeeRepository) GetEmployees() ([]model.Employee, error) {
	return m.employees, m.err
}

func (m *mockEmployeeRepository) GetEmployeeByName(name string) (*model.Employee, error) {
	for _, employee := range m.employees {
		if employee.Name == name {
			return &employee, nil
		}
	}
	return nil, errors.New("employee not found")
}

func TestGetEmployees_RepositoryNotInitialized(t *testing.T) {
	// Ensure repository is not initialized
	ConfigureRepositoryForTest(t, nil)

	employees, err := GetEmployees()
	assert.Nil(t, employees)
	assert.Error(t, err)
	assert.EqualError(t, err, "repository has not been initialized")
}

func TestGetEmployees_RepositoryInitialized(t *testing.T) {
	mockEmployees := []model.Employee{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	mockRepo := &mockEmployeeRepository{
		employees: mockEmployees,
		err:       nil,
	}
	ConfigureRepositoryForTest(t, mockRepo)

	employees, err := GetEmployees()
	assert.NoError(t, err)
	assert.Equal(t, mockEmployees, employees)
}

func TestGetEmployees_RepositoryReturnsError(t *testing.T) {
	mockRepo := &mockEmployeeRepository{
		employees: nil,
		err:       errors.New("database error"),
	}
	ConfigureRepositoryForTest(t, mockRepo)

	employees, err := GetEmployees()
	assert.Nil(t, employees)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
}

func TestGetEmployeeByName_Found(t *testing.T) {
	mockEmployees := []model.Employee{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	mockRepo := &mockEmployeeRepository{
		employees: mockEmployees,
		err:       nil,
	}
	ConfigureRepositoryForTest(t, mockRepo)

	employee, err := GetEmployeeByName("Bob")
	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, "Bob", employee.Name)
}

func TestGetEmployeeByName_NotFound(t *testing.T) {
	mockEmployees := []model.Employee{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	mockRepo := &mockEmployeeRepository{
		employees: mockEmployees,
		err:       nil,
	}
	ConfigureRepositoryForTest(t, mockRepo)

	employee, err := GetEmployeeByName("Charlie")
	assert.Nil(t, employee)
	assert.Error(t, err)
	assert.EqualError(t, err, "employee not found")
}

func TestGetEmployeeByName_RepositoryNotInitialized(t *testing.T) {
	ConfigureRepositoryForTest(t, nil)

	employee, err := GetEmployeeByName("Alice")
	assert.Nil(t, employee)
	assert.Error(t, err)
	assert.EqualError(t, err, "repository has not been initialized")
}

func TestGetEmployeeByName_RepositoryReturnsError(t *testing.T) {
	mockRepo := &mockEmployeeRepository{
		employees: nil,
		err:       errors.New("database error"),
	}
	ConfigureRepositoryForTest(t, mockRepo)

	employee, err := GetEmployeeByName("Alice")
	assert.Nil(t, employee)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
}
