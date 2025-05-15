package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmployee_AddScore(t *testing.T) {
	e := &Employee{Name: "Alice", Position: "Developer"}

	e.AddScore("performance", 95)
	if val, ok := e.Scores["performance"]; !ok || val != 95 {
		t.Errorf("expected performance score to be 95, got %v", val)
	}

	e.AddScore("attendance", "excellent")
	assert.Equal(t, "excellent", e.Scores["attendance"], "Expected attendance score to be 'excellent'")

	e.AddScore("performance", 100)
	assert.Equal(t, 100, e.Scores["performance"], "Expected performance score to be updated to 100")

	e = &Employee{}
	e.AddScore("teamwork", 88)
	expected := map[string]any{"teamwork": 88}
	assert.Equal(t, expected, e.Scores, "Expected scores map to be initialized with teamwork score")
}
