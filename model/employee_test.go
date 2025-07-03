package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmployee_AddScore(t *testing.T) {
	e := &Employee{Name: "Alice", Position: "Developer"}

	e.AddScore("performance", "95")
	e.AddScore("attendance", "excellent")
	e.AddScore("performance", "100")

	assert.Equal(t, 3, len(e.Scores), "Expected 3 scores to be added")
	assert.Equal(t, "performance", e.Scores[0].Key, "First score key should be 'performance'")
	assert.Equal(t, "95", e.Scores[0].Value, "First score value should be '95'")

	assert.Equal(t, "attendance", e.Scores[1].Key, "Second score key should be 'attendance'")
	assert.Equal(t, "excellent", e.Scores[1].Value, "Second score value should be 'excellent'")

	assert.Equal(t, "performance", e.Scores[2].Key, "Third score key should be 'performance'")
	assert.Equal(t, "100", e.Scores[2].Value, "Third score value should be '100'")
}
