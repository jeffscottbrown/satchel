package repository

import (
	_ "embed"
	"errors"

	"github.com/jeffscottbrown/satchel/model"
)

var employeeRepository EmployeeRepository

func SetRepository(r EmployeeRepository) {
	employeeRepository = r
}

type EmployeeRepository interface {
	GetEmployees() ([]model.Employee, error)
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
