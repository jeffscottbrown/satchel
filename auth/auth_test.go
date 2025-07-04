package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDomainAllowed(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"", false}, // Empty email
		{"someone@objectcomputing.com", true},
		{"someone@someotherdomain.com", false}}
	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			result := isAllowedDomain(test.email)
			assert.Equal(t, test.expected, result, "Expected domain check to match")
		})
	}
}
