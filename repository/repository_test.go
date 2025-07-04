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

func TestDeleteEmployee(t *testing.T) {
	email := "test@somedomain.com"
	emp, err := GetEmployeeByEmail(email)
	assert.Error(t, err)
	assert.Nil(t, emp)

	err = SaveEmployee(&model.Employee{
		Email: email,
	})
	assert.NoError(t, err)

	emp, err = GetEmployeeByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, emp)

	err = DeleteEmployee(email)
	assert.NoError(t, err)
	emp, err = GetEmployeeByEmail(email)
	assert.Error(t, err)
	assert.Nil(t, emp)
}

func TestSavePosition(t *testing.T) {
	email := "someone@somewhere.com"
	t.Cleanup(func() {
		err := DeleteEmployee(email)
		assert.NoError(t, err)
	})
	SaveEmployee(&model.Employee{
		Email: email,
	})
	emp, err := GetEmployeeByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, "", emp.Position)

	err = SavePosition(email, "Some New Position")

	emp, err = GetEmployeeByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, "Some New Position", emp.Position)
}

func TestSaveBio(t *testing.T) {
	email := "someone@somewhere.com"
	t.Cleanup(func() {
		err := DeleteEmployee(email)
		assert.NoError(t, err)
	})
	SaveEmployee(&model.Employee{
		Email: email,
	})
	emp, err := GetEmployeeByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, "", emp.Bio)

	err = SaveBio(email, "Some New Bio")

	emp, err = GetEmployeeByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, "Some New Bio", emp.Bio)
}

func TestUpdatingReflections(t *testing.T) {
	email := "someone@someplace.com"
	t.Cleanup(func() {
		err := DeleteEmployee(email)
		assert.NoError(t, err)
	})
	SaveEmployee(&model.Employee{
		Email: email,
	})
	emp, err := GetEmployeeByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Empty(t, emp.Reflections)

	AddReflection(email, "Favorite Band", "Grateful Dead")
	AddReflection(email, "Home", "Here")

	assert.NoError(t, err)

	emp, err = GetEmployeeByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Len(t, emp.Reflections, 2)

	var homeReflectionID uint
	for _, r := range emp.Reflections {
		if r.Key == "Home" {
			homeReflectionID = r.ID
			break
		}
	}

	err = DeleteReflection(email, homeReflectionID)
	assert.NoError(t, err)

	emp, err = GetEmployeeByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Len(t, emp.Reflections, 1)
	assert.Equal(t, "Favorite Band", emp.Reflections[0].Key)

}

func TestMain(m *testing.M) {
	RunTestsWithTestContainer(m)
}
