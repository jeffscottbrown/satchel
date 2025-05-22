package repository

import (
	"errors"

	"github.com/jeffscottbrown/satchel/model"
)

var employeeRepository EmployeeRepository

type EmployeeRepository interface {
	GetEmployees() ([]model.Employee, error)
	GetEmployeeByName(string) (*model.Employee, error)
}

func SetRepository(repo EmployeeRepository) {
	employeeRepository = repo
}

func GetEmployees() ([]model.Employee, error) {
	if employeeRepository == nil {
		return nil, errors.New("repository has not been initialized")
	}
	employees, err := employeeRepository.GetEmployees()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func GetEmployeeByName(name string) (*model.Employee, error) {
	employees, err := GetEmployees()
	if err != nil {
		return nil, err
	}
	for _, employee := range employees {
		if employee.Name == name {
			return &employee, nil
		}
	}
	return nil, errors.New("employee not found")
}
