package repository

import (
	"errors"

	"github.com/jeffscottbrown/satchel/model"
)

var employeeRepository EmployeeRepository

func SetRepository(r EmployeeRepository) {
	employeeRepository = r
}

type EmployeeRepository interface {
	GetEmployees() ([]model.Employee, error)
	GetEmployeeByEmail(string) (model.Employee, error)
	SaveEmployee(employee *model.Employee) error
}

func SaveEmployee(employee *model.Employee) error {
	if employeeRepository == nil {
		return errors.New("repository has not been initialized")
	}
	return employeeRepository.SaveEmployee(employee)
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

func GetEmployeeByEmail(email string) (*model.Employee, error) {
	if employeeRepository == nil {
		return nil, errors.New("repository has not been initialized")
	}
	employee, err := employeeRepository.GetEmployeeByEmail(email)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}
