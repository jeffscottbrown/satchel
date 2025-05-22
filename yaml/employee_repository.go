package yaml

import (
	_ "embed"

	"github.com/jeffscottbrown/satchel/model"
	"github.com/jeffscottbrown/satchel/repository"
	"gopkg.in/yaml.v3"
)

//go:embed employees.yaml
var employeesYAML []byte

type YamlEmployeeRepository struct{}

func (r *YamlEmployeeRepository) GetEmployees() ([]model.Employee, error) {
	var employees []model.Employee
	if err := yaml.Unmarshal(employeesYAML, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *YamlEmployeeRepository) GetEmployeeByName(name string) (*model.Employee, error) {
	employees, err := r.GetEmployees()
	if err != nil {
		return nil, err
	}
	for _, employee := range employees {
		if employee.Name == name {
			return &employee, nil
		}
	}
	return nil, nil
}

func init() {
	repository.SetRepository(&YamlEmployeeRepository{})
}
