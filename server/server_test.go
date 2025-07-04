package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/jeffscottbrown/satchel/model"
	"github.com/jeffscottbrown/satchel/repository"
	"github.com/stretchr/testify/assert"
)

func TestRootEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := createRouter()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")

	assert.Contains(t, recorder.Body.String(), "<title>Satchel</title>", "Page title should be 'Satchel'")
}

func TestRootHandler_GetEmployeeesError(t *testing.T) {
	repository.ConfigureRepositoryForTest(t, &errorThrowingEmployeeRepository{})

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	rootHandler(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status(), "Expected status code 500")
	assert.Equal(t, w.Body.String(), "Error retrieving employees: An error occurred retrieving employees", "Expected error message in response")
}

func TestEmployeeHandler(t *testing.T) {
	repository.SaveEmployee(&model.Employee{
		// Name:      "Henry David Thoreau",
		Email: "henry@thewods.org",
		Reflections: []model.Reflection{
			{
				Key:   "Contemplative",
				Value: "10",
			},
			{
				Key:   "Social",
				Value: "3",
			},
			{
				Key:   "Favorite Place",
				Value: "The Pond",
			},
		},
	})
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "employeeEmail", Value: "henry@thewods.org"},
	}
	c.Request = httptest.NewRequest("GET", "/", nil)

	employeeHandler(c)
	assert.Equal(t, http.StatusOK, c.Writer.Status(), "Expected status code 200")

	doc, err := goquery.NewDocumentFromReader(w.Body)
	assert.NoError(t, err, "Expected no error parsing HTML")
	expected := []struct {
		Label string
		Value string
	}{
		{"Contemplative:", "10"},
		{"Favorite Place:", "Walden Pond"},
		{"Social:", "2"},
	}
	doc.Find(".card-stats .stat").Each(func(i int, s *goquery.Selection) {
		if i >= len(expected) {
			t.Errorf("Unexpected extra .stat element at index %d", i)
			return
		}

		label := s.Find(".label").Text()
		value := s.Find(".value").Text()

		if label != expected[i].Label {
			t.Errorf("Label mismatch at index %d: got %q, want %q", i, label, expected[i].Label)
		}
		if value != expected[i].Value {
			t.Errorf("Value mismatch at index %d: got %q, want %q", i, value, expected[i].Value)
		}
	})

	if got := doc.Find(".card-stats table tbody tr").Length(); got != len(expected) {
		t.Errorf("Expected %d stat rows, got %d", len(expected), got)
	}
}

type errorThrowingEmployeeRepository struct {
}

// DeleteReflection implements repository.EmployeeRepository.
func (m *errorThrowingEmployeeRepository) DeleteReflection(reflectionId uint) error {
	panic("unimplemented")
}

// SaveEmployee implements repository.EmployeeRepository.
func (m *errorThrowingEmployeeRepository) SaveEmployee(employee *model.Employee) error {
	panic("unimplemented")
}

func (m *errorThrowingEmployeeRepository) GetEmployees() ([]model.Employee, error) {
	return nil, errors.New("An error occurred retrieving employees")
}

func (m *errorThrowingEmployeeRepository) GetEmployeeByEmail(email string) (model.Employee, error) {
	return model.Employee{}, errors.New("An error occurred retrieving employee by email")
}

func TestMain(m *testing.M) {
	repository.RunTestsWithTestContainer(m)
}
