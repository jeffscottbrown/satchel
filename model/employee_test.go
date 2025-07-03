package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmployee_AddReflection(t *testing.T) {
	e := &Employee{Name: "Alice", Position: "Developer"}

	e.AddReflection("performance", "95")
	e.AddReflection("attendance", "excellent")
	e.AddReflection("performance", "100")

	assert.Equal(t, 3, len(e.Reflections), "Expected 3 scores to be added")
	assert.Equal(t, "performance", e.Reflections[0].Key, "First score key should be 'performance'")
	assert.Equal(t, "95", e.Reflections[0].Value, "First score value should be '95'")

	assert.Equal(t, "attendance", e.Reflections[1].Key, "Second score key should be 'attendance'")
	assert.Equal(t, "excellent", e.Reflections[1].Value, "Second score value should be 'excellent'")

	assert.Equal(t, "performance", e.Reflections[2].Key, "Third score key should be 'performance'")
	assert.Equal(t, "100", e.Reflections[2].Value, "Third score value should be '100'")
}
