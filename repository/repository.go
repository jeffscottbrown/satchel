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
	DeleteReflection(reflectionId uint) error
	DeleteEmployee(email string) error
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

func DeleteEmployee(email string) error {
	if employeeRepository == nil {
		return errors.New("repository has not been initialized")
	}
	return employeeRepository.DeleteEmployee(email)
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

func SavePosition(email string, position string) error {
	employee, err := GetEmployeeByEmail(email)
	if err != nil {
		return err
	}
	employee.Position = position
	return SaveEmployee(employee)
}

func SaveBio(email string, bio string) error {
	employee, err := GetEmployeeByEmail(email)
	if err != nil {
		return err
	}
	employee.Bio = bio
	return SaveEmployee(employee)
}

func DeleteReflection(email string, reflectionId uint) error {
	employee, err := GetEmployeeByEmail(email)
	if err != nil {
		return err
	}
	found := false
	for _, reflection := range employee.Reflections {
		if reflection.ID == reflectionId {
			found = true
			break
		}
	}
	if !found {
		return errors.New("reflection not found")
	}
	return employeeRepository.DeleteReflection(reflectionId)
}
func AddReflection(email string, name string, value string) error {
	employee, err := GetEmployeeByEmail(email)
	if err != nil {
		return err
	}
	employee.AddReflection(name, value)
	return SaveEmployee(employee)
}
