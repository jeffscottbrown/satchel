package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/jeffscottbrown/satchel/yaml"
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
