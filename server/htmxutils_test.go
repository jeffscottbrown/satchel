package server

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTemplateToString(t *testing.T) {

	testTemplate := template.Must(template.New("test").Parse(`
		{{define "example"}}Hello, {{.Name}}!{{end}}
	`))

	data := map[string]string{"Name": "World"}

	result := renderTemplateToString("example", data, testTemplate)

	expected := "Hello, World!"
	assert.Equal(t, expected, result, "The rendered template string should match the expected output")
}
